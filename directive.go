package ntgo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

type DirectiveType int

const (
	DirectiveTypeUnknown    DirectiveType = iota
	DirectiveTypeString     DirectiveType = iota
	DirectiveTypeText       DirectiveType = iota
	DirectiveTypeList       DirectiveType = iota
	DirectiveTypeDictionary DirectiveType = iota
	DirectiveTypeComment    DirectiveType = iota
)

var (
	NextDirectiveAppearedError        = errors.New("ntgo: next directive appeared")
	DifferentTypesOnTheSameLevelError = errors.New("ntgo: can not place different types of entities on the same level")
	DictionaryKeyNestedQuotesError    = errors.New("ntgo: quoted dictionary key can not contain any quotes")
	EmptyDataError                    = errors.New("ntgo: data can not be empty")
	RootLevelHasIndentError           = errors.New("ntgo: root level must not be indented")
	TabInIndentationError             = errors.New("ntgo: indent can not contain tab")
	RootStringError                   = errors.New("ntgo: no string allowed on root level")
	StringHasChildError               = errors.New("ntgo: string type can not have child")
	TextHasChildError                 = errors.New("ntgo: text type can not have child")
	DifferentLevelOnSameChildError    = errors.New("ntgo: child elements have dirfferent leves")
	StringWithNewLineError            = errors.New("ntgo: string type can not have line break")
	DictionaryDuplicateKeyError       = errors.New("ntgo: dictionary type can not have the same key")
)

type Directive struct {
	Type DirectiveType

	String     string
	Text       []string
	List       []*Directive
	Dictionary map[string]*Directive

	IndentSize int
	Depth      int
}

func (d *Directive) Unmarshal() string {
	str := ""

	baseIndent := fmt.Sprintf("%*s", d.IndentSize*d.Depth, "")
	switch d.Type {
	case DirectiveTypeString:
		{
			str = d.String
		}
	case DirectiveTypeText:
		{
			for i := 0; i < len(d.Text); i++ {
				str = fmt.Sprintf("%s%s> %s", str, baseIndent, d.Text[i])
			}
		}
	case DirectiveTypeList:
		{
			for i := 0; i < len(d.List); i++ {
				dataLn := string(LineBreak)
				tailLn := string(LineBreak)
				if i == len(d.List)-1 {
					tailLn = ""
				}

				child := d.List[i]
				if child.Type == DirectiveTypeString {
					dataLn = ""
				}

				str = fmt.Sprintf("%s%s- %s%s%s", str, baseIndent, dataLn, child.Unmarshal(), tailLn)
			}
		}
	case DirectiveTypeDictionary:
		{
			it := 0
			for k, v := range d.Dictionary {
				dataLn := string(LineBreak)
				tailLn := string(LineBreak)
				if it == len(d.Dictionary)-1 {
					tailLn = ""
				}

				if v.Type == DirectiveTypeString {
					dataLn = ""
				}

				str = fmt.Sprintf("%s%s%s: %s%s%s", str, baseIndent, k, dataLn, v.Unmarshal(), tailLn)

				it++
			}
		}
	}

	return str
}

func (d *Directive) Marshal(content []byte) (err error) {
	d.Type = DirectiveTypeUnknown

	buffer := bytes.NewBuffer(removeTrailingLineBreaks(content))

	var currentLine []byte
	var index int
	loadedNextLine := false

	for eof := false; !eof; {
		if !loadedNextLine {
			if currentLine, err = buffer.ReadBytes(byte(LineBreak)); err == io.EOF {
				err = nil
				eof = true
			} else if err != nil {
				break
			}
		}

		var directiveType DirectiveType
		if directiveType, index, err = detectDirectiveType(currentLine); err != nil {
			break
		}

		switch directiveType {
		case DirectiveTypeUnknown, DirectiveTypeComment:
			{
				loadedNextLine = false
				err = nil
			}
		case DirectiveTypeText:
			currentLine, loadedNextLine, err = d.readTextDirective(index, currentLine, buffer)
		case DirectiveTypeList:
			currentLine, loadedNextLine, err = d.readListDirective(index, currentLine, buffer)
		default:
			currentLine, loadedNextLine, err = d.readDictionaryDirective(index, currentLine, buffer)
		}

		if err != nil {
			break
		}
	}

	if err == nil && d.Type == DirectiveTypeUnknown {
		err = EmptyDataError
	}

	return
}

func removeTrailingLineBreaks(data []byte) []byte {
	for len(data) > 0 && data[len(data)-1] == LineBreak {
		data = data[:len(data)-1]
	}
	return data
}

func detectDirectiveType(line []byte) (DirectiveType, int, error) {
	chars, newIndex := readFirstMeaningfulTwoCharacters(line)

	if chars[0] == Tab {
		return DirectiveTypeUnknown, newIndex, TabInIndentationError
	}

	if chars[0] == CommentSymbol {
		return DirectiveTypeComment, newIndex, nil
	}

	directiveType := DirectiveTypeUnknown

	switch chars[0] {
	case EmptyChar:
	case CommentSymbol:
	case TextSymbol:
		{
			if chars[1] == Space || chars[1] == LineBreak || chars[1] == EmptyChar {
				directiveType = DirectiveTypeText
			} else {
				directiveType = DirectiveTypeString
			}
		}
	case ListSymbol:
		{
			if chars[1] == Space || chars[1] == LineBreak || chars[1] == EmptyChar {
				directiveType = DirectiveTypeList
			} else {
				directiveType = DirectiveTypeString
			}
		}
	default:
		directiveType = DirectiveTypeString
	}

	return directiveType, newIndex, nil
}

func (d *Directive) readTextDirective(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	if d.Type != DirectiveTypeUnknown {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}

	d.Type = DirectiveTypeText

	var readErr error
	currentLine := initialLine

	for {
		char, newIndentSpaces := readFirstMeaningfulCharacter(currentLine, false)

		if char != CommentSymbol && char != LineBreak && newIndentSpaces != NotFoundIndex {
			// validate
			var err error
			{
				if newIndentSpaces > baseIndentSpaces {
					// deeper
					if char == TextSymbol {
						err = DifferentLevelOnSameChildError
					} else {
						err = TextHasChildError
					}
				} else if newIndentSpaces < baseIndentSpaces {
					// shallower
					if char == TextSymbol {
						err = DifferentLevelOnSameChildError
					} else {
						hasNext = true
					}
				} else if char != TextSymbol {
					err = DifferentLevelOnSameChildError
				}
			}
			if hasNext {
				if len(d.Text) >= 1 {
					removeTrailingLineBreak(&d.Text[len(d.Text)-1])
				}
				return currentLine, hasNext, nil
			}
			if err != nil {
				return nil, hasNext, err
			}

			// append text
			var appendText string
			{
				if len(currentLine) <= newIndentSpaces+1 {
					// text ends with symbol
					appendText = ""
				} else if currentLine[newIndentSpaces+1] == LineBreak {
					// text symbol with no space
					appendText = string(LineBreak)
				} else {
					// after text symbol(>) and space
					appendText = string(currentLine[newIndentSpaces+2:])
				}
			}
			d.Text = append(d.Text, appendText)
		}

		if readErr == io.EOF {
			if len(d.Text) >= 1 {
				removeTrailingLineBreak(&d.Text[len(d.Text)-1])
			}
			return nil, hasNext, nil
		} else if readErr != nil {
			return nil, hasNext, readErr
		}

		currentLine, readErr = buffer.ReadBytes(byte(LineBreak))
	}

	return currentLine, hasNext, nil
}

func (d *Directive) readListDirective(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	if d.Type != DirectiveTypeUnknown && d.Type != DirectiveTypeList {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}

	d.Type = DirectiveTypeList
	elementContent := initialLine[baseIndentSpaces+1:]

	// detect string
	elementContentChar, _ := readFirstMeaningfulCharacter(elementContent, true)

	// string case
	if elementContentChar != EmptyChar {
		var nextLine []byte
		var err error

		if nextLine, err = buffer.ReadBytes(byte(LineBreak)); err == io.EOF {
			err = nil
		} else if err != nil {
			return nil, hasNext, err
		}

		child := &Directive{
			Type:       DirectiveTypeString,
			IndentSize: d.IndentSize,
			Depth:      d.Depth + 1,
		}

		if initialLine[len(initialLine)-1] == LineBreak {
			child.String = string(initialLine[baseIndentSpaces+2 : len(initialLine)-1])
		} else {
			child.String = string(initialLine[baseIndentSpaces+2:])
		}

		d.List = append(d.List, child)

		char, newIndentSpaces := readFirstMeaningfulCharacter(nextLine, true)

		// validate
		if char == Tab {
			err = TabInIndentationError
		} else if newIndentSpaces != NotFoundIndex {
			if newIndentSpaces == baseIndentSpaces {
				hasNext = true
			} else if newIndentSpaces > baseIndentSpaces {
				err = StringHasChildError
			} else if newIndentSpaces < baseIndentSpaces {
				err = DifferentLevelOnSameChildError
			}
		}

		return nextLine, hasNext, err
	}

	var currentLine []byte
	var err error

	// collect child content lines
	for err != io.EOF {
		if currentLine, err = buffer.ReadBytes(byte(LineBreak)); err != nil && err != io.EOF {
			break
		}

		char, newIndentSpaces := readFirstMeaningfulCharacter(currentLine, true)

		// validate
		if char == Tab {
			err = TabInIndentationError
			break
		}
		if newIndentSpaces == baseIndentSpaces {
			hasNext = true
			break
		}

		elementContent = append(elementContent, currentLine...)
	}

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return nil, hasNext, err
	}

	child := &Directive{
		IndentSize: d.IndentSize,
		Depth:      d.Depth + 1,
	}

	// marshal child
	// TODO: elementContent internally converted to bytes.Buufer, inpsect its performance cost
	if marshalErr := child.Marshal(elementContent); marshalErr != nil {
		if marshalErr != EmptyDataError {
			return nil, hasNext, marshalErr
		}
		// treat empty data as empty string
		child.Type = DirectiveTypeString
		child.String = ""
	}

	d.List = append(d.List, child)

	return currentLine, hasNext, err
}

func (d *Directive) readDictionaryDirective(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	var err error

	// dictionary
	if d.Type != DirectiveTypeUnknown && d.Type != DirectiveTypeDictionary {
		err = DifferentTypesOnTheSameLevelError
	} else if d.Depth == 0 && baseIndentSpaces > 0 {
		err = RootLevelHasIndentError
	}

	if err != nil {
		return nil, hasNext, err
	}

	key, valueIndex := detectKeyBytes(initialLine)

	// unexpected string
	if key == nil && valueIndex == NotFoundIndex {
		return nil, hasNext, RootStringError
	}

	if err = sanitizeDictionaryKey(&key); err != nil {
		return nil, hasNext, err
	}
	if d.Dictionary != nil {
		if _, exists := d.Dictionary[string(key)]; exists {
			err = DictionaryDuplicateKeyError
		}
	}

	d.Type = DirectiveTypeDictionary

	currentLine := initialLine
	elementContent := currentLine[valueIndex:]

	firstChar, _ := readFirstMeaningfulCharacter(elementContent, true)

	if firstChar != EmptyChar {
		var char byte
		var nextIndex int

		for err != io.EOF {
			if currentLine, err = buffer.ReadBytes(byte(LineBreak)); err != nil && err != io.EOF {
				break
			}

			char, nextIndex = readFirstMeaningfulCharacter(currentLine, true)

			if char == Tab {
				err = TabInIndentationError
				break
			}

			if nextIndex != NotFoundIndex {
				// returned to same level
				if nextIndex == baseIndentSpaces {
					// it is next element
					hasNext = true
					break
				}

				if nextIndex > baseIndentSpaces && currentLine[nextIndex] != CommentSymbol {
					// string has child
					err = StringHasChildError
					break
				}

				elementContent = append(elementContent, currentLine...)
			}
		}

		if d.Dictionary != nil {
			if _, exists := d.Dictionary[string(key)]; exists {
				err = DictionaryDuplicateKeyError
			}
		}

		if err == io.EOF {
			err = nil
		}

		if err != nil {
			return nil, hasNext, err
		}

		// char after line break
		firstChar, _ = readFirstMeaningfulCharacter(elementContent, true)

		if firstChar != EmptyChar {
			child := &Directive{
				IndentSize: d.IndentSize,
				Depth:      d.Depth + 1,
			}

			child.Type = DirectiveTypeString
			if initialLine[len(initialLine)-1] == LineBreak {
				child.String = string(initialLine[valueIndex : len(initialLine)-1])
			} else {
				child.String = string(initialLine[valueIndex:])
			}

			if d.Dictionary == nil {
				d.Dictionary = make(map[string]*Directive)
			}

			d.Dictionary[string(key)] = child
		} else {
			// empty case
			child := &Directive{
				Type:   DirectiveTypeString,
				String: "",
			}
			if d.Dictionary == nil {
				d.Dictionary = make(map[string]*Directive)
			}
			d.Dictionary[string(key)] = child
		}

		return currentLine, hasNext, err
	}

	var char byte
	var nextIndex int

	for err != io.EOF {
		if currentLine, err = buffer.ReadBytes(byte(LineBreak)); err != nil && err != io.EOF {
			break
		}

		char, nextIndex = readFirstMeaningfulCharacter(currentLine, true)

		if char == Tab {
			err = TabInIndentationError
			break
		}

		if nextIndex != NotFoundIndex {
			// returned to same level
			if nextIndex == baseIndentSpaces {
				// it is next element
				hasNext = true
				break
			}

			if char != EmptyChar && char != ListSymbol && char != TextSymbol && char != CommentSymbol {
				_, valueIndex := detectKeyBytes(currentLine)
				// sepIndex := getDictionarySeparatorIndex(line)
				if valueIndex == NotFoundIndex {
					// string has line break
					err = StringWithNewLineError
					break
				}
			}

			elementContent = append(elementContent, currentLine...)
		}
	}

	if firstChar == Tab {
		err = TabInIndentationError
	} else if d.Dictionary != nil {
		if _, exists := d.Dictionary[string(key)]; exists {
			err = DictionaryDuplicateKeyError
		}
	}

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return nil, hasNext, err
	}

	// char after line break
	firstChar, _ = readFirstMeaningfulCharacter(elementContent, true)

	if firstChar != EmptyChar {
		child := &Directive{
			IndentSize: d.IndentSize,
			Depth:      d.Depth + 1,
		}

		if marshalErr := child.Marshal(elementContent); marshalErr != nil {
			if marshalErr == EmptyDataError {
				child.Type = DirectiveTypeString
				child.String = ""
				marshalErr = nil
			} else {
				return nil, hasNext, marshalErr
			}
		}

		if d.Dictionary == nil {
			d.Dictionary = make(map[string]*Directive)
		}

		d.Dictionary[string(key)] = child
	} else {
		// empty case
		child := &Directive{
			Type:   DirectiveTypeString,
			String: "",
		}
		if d.Dictionary == nil {
			d.Dictionary = make(map[string]*Directive)
		}
		d.Dictionary[string(key)] = child
	}

	return currentLine, hasNext, err
}

func removeTrailingLineBreak(s *string) {
	l := len(*s)
	if l-1 >= 0 && (*s)[l-1] == LineBreak {
		*s = (*s)[:l-1]
	}
}

func sanitizeDictionaryKey(key *[]byte) error {
	quoted := false

	// remove sorrounding quotes
	{
		if len(*key) >= 2 {
			if ((*key)[0] == '"' && (*key)[len(*key)-1] == '"') || ((*key)[0] == '\'' && (*key)[len(*key)-1] == '\'') {
				quoted = true
				*key = (*key)[1 : len(*key)-1]
				quotedKey := *key
				// nested quote is allowed
				if len(*key) >= 2 {
					if ((*key)[0] == '"' && (*key)[len(*key)-1] == '"') || ((*key)[0] == '\'' && (*key)[len(*key)-1] == '\'') {
						quotedKey = (*key)[1 : len(*key)-1]
					}
				}
				for _, b := range quotedKey {
					if b == '"' || b == '\'' {
						return DictionaryKeyNestedQuotesError
					}
				}
			}
		}
	}

	// remove trailing space
	if !quoted {
		index := len(*key) - 1
		trimSize := 0

		for index >= 0 {
			if unicode.IsSpace(rune((*key)[index])) {
				trimSize++
			} else {
				break
			}
			index--
		}

		*key = (*key)[:len(*key)-trimSize]
	}

	return nil
}

func readFirstMeaningfulTwoCharacters(line []byte) ([]byte, int) {
	var index int

	for len(line) > index {
		char := line[index]

		if char != Space && char != LineBreak {
			if index+1 == len(line) {
				return []byte{char, EmptyChar}, index
			}

			return []byte{char, line[index+1]}, index
		}

		index++
	}

	return []byte{EmptyChar, EmptyChar}, NotFoundIndex
}

func readFirstMeaningfulCharacter(line []byte, skipLineBreak bool) (byte, int) {
	var char byte
	var index int

	for len(line) > index {
		char = line[index]

		if char != Space {
			if char == LineBreak {
				if !skipLineBreak {
					return char, index
				}
			} else {
				return char, index
			}
		}

		index++
	}

	return EmptyChar, NotFoundIndex
}

/**
 * Specifiaction:
 *   dictionary key is inspected with the following rules;
 *   1. line starts with a symbol of multi line text (>), list (-) and comment (#) with trailing space, are not inspected
 *   2. line with colon, the key separator symbol (:) with trailing space is subject to search dictionary key
 *   3. colon must be outside of quotes, otherwise it is not considered as key separator
 *   4. bytes from first meaningful character to the index before colon are considered as key
 *
 * This function expects line that is already considerd as dictionary.
 * So the test of first meaningful character is skipped.
 */
func detectKeyBytes(line []byte) ([]byte, int) {
	index := 0
	var char byte

	meaningfulIndex := NotFoundIndex
	quote := EmptyChar

	for index < len(line) {
		char = line[index]

		if (char == Quote || char == DoubleQuote) && char == quote {
			quote = EmptyChar
		}

		// 3.
		if meaningfulIndex == NotFoundIndex && !unicode.IsSpace(rune(char)) {
			meaningfulIndex = index
			if quote == EmptyChar && (char == Quote || char == DoubleQuote) {
				quote = char
			}
		}

		if DictionaryKeySeparator == char && quote == EmptyChar {
			if index < (len(line) - 1) {
				// 2.
				if unicode.IsSpace(rune(line[index+1])) {
					// 4.
					return line[meaningfulIndex:index], index + 2
				}
			} else {
				// 4.
				return line[meaningfulIndex:index], index + 1
			}
		}

		index++
	}

	return nil, NotFoundIndex
}

func getDictionarySeparatorIndex(line []byte) int {
	index := len(line) - 1
	var char byte

	for index >= 0 {
		char = line[index]

		if DictionaryKeySeparator == char {
			return index
		}

		index--
	}

	return NotFoundIndex
}
