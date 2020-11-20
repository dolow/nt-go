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
	Text       MultiLineText
	List       []*Directive
	Dictionary map[string]*Directive

	IndentSize int
	Depth      int
}

func (d *Directive) ToString() string {
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
				// TODO: user prefered line breal code
				dataLn := string(LF)
				tailLn := string(LF)
				if i == len(d.List)-1 {
					tailLn = ""
				}

				child := d.List[i]
				if child.Type == DirectiveTypeString {
					dataLn = ""
				}

				str = fmt.Sprintf("%s%s- %s%s%s", str, baseIndent, dataLn, child.ToString(), tailLn)
			}
		}
	case DirectiveTypeDictionary:
		{
			it := 0
			for k, v := range d.Dictionary {
				dataLn := string(LF)
				tailLn := string(LF)
				if it == len(d.Dictionary)-1 {
					tailLn = ""
				}

				if v.Type == DirectiveTypeString {
					dataLn = ""
				}

				str = fmt.Sprintf("%s%s%s: %s%s%s", str, baseIndent, k, dataLn, v.ToString(), tailLn)

				it++
			}
		}
	}

	return str
}

func readLine(buffer *bytes.Buffer) (line []byte, err error) {
	var b byte
	for err != io.EOF {
		if b, err = buffer.ReadByte(); err != nil && err != io.EOF {
			break
		} else if b == EmptyChar {
			break
		}
		line = append(line, b)
		// skip CRLF case, it will treated as empty line
		if b == CR || b == LF {
			break
		}
	}
	return
}

func (d *Directive) Parse(content []byte) (err error) {
	d.Type = DirectiveTypeUnknown

	removeBytesTrailingLineBreaks(&content)
	buffer := bytes.NewBuffer(content)

	var currentLine []byte
	var index int
	loadedNextLine := false

	for eof := false; !eof; {
		if !loadedNextLine {
			if currentLine, err = readLine(buffer); err == io.EOF {
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
		case DirectiveTypeDictionary, DirectiveTypeString:
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

func detectDirectiveType(line []byte) (DirectiveType, int, error) {
	directiveType := DirectiveTypeUnknown
	index := 0

	chars := []byte{EmptyChar, EmptyChar}
	for ; len(line) > index; index++ {
		char := line[index]

		if char != Space && char != CR && char != LF {
			chars[0] = char
			if index+1 == len(line) {
				chars[1] = EmptyChar
			} else {
				chars[1] = line[index+1]
			}
			break
		}
	}

	switch chars[0] {
	case EmptyChar:
		index = NotFoundIndex
	case CommentSymbol:
		directiveType = DirectiveTypeComment
	case Tab:
		return DirectiveTypeUnknown, index, TabInIndentationError
	case TextSymbol:
		{
			switch chars[1] {
			case Space, CR, LF, EmptyChar:
				directiveType = DirectiveTypeText
			default:
				directiveType = DirectiveTypeString
			}
		}
	case ListSymbol:
		{
			switch chars[1] {
			case Space, CR, LF, EmptyChar:
				directiveType = DirectiveTypeList
			default:
				directiveType = DirectiveTypeString
			}
		}
	default:
		//
		_, keyIndex := detectKeyBytes(line)
		if keyIndex == NotFoundIndex {
			directiveType = DirectiveTypeString
		} else {
			directiveType = DirectiveTypeDictionary
		}
	}

	return directiveType, index, nil
}

func (d *Directive) readTextDirective(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	if d.Type != DirectiveTypeUnknown {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}

	d.Type = DirectiveTypeText

	var err error
	currentLine := initialLine

	for {
		char, nextIndex := readFirstMeaningfulCharacter(currentLine, false)

		if char != CommentSymbol && char != CR && char != LF && nextIndex != NotFoundIndex {
			// validate
			if char == TextSymbol {
				if nextIndex != baseIndentSpaces {
					return nil, hasNext, DifferentLevelOnSameChildError
				}
			} else {
				if nextIndex > baseIndentSpaces {
					return nil, hasNext, TextHasChildError
				}
				if nextIndex == baseIndentSpaces {
					return nil, hasNext, DifferentTypesOnTheSameLevelError
				}
				if nextIndex < baseIndentSpaces {
					hasNext = true
					break
				}
			}

			// append text
			if len(currentLine) <= nextIndex+1 {
				// text ends with symbol
				d.Text = append(d.Text, "")
			} else if currentLine[nextIndex+1] == CR || currentLine[nextIndex+1] == LF {
				// text symbol with no space
				// TODO: CRLF case
				d.Text = append(d.Text, string(currentLine[nextIndex+1]))
			} else {
				// after text symbol(>) and space
				d.Text = append(d.Text, string(currentLine[nextIndex+2:]))
			}
		}

		if err == io.EOF {
			if len(d.Text) >= 1 {
				removeStringTrailingLineBreaks(&d.Text[len(d.Text)-1])
			}
			return nil, hasNext, nil
		}

		if currentLine, err = readLine(buffer); err != nil && err != io.EOF {
			return nil, hasNext, err
		}
	}

	if hasNext {
		if len(d.Text) >= 1 {
			removeStringTrailingLineBreaks(&d.Text[len(d.Text)-1])
		}
	}

	return currentLine, hasNext, nil
}

func (d *Directive) readListDirective(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	if d.Type != DirectiveTypeUnknown && d.Type != DirectiveTypeList {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}

	d.Type = DirectiveTypeList

	currentLine := initialLine
	elementContent := currentLine[baseIndentSpaces+1:]

	firstChar, _ := readFirstMeaningfulCharacter(elementContent, true)

	var child *Directive

	// string case
	if firstChar != EmptyChar {
		for eof := false; !eof; {
			var err error
			if currentLine, err = readLine(buffer); err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, hasNext, err
			}

			char, nextIndex := readFirstMeaningfulCharacter(currentLine, true)
			if char == Tab {
				return currentLine, hasNext, TabInIndentationError
			}

			if nextIndex == NotFoundIndex {
				continue
			}

			if nextIndex == baseIndentSpaces {
				hasNext = true
				break
			}

			// validate
			if currentLine[nextIndex] != CommentSymbol {
				if nextIndex > baseIndentSpaces {
					return currentLine, hasNext, StringHasChildError
				}
				// parent should not be contained
				if nextIndex < baseIndentSpaces {
					return currentLine, hasNext, DifferentLevelOnSameChildError
				}
			}

			elementContent = append(elementContent, currentLine...)
		}

		child = &Directive{Type: DirectiveTypeString}

		if firstChar, _ = readFirstMeaningfulCharacter(elementContent, true); firstChar != EmptyChar {
			child.IndentSize = d.IndentSize
			child.Depth = d.Depth + 1

			if initialLine[len(initialLine)-1] == CR || initialLine[len(initialLine)-1] == LF {
				child.String = string(initialLine[baseIndentSpaces+2 : len(initialLine)-1])
			} else {
				child.String = string(initialLine[baseIndentSpaces+2:])
			}
		}
	} else {
		// collect child content lines
		for eof := false; !eof; {
			var err error
			if currentLine, err = readLine(buffer); err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, hasNext, err
			}

			char, newIndex := readFirstMeaningfulCharacter(currentLine, true)
			if char == Tab {
				return nil, hasNext, TabInIndentationError
			}

			if newIndex == baseIndentSpaces {
				hasNext = true
				break
			}

			elementContent = append(elementContent, currentLine...)
		}

		child = &Directive{
			IndentSize: d.IndentSize,
			Depth:      d.Depth + 1,
		}

		// Parse child
		// TODO: elementContent internally converted to bytes.Buffer, inpsect its performance cost
		if err := child.Parse(elementContent); err != nil {
			if err != EmptyDataError {
				return nil, hasNext, err
			}
			// treat empty data as empty string
			child.Type = DirectiveTypeString
			child.String = ""
		}
	}

	d.List = append(d.List, child)

	return currentLine, hasNext, nil
}

func (d *Directive) readDictionaryDirective(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	var err error

	// dictionary
	if d.Type != DirectiveTypeUnknown && d.Type != DirectiveTypeDictionary {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}
	if d.Depth == 0 && baseIndentSpaces > 0 {
		return nil, hasNext, RootLevelHasIndentError
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
			return nil, hasNext, DictionaryDuplicateKeyError
		}
	}

	d.Type = DirectiveTypeDictionary

	currentLine := initialLine
	elementContent := currentLine[valueIndex:]

	firstChar, _ := readFirstMeaningfulCharacter(elementContent, true)

	var child *Directive

	// child is string
	if firstChar != EmptyChar {
		for eof := false; !eof; {
			if currentLine, err = readLine(buffer); err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, hasNext, err
			}

			char, nextIndex := readFirstMeaningfulCharacter(currentLine, true)

			if char == Tab {
				return nil, hasNext, TabInIndentationError
			}

			if nextIndex == NotFoundIndex {
				continue
			}

			// returned to same level
			if nextIndex == baseIndentSpaces {
				// it is next element
				hasNext = true
				break
			}

			if currentLine[nextIndex] != CommentSymbol {
				if nextIndex > baseIndentSpaces {
					return nil, hasNext, StringHasChildError
				}
				// parent should not be contained
				if nextIndex < baseIndentSpaces {
					return nil, hasNext, DifferentLevelOnSameChildError
				}
			}

			elementContent = append(elementContent, currentLine...)
		}

		child = &Directive{Type: DirectiveTypeString}

		// char after line break
		if firstChar, _ = readFirstMeaningfulCharacter(elementContent, true); firstChar != EmptyChar {
			child.IndentSize = d.IndentSize
			child.Depth = d.Depth + 1

			if initialLine[len(initialLine)-1] == CR || initialLine[len(initialLine)-1] == LF {
				child.String = string(initialLine[valueIndex : len(initialLine)-1])
			} else {
				child.String = string(initialLine[valueIndex:])
			}
		}
	} else {
		for eof := false; !eof; {
			if currentLine, err = readLine(buffer); err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, hasNext, err
			}

			char, nextIndex := readFirstMeaningfulCharacter(currentLine, true)

			if char == Tab {
				return nil, hasNext, TabInIndentationError
			}

			if nextIndex == NotFoundIndex {
				continue
			}
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
					return nil, hasNext, StringWithNewLineError
				}
			}

			elementContent = append(elementContent, currentLine...)
		}

		// char after line break
		firstChar, _ = readFirstMeaningfulCharacter(elementContent, true)

		child = &Directive{Depth: d.Depth + 1}

		// empty case
		if firstChar == EmptyChar {
			child.Type = DirectiveTypeString
			child.String = ""
		} else {
			child.IndentSize = d.IndentSize

			if err = child.Parse(elementContent); err == EmptyDataError {
				child.Type = DirectiveTypeString
				child.String = ""
			} else if err != nil {
				return nil, hasNext, err
			}
		}
	}

	if d.Dictionary == nil {
		d.Dictionary = make(map[string]*Directive)
	}

	d.Dictionary[string(key)] = child

	return currentLine, hasNext, nil
}

func removeStringTrailingLineBreaks(s *string) {
	l := len(*s)
	if l > 0 && ((*s)[l-1] == CR || (*s)[l-1] == LF) {
		*s = (*s)[:l-1]
	}
}

func removeBytesTrailingLineBreaks(b *[]byte) {
	l := len(*b)
	if l > 0 && ((*b)[l-1] == CR || (*b)[l-1] == LF) {
		*b = (*b)[:l-1]
	}
}

func sanitizeDictionaryKey(key *[]byte) error {
	quoted := false

	// remove sorrounding quotes
	if len(*key) >= 2 {
		if ((*key)[0] == DoubleQuote && (*key)[len(*key)-1] == DoubleQuote) || ((*key)[0] == Quote && (*key)[len(*key)-1] == Quote) {
			quoted = true
			*key = (*key)[1 : len(*key)-1]
			quotedKey := *key
			// nested quote is allowed
			if len(*key) >= 2 {
				if ((*key)[0] == DoubleQuote && (*key)[len(*key)-1] == DoubleQuote) || ((*key)[0] == Quote && (*key)[len(*key)-1] == Quote) {
					quotedKey = (*key)[1 : len(*key)-1]
				}
			}
			for _, b := range quotedKey {
				if b == DoubleQuote || b == Quote {
					return DictionaryKeyNestedQuotesError
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

func readFirstMeaningfulCharacter(line []byte, skipLineBreak bool) (byte, int) {
	var char byte
	var index int

	for len(line) > index {
		char = line[index]

		switch char {
		case Space:
		case CR, LF:
			{
				if !skipLineBreak {
					return char, index
				}
			}
		default:
			return char, index
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
