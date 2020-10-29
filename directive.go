package nestedtext

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

	EmptyChar byte = 0x00
	NotFoundIndex int = -1
)

var NextDirectiveAppearedError = errors.New("nestedtext: next directive appeared")
var DifferentTypesOnTheSameLevelError = errors.New("nestedtext: can not place different types of entities on the same level")
var DictionaryKeyNestedQuotesError = errors.New("nestedtext: quoted dictionary key can not contain any quotes")
var EmptyDataError = errors.New("nestedtext: data can not be empty")
var RootLevelHasIndentError = errors.New("nestedtext: root level must not be indented")
var TabInIndentationError = errors.New("nestedtext: indent can not contain tab")
var RootStringError = errors.New("nestedtext: no string allowed on root level")
var StringHasChildError = errors.New("nestedtext: string type can not have child")
var DifferentLevelOnSameChildError = errors.New("nestedtext: child elements have dirfferent leves")
var StringWithNewLineError = errors.New("nestedtext: string type can not have line break")
var DictionaryDuplicateKeyError = errors.New("nestedtext: dictionary type can not have the same key")

type Directive struct {
	Type DirectiveType

	String     string
	Text       []string
	List       []*Directive
	Dictionary map[string]*Directive

	IndentSize int
	Depth      int
}

type DirectiveMarshalError struct {
	error

	Line   int
	Offset int
}


func (d *Directive) Unmarshal() string {
	str := ""

	baseIndent := fmt.Sprintf("%*s", d.IndentSize * d.Depth, "")
	switch d.Type {
	case DirectiveTypeString: {
		str = d.String
	}
	case DirectiveTypeText: {
		for i := 0; i < len(d.Text); i++ {
			str = fmt.Sprintf("%s%s> %s", str, baseIndent, d.Text[i])
		}
	}
	case DirectiveTypeList: {
		for i := 0; i < len(d.List); i++ {
			dataLn := "\n"
			tailLn := "\n"
			if i == len(d.List) - 1 {
				tailLn = ""
			}

			child := d.List[i]
			if child.Type == DirectiveTypeString {
				dataLn = ""
			}

			str = fmt.Sprintf("%s%s- %s%s%s", str, baseIndent, dataLn, child.Unmarshal(), tailLn)
		}
	}
	case DirectiveTypeDictionary: {
		it := 0
		for k, v := range d.Dictionary {
			dataLn := "\n"
			tailLn := "\n"
			if it == len(d.Dictionary) - 1 {
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

func (d *Directive) Marshal(content []byte) *DirectiveMarshalError {
	var marshalErr *DirectiveMarshalError = nil
	d.Type = DirectiveTypeUnknown

	// remove trailing line breaks
	for len(content) > 0 && content[len(content) - 1] == '\n' {
		content = content[:len(content) - 1]
	}

	buffer := bytes.NewBuffer(content)

	var line []byte
	var index, lastIndex int
	var readBytesErr error

	ReadLineLoop: for {
		eof := false

		if readBytesErr != NextDirectiveAppearedError {
			line, readBytesErr = buffer.ReadBytes(byte('\n'))

			eof = (readBytesErr == io.EOF)
			if !eof && readBytesErr != nil {
				marshalErr = &DirectiveMarshalError{ error: readBytesErr }
				break
			}
		}

		firstMeaningfulChar, newIndex := readFirstMeaningfulCharacter(line, true)
		if firstMeaningfulChar == Tab {
			marshalErr = &DirectiveMarshalError{
				error: TabInIndentationError,
			}
			break
		}

		if firstMeaningfulChar != CommentSymbol {
			lastIndex = index
			index = newIndex
		}

		directiveType := DirectiveTypeUnknown

		switch firstMeaningfulChar {
		case EmptyChar:
		case CommentSymbol:
		case TextSymbol: {
			if len(line) <= index + 1 {
				directiveType = DirectiveTypeText
			} else {
				nextChar := line[index + 1]
				if nextChar == ' ' || nextChar == '\n' {
					directiveType = DirectiveTypeText
				} else {
					directiveType = DirectiveTypeString
				}
			}
		}
		case ListSymbol: {
			if len(line) <= index + 1 {
				directiveType = DirectiveTypeList
			} else {
				nextChar := line[index + 1]
				if nextChar == ' ' || nextChar == '\n' {
					directiveType = DirectiveTypeList
				} else {
					directiveType = DirectiveTypeString
				}
			}
		}
		default: directiveType = DirectiveTypeString
		}

		switch directiveType {
		case DirectiveTypeUnknown:
		case DirectiveTypeText:  // multi line text
			{
				if d.Type == DirectiveTypeText {
					if index != lastIndex {
						marshalErr = &DirectiveMarshalError{
							error: DifferentLevelOnSameChildError,
						}
						break ReadLineLoop
					}
				} else if d.Type != DirectiveTypeUnknown {
					marshalErr = &DirectiveMarshalError{
						error: DifferentTypesOnTheSameLevelError,
					}
					break ReadLineLoop
				}

				d.Type = DirectiveTypeText

				TextChildReadLineLoop: for {
					char, newIndex := readFirstMeaningfulCharacter(line, false)
					trailingBlankLine := false

					if char != CommentSymbol && newIndex != NotFoundIndex {
						if newIndex > index {
							// deeper
							marshalErr = &DirectiveMarshalError{
								error: DifferentLevelOnSameChildError,
							}
							break TextChildReadLineLoop
						} else if newIndex < index {
							// shallower
							readBytesErr = NextDirectiveAppearedError
							break TextChildReadLineLoop
						}

						if char != TextSymbol {
							marshalErr = &DirectiveMarshalError{
								error: DifferentLevelOnSameChildError,
							}
							break TextChildReadLineLoop
						}

						_, contentIndex := readFirstMeaningfulCharacter(line[newIndex + 1:], false)

						if contentIndex == NotFoundIndex {
							d.Text = append(d.Text, "")
						} else {
							d.Text = append(d.Text, string(line[newIndex + 1 + contentIndex:]))
						}
					} else {
						trailingBlankLine = true
					}

					if readBytesErr != nil {
						// if directive contains trailing blank line, last line of text directive may have trailing line break
						if trailingBlankLine && len(d.Text) >= 1 {
							lastLine := d.Text[len(d.Text) - 1]
							if len(lastLine) >= 1 && lastLine[len(lastLine) - 1] == LineBreak {
								d.Text[len(d.Text) - 1] = lastLine[:len(lastLine) - 1]
							}
						}

						if readBytesErr == io.EOF {
							break TextChildReadLineLoop
						}
						marshalErr = &DirectiveMarshalError{
							error: readBytesErr,
						}
						break TextChildReadLineLoop
					}

					line, readBytesErr = buffer.ReadBytes(byte('\n'))
				}
			}
		case DirectiveTypeList: // list
			{
				if d.Type != DirectiveTypeUnknown && d.Type != DirectiveTypeList {
					marshalErr = &DirectiveMarshalError{
						error: DifferentTypesOnTheSameLevelError,
					}
					break ReadLineLoop
				}

				d.Type = DirectiveTypeList
				elementContent := line[index+1:]

				// detect string
				elementContentChar, _ := readFirstMeaningfulCharacter(elementContent, true)
				childIsString := elementContentChar != EmptyChar

				firstLine := line

				// TODO: almost same as dictionary
				ListChildReadLineLoop: for {
					line, readBytesErr = buffer.ReadBytes(byte('\n'))

					childEof := (readBytesErr == io.EOF)
					if !childEof && readBytesErr != nil {
						marshalErr = &DirectiveMarshalError{
							error: readBytesErr,
						}
						break ListChildReadLineLoop
					}

					char, nextIndex := readFirstMeaningfulCharacter(line, true)

					if char == Tab {
						marshalErr = &DirectiveMarshalError{
							error: TabInIndentationError,
						}
						break ListChildReadLineLoop
					}

					if nextIndex == index {
						// it is next element
						readBytesErr = NextDirectiveAppearedError
						break ListChildReadLineLoop
					} else {
						var err error

						if childIsString && nextIndex != NotFoundIndex {
							if nextIndex > index && line[nextIndex] != CommentSymbol {
								err = StringHasChildError
							} else if nextIndex < index {
								err = DifferentLevelOnSameChildError
							}
						}
						if err != nil {
							marshalErr = &DirectiveMarshalError{
								error: err,
							}
							break ReadLineLoop
						}
					}

					elementContent = append(elementContent, line...)

					if childEof {
						break ListChildReadLineLoop
					}
				}

				child := &Directive{
					IndentSize: d.IndentSize,
					Depth: d.Depth + 1,
				}

				if childIsString {
					child.Type = DirectiveTypeString
					if firstLine[len(firstLine) - 1] == '\n' {
						child.String = string(firstLine[index + 2:len(firstLine) - 1])
					} else {
						child.String = string(firstLine[index + 2:])
					}
				} else {
					// TODO: elementContent internally converted to bytes.Buufer, inpsect its performance cost
					if marshalErr = child.Marshal(elementContent); marshalErr != nil {
						if marshalErr.error == EmptyDataError {
							child.Type = DirectiveTypeString
							child.String = ""
							marshalErr = nil
						} else {
							break ReadLineLoop
						}
					}
				}

				d.List = append(d.List, child)
			}
		default: // string or dictionary
			{
				dirtyKey, valueIndex := detectKeyBytes(line)
				
				if dirtyKey == nil && valueIndex == NotFoundIndex {
					// string
					if d.Type != DirectiveTypeUnknown {
						marshalErr = &DirectiveMarshalError{
							error: DifferentTypesOnTheSameLevelError,
						}
						break ReadLineLoop
					}

					if d.Depth == 0 {
						marshalErr = &DirectiveMarshalError{
							error: RootStringError,
						}
						break ReadLineLoop	
					}

					d.Type = DirectiveTypeString
					// remove trailing line break
					if line[len(line)-1] == '\n' {
						d.String = string(line[index : len(line)-1])
					} else {
						d.String = string(line[index:])
					}
					break ReadLineLoop
				}

				// dictionary
				if d.Type != DirectiveTypeUnknown && d.Type != DirectiveTypeDictionary {
					marshalErr = &DirectiveMarshalError{
						error: DifferentTypesOnTheSameLevelError,
					}
					break ReadLineLoop
				}

				if d.Depth == 0 && index > 0 {
					marshalErr = &DirectiveMarshalError{
						error: RootLevelHasIndentError,
					}
					break ReadLineLoop
				}

				d.Type = DirectiveTypeDictionary
				key, err := dictionaryKeySanitize(dirtyKey)

				if err != nil {
					marshalErr = &DirectiveMarshalError{
						error: err,
					}
					break ReadLineLoop
				}

				firstLine := line
				elementContent := line[valueIndex:]

				firstChar, _ := readFirstMeaningfulCharacter(elementContent, true)
				childIsString := firstChar != EmptyChar

				var char byte
				var nextIndex int
				DictionaryChildReadLineLoop: for {
					line, readBytesErr = buffer.ReadBytes(byte('\n'))

					childEof := (readBytesErr == io.EOF)
					if !childEof && readBytesErr != nil {
						marshalErr = &DirectiveMarshalError{
							error: readBytesErr,
						}
						break DictionaryChildReadLineLoop
					}

					char, nextIndex = readFirstMeaningfulCharacter(line, true)

					if char == Tab {
						marshalErr = &DirectiveMarshalError{
							error: TabInIndentationError,
						}
						break DictionaryChildReadLineLoop
					}

					if nextIndex != NotFoundIndex {
						// returned to same level
						if nextIndex == index {
							if char == ListSymbol || char == TextSymbol {
								// TODO: irregular case
							}
							// it is next element
							readBytesErr = NextDirectiveAppearedError
							break DictionaryChildReadLineLoop
						}

						if childIsString {
							if nextIndex > index && line[nextIndex] != CommentSymbol {
								// string has child
								marshalErr = &DirectiveMarshalError{
									error: StringHasChildError,
								}
								break ReadLineLoop
							}
						} else {
							if char != EmptyChar && char != ListSymbol && char != TextSymbol && char != CommentSymbol {
								_, valueIndex := detectKeyBytes(line)
								// sepIndex := getDictionarySeparatorIndex(line)
								if valueIndex == NotFoundIndex {
									// string has line break
									marshalErr = &DirectiveMarshalError{
										error: StringWithNewLineError,
									}
									break ReadLineLoop
								}
							}
						}

						elementContent = append(elementContent, line...)
					}

					if childEof {
						break DictionaryChildReadLineLoop
					}
				}

				if firstChar == Tab && !childIsString {
					marshalErr = &DirectiveMarshalError{
						error: TabInIndentationError,
					}
					break ReadLineLoop
				}

				if d.Dictionary != nil {
					if _, exists := d.Dictionary[string(key)]; exists {
						marshalErr = &DirectiveMarshalError{
							error: DictionaryDuplicateKeyError,
						}
						break ReadLineLoop
					}
				}

				firstChar, _ = readFirstMeaningfulCharacter(elementContent, true)

				if firstChar != EmptyChar {
					child := &Directive{
						IndentSize: d.IndentSize,
						Depth: d.Depth + 1,
					}

					if childIsString {
						child.Type = DirectiveTypeString
						if firstLine[len(firstLine) - 1] == '\n' {
							child.String = string(firstLine[valueIndex:len(firstLine) - 1])
						} else {
							child.String = string(firstLine[valueIndex:])
						}
					} else {
						if marshalErr = child.Marshal(elementContent); marshalErr != nil {
							if marshalErr.error == EmptyDataError {
								child.Type = DirectiveTypeString
								child.String = ""
								marshalErr = nil
							} else {
								break ReadLineLoop
							}
						}
					}

					if d.Dictionary == nil {
						d.Dictionary = make(map[string]*Directive)
					}
					
					d.Dictionary[string(key)] = child
				} else {
					// empty case
					child := &Directive{
						Type: DirectiveTypeString,
						String: "",
					}
					if d.Dictionary == nil {
						d.Dictionary = make(map[string]*Directive)
					}
					d.Dictionary[string(key)] = child
				}

			}
		}

		if eof {
			break
		}
	}

	if marshalErr == nil && d.Type == DirectiveTypeUnknown {
		marshalErr = &DirectiveMarshalError{
			error: EmptyDataError,
		}
	}

	return marshalErr
}

func (d *Directive) ReadDictionaryDirective(contentBuffer *bytes.Buffer, elementContent *[]byte, indentationIndex int) *DirectiveMarshalError {
	var marshalErr *DirectiveMarshalError = nil

	for {
		line, err := contentBuffer.ReadBytes(byte('\n'))

		eof := (err == io.EOF)
		if !eof && err != nil {
			return &DirectiveMarshalError{
				error: err,
			}
		}

		char, nextIndex := readFirstMeaningfulCharacter(line, true)

		if char == Tab {
			return &DirectiveMarshalError{
				error: TabInIndentationError,
			}
		}

		if nextIndex == indentationIndex {
			if char == ListSymbol || char == TextSymbol {
				// TODO: irregular case
			}
			// it is next element
			return &DirectiveMarshalError{
				error: NextDirectiveAppearedError,
			}
		}

		if nextIndex == NotFoundIndex {
			*elementContent = append(*elementContent, line[0:]...)
		} else {
			*elementContent = append(*elementContent, line[nextIndex:]...)
		}

		if eof {
			break
		}
	}

	return marshalErr
}

func (d *Directive) ReadListDirective(contentBuffer *bytes.Buffer, elementContent *[]byte, indentationIndex int) *DirectiveMarshalError {
	var marshalErr *DirectiveMarshalError = nil

	for {
		line, err := contentBuffer.ReadBytes(byte('\n'))

		eof := (err == io.EOF)
		if !eof && err != nil {
			return &DirectiveMarshalError{
				error: err,
			}
		}

		char, nextIndex := readFirstMeaningfulCharacter(line, true)

		if char == Tab {
			return &DirectiveMarshalError{
				error: TabInIndentationError,
			}
		}

		if nextIndex == indentationIndex {
			if char != ListSymbol {
				// TODO: irregular case
			}
			// it is next element
			return &DirectiveMarshalError{
				error: NextDirectiveAppearedError,
			}
		}

		if nextIndex == NotFoundIndex {
			*elementContent = append(*elementContent, line[0:]...)
		} else {
			*elementContent = append(*elementContent, line[nextIndex:]...)
		}

		if eof {
			break
		}
	}

	return marshalErr
}

func dictionaryKeySanitize(key []byte) ([]byte, error) {
	sanitizedKey := key

	quoted := false

	// remove sorrounding quotes
	{
		if len(sanitizedKey) >= 2 {
			if (key[0] == '"' && key[len(key) - 1] == '"') || (key[0] == '\'' && key[len(key) - 1] == '\'') {
				quoted = true
				sanitizedKey = key[1:len(key) - 1]
				quotedKey := sanitizedKey
				// nested quote is allowed
				if len(sanitizedKey) >= 2 {
					if (sanitizedKey[0] == '"' && sanitizedKey[len(sanitizedKey) - 1] == '"') || (sanitizedKey[0] == '\'' && sanitizedKey[len(sanitizedKey) - 1] == '\'') {
						quotedKey = sanitizedKey[1:len(sanitizedKey) - 1]
					}
				}
				for _, b := range quotedKey {
					if b == '"' || b ==  '\'' {
						return nil, DictionaryKeyNestedQuotesError
					}
				}
			}
		}
	}

	// remove trailing space
	if !quoted {
		index := len(sanitizedKey) - 1
		trimSize := 0

		for index >= 0 {
			if unicode.IsSpace(rune(sanitizedKey[index])) {
				trimSize++
			} else {
				break
			}
			index--
		}
		
		sanitizedKey = sanitizedKey[:len(sanitizedKey) - trimSize]
	}

	return sanitizedKey, nil
}

func readFirstMeaningfulCharacter(line []byte, skipLineBreak bool) (byte, int) {
	var char byte
	var index int

	for len(line) > index {
		char = line[index]

		//if !unicode.IsSpace(rune(char)) {
		if char != Space {
			if char == LineBreak {
				if !skipLineBreak {
					return char, index
				}
			} else {
				return char, index
			}
		}

		if char == Tab {
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
				if unicode.IsSpace(rune(line[index + 1])) {
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
