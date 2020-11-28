package ntgo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type ValueType int

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeString
	ValueTypeText
	ValueTypeList
	ValueTypeDictionary
	ValueTypeComment
)

var (
	NextValueAppearedError            = errors.New("ntgo: next value appeared")
	DifferentTypesOnTheSameLevelError = errors.New("ntgo: can not place different types of entities on the same level")
	DictionaryKeyWithUnpairedQuotes   = errors.New("ntgo: quoted dictionary key can not contain unpaired quotes")
	EmptyDataError                    = errors.New("ntgo: data can not be empty")
	RootLevelHasIndentError           = errors.New("ntgo: root level must not be indented")
	TabInIndentationError             = errors.New("ntgo: indent can not contain tab")
	RootStringError                   = errors.New("ntgo: no string allowed on root level")
	StringHasChildError               = errors.New("ntgo: string type can not have child")
	TextHasChildError                 = errors.New("ntgo: text type can not have child")
	DifferentLevelOnSameChildError    = errors.New("ntgo: child elements have dirfferent levels")
	StringWithNewLineError            = errors.New("ntgo: string type can not have line break")
	DictionaryDuplicateKeyError       = errors.New("ntgo: dictionary type can not have the same key")
)

func (t ValueType) String() string {
	switch t {
	case ValueTypeUnknown:
		return "unknown"
	case ValueTypeString:
		return "string"
	case ValueTypeText:
		return "text"
	case ValueTypeList:
		return "list"
	case ValueTypeDictionary:
		return "dictionary"
	case ValueTypeComment:
		return "comment"
	}
	return ""
}

type MultilineStrings []string

func (t MultilineStrings) String() string {
	return strings.Join(t, "")
}

type Value struct {
	Type ValueType

	String     string
	Text       MultilineStrings
	List       []*Value
	Dictionary map[string]*Value

	IndentSize int
	Depth      int
}

func (v *Value) ToNestedText() string {
	str := ""

	if v.IndentSize <= 0 {
		// default size
		v.IndentSize = UnmarshalDefaultIndentSize
	}

	baseIndent := fmt.Sprintf("%*s", v.IndentSize*v.Depth, "")

	switch v.Type {
	case ValueTypeString:
		str = v.String
	case ValueTypeText:
		for i := 0; i < len(v.Text); i++ {
			str = fmt.Sprintf("%s%s> %s", str, baseIndent, v.Text[i])
		}
	case ValueTypeList:
		for i := 0; i < len(v.List); i++ {
			// TODO: user prefered line break code
			dataLn := string(LF)

			child := v.List[i]
			if child.Type == ValueTypeString {
				dataLn = string(Space)
			}

			// TODO: linear recursion
			str = fmt.Sprintf("%s%s-%s%s\n", str, baseIndent, dataLn, child.ToNestedText())
		}
	case ValueTypeDictionary:
		it := 0
		for k, v := range v.Dictionary {
			dataLn := string(LF)

			if v.Type == ValueTypeString {
				dataLn = string(Space)
			}

			str = fmt.Sprintf("%s%s%s:%s%s\n", str, baseIndent, k, dataLn, v.ToNestedText())

			it++
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
		// CRLF can be ignored under parsing data structure
		// It only should be considered under parsing multi line text
		if b == CR || b == LF {
			break
		}
	}
	return
}

func (v *Value) Parse(content []byte) (err error) {
	v.Type = ValueTypeUnknown

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

		var valueType ValueType
		if valueType, index, err = detectValueType(currentLine); err != nil {
			break
		}

		if v.Depth == 0 && index > 0 {
			err = RootLevelHasIndentError
			break
		}

		switch valueType {
		case ValueTypeUnknown, ValueTypeComment:
			{
				loadedNextLine = false
				err = nil
			}
		case ValueTypeText:
			currentLine, loadedNextLine, err = v.readTextValue(index, currentLine, buffer)
		case ValueTypeList:
			currentLine, loadedNextLine, err = v.readListValue(index, currentLine, buffer)
		case ValueTypeDictionary, ValueTypeString:
			currentLine, loadedNextLine, err = v.readDictionaryValue(index, currentLine, buffer)
		}

		if err != nil {
			break
		}
	}

	if err == nil && v.Type == ValueTypeUnknown {
		err = EmptyDataError
	}

	return
}

func detectValueType(line []byte) (ValueType, int, error) {
	valueType := ValueTypeUnknown
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
		valueType = ValueTypeComment
	case Tab:
		return ValueTypeUnknown, index, TabInIndentationError
	case TextSymbol:
		{
			switch chars[1] {
			case Space, CR, LF, EmptyChar:
				valueType = ValueTypeText
			default:
				valueType = ValueTypeString
			}
		}
	case ListSymbol:
		{
			switch chars[1] {
			case Space, CR, LF, EmptyChar:
				valueType = ValueTypeList
			default:
				valueType = ValueTypeString
			}
		}
	default:
		//
		_, keyIndex := detectKeyBytes(line)
		if keyIndex == NotFoundIndex {
			valueType = ValueTypeString
		} else {
			valueType = ValueTypeDictionary
		}
	}

	return valueType, index, nil
}

func (v *Value) readTextValue(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	if v.Type != ValueTypeUnknown {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}

	v.Type = ValueTypeText

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
				v.Text = append(v.Text, "")
			} else if currentLine[nextIndex+1] == CR || currentLine[nextIndex+1] == LF {
				// text symbol with no space
				v.Text = append(v.Text, string(currentLine[nextIndex+1]))
			} else {
				// after text symbol(>) and space
				v.Text = append(v.Text, string(currentLine[nextIndex+2:]))
			}
		}

		if err == io.EOF {
			if len(v.Text) >= 1 {
				removeStringTrailingLineBreaks(&v.Text[len(v.Text)-1])
			}
			return nil, hasNext, nil
		}

		if currentLine, err = readLine(buffer); err != nil && err != io.EOF {
			return nil, hasNext, err
		}

		// CRLF
		if len(currentLine) == 1 && currentLine[0] == LF {
			if len(v.Text) > 0 {
				lastLine := v.Text[len(v.Text)-1]
				if lastLine[len(lastLine)-1] == CR {
					v.Text[len(v.Text)-1] += string(LF)
				}
			}
		}
	}

	if hasNext {
		if len(v.Text) >= 1 {
			removeStringTrailingLineBreaks(&v.Text[len(v.Text)-1])
		}
	}

	return currentLine, hasNext, nil
}

func (v *Value) readListValue(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	if v.Type != ValueTypeUnknown && v.Type != ValueTypeList {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}

	v.Type = ValueTypeList

	currentLine := initialLine
	elementContent := currentLine[baseIndentSpaces+1:]

	firstChar, _ := readFirstMeaningfulCharacter(elementContent, true)

	var child *Value

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

		child = &Value{Type: ValueTypeString}

		if firstChar, _ = readFirstMeaningfulCharacter(elementContent, true); firstChar != EmptyChar {
			child.IndentSize = v.IndentSize
			child.Depth = v.Depth + 1

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
			// TODO: reading twice
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

		child = &Value{
			IndentSize: v.IndentSize,
			Depth:      v.Depth + 1,
		}

		// Parse child
		// TODO: elementContent internally converted to bytes.Buffer, inpsect its performance cost
		if err := child.Parse(elementContent); err != nil {
			if err != EmptyDataError {
				return nil, hasNext, err
			}
			// treat empty data as empty string
			child.Type = ValueTypeString
			child.String = ""
		}
	}

	v.List = append(v.List, child)

	return currentLine, hasNext, nil
}

func (v *Value) readDictionaryValue(baseIndentSpaces int, initialLine []byte, buffer *bytes.Buffer) ([]byte, bool, error) {
	hasNext := false
	var err error

	// dictionary
	if v.Type != ValueTypeUnknown && v.Type != ValueTypeDictionary {
		return nil, hasNext, DifferentTypesOnTheSameLevelError
	}

	key, valueIndex := detectKeyBytes(initialLine)

	// unexpected string
	if key == nil && valueIndex == NotFoundIndex {
		return nil, hasNext, RootStringError
	}
	if err = sanitizeDictionaryKey(&key); err != nil {
		return nil, hasNext, err
	}

	if v.Dictionary != nil {
		if _, exists := v.Dictionary[string(key)]; exists {
			return nil, hasNext, DictionaryDuplicateKeyError
		}
	}

	v.Type = ValueTypeDictionary

	currentLine := initialLine
	elementContent := currentLine[valueIndex:]

	firstChar, _ := readFirstMeaningfulCharacter(elementContent, true)

	var child *Value

	// child is string
	if len(initialLine) > valueIndex {
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

			if nextIndex == NotFoundIndex || char == CommentSymbol {
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

		child = &Value{Type: ValueTypeString}

		// char after line break
		if firstChar, _ = readFirstMeaningfulCharacter(elementContent, true); firstChar != EmptyChar {
			child.IndentSize = v.IndentSize
			child.Depth = v.Depth + 1

			if initialLine[len(initialLine)-1] == CR || initialLine[len(initialLine)-1] == LF {
				child.String = string(initialLine[valueIndex : len(initialLine)-1])
			} else {
				child.String = string(initialLine[valueIndex:])
			}
		}
	} else {
		levels := []int{}
		for eof := false; !eof; {
			lastLine := currentLine
			if currentLine, err = readLine(buffer); err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, hasNext, err
			}

			// CRLF
			if len(currentLine) == 1 && currentLine[0] == LF {
				if len(lastLine) > 0 && lastLine[len(lastLine)-1] == CR {
					elementContent = append(elementContent, currentLine[0])
				}
			}

			char, nextIndex := readFirstMeaningfulCharacter(currentLine, true)

			if char == Tab {
				return nil, hasNext, TabInIndentationError
			}

			if nextIndex == NotFoundIndex {
				continue
			}

			if char == CommentSymbol {
				continue
			}

			// returned to same level
			if nextIndex == baseIndentSpaces {
				// it is next element
				hasNext = true
				break
			}

			// inspect indent level validity
			if len(levels) == 0 {
				levels = append(levels, nextIndex)
			} else {
				if levels[len(levels)-1] < nextIndex {
					levels = append(levels, nextIndex)
				} else {
					found := false
					for i, l := range levels {
						if l == nextIndex {
							levels = levels[:i+1]
							found = true
							break
						}
					}
					if !found {
						return nil, hasNext, DifferentLevelOnSameChildError
					}
				}
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

		child = &Value{Depth: v.Depth + 1}

		// empty case
		if firstChar == EmptyChar {
			child.Type = ValueTypeString
			child.String = ""
		} else {
			child.IndentSize = v.IndentSize

			if err = child.Parse(elementContent); err == EmptyDataError {
				child.Type = ValueTypeString
				child.String = ""
			} else if err != nil {
				return nil, hasNext, err
			}
		}
	}

	if v.Dictionary == nil {
		v.Dictionary = make(map[string]*Value)
	}

	v.Dictionary[string(key)] = child

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
	keyLen := len(*key)

	index := len(*key) - 1
	trimSize := 0

	// remove trailing space
	for index >= 0 {
		if unicode.IsSpace(rune((*key)[index])) {
			trimSize++
		} else {
			break
		}
		index--
	}

	*key = (*key)[:len(*key)-trimSize]
	keyLen = len(*key)

	// remove surrounding quotes
	if keyLen >= 2 {
		if ((*key)[0] == DoubleQuote && (*key)[keyLen-1] == DoubleQuote) || ((*key)[0] == Quote && (*key)[keyLen-1] == Quote) {
			*key = (*key)[1 : keyLen-1]
			keyLen = len(*key)
		}
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
 *   2. line with the dictionary key delimiter symbol (':') is subject to search dictionary key
 *   3. treat delimiter inside of quotes as a part of key
 *   4. quotes detection is prior to delimiter detection
 *   5. bytes from first meaningful character to the index before delimiter are considered as key
 *
 * This function expects line that is already considerd as dictionary.
 * So the test of first meaningful character is skipped.
 */
func detectKeyBytes(line []byte) ([]byte, int) {
	var char byte

	meaningfulIndex := NotFoundIndex
	quoteClosingIndex := NotFoundIndex
	delimiterBeginIndex := NotFoundIndex
	delimiterEndIndex := NotFoundIndex
	quote := EmptyChar

	// 4.
	for index := 0; index < len(line); index++ {
		char = line[index]

		if quote != EmptyChar && char == quote {
			quoteClosingIndex = index
		}
		// 3.
		if meaningfulIndex == NotFoundIndex && !unicode.IsSpace(rune(char)) {
			meaningfulIndex = index
			if char == Quote || char == DoubleQuote {
				quote = char
			}
		}
		// 2.
		if char == DictionaryKeySeparator {
			if index > quoteClosingIndex && (delimiterBeginIndex == NotFoundIndex || delimiterBeginIndex < quoteClosingIndex) {
				// ':' with line break
				if index >= (len(line) - 1) {
					delimiterBeginIndex = index
					delimiterEndIndex = index + 1
				} else {
					// ':' with space
					if unicode.IsSpace(rune(line[index+1])) {
						delimiterBeginIndex = index
						delimiterEndIndex = index + 2
					}
				}
			}
		}
	}

	if delimiterBeginIndex == NotFoundIndex {
		return nil, NotFoundIndex
	}
	// 5.
	return line[meaningfulIndex:delimiterBeginIndex], delimiterEndIndex
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
