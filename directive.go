package nestedtext

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"unicode"
)

type DirectiveType int

const (
	DirectiveTypeUnknown DirectiveType = iota
	DirectiveTypeString  DirectiveType = iota
	DirectiveTypeText    DirectiveType = iota
	DirectiveTypeList    DirectiveType = iota
	DirectiveTypeMap     DirectiveType = iota
)

type Directive struct {
	Type DirectiveType

	String string
	Text   []string
	List   []*Directive
	Map    map[string]*Directive

	IndentSize int
	Depth      int
}

func (d *Directive) ToString() string {
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

			str = fmt.Sprintf("%s%s- %s%s%s", str, baseIndent, dataLn, child.ToString(), tailLn)
		}
	}
	case DirectiveTypeMap: {
		it := 0
		for k, v := range d.Map {
			dataLn := "\n"
			tailLn := "\n"
			if it == len(d.Map) - 1 {
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

func (d *Directive) Parse(content []byte) (*Directive, error) {
	d.Type = DirectiveTypeUnknown

	// remove trailing line breaks
	for len(content) > 0 && content[len(content) - 1] == '\n' {
		content = content[:len(content) - 1]
	}

	buffer := bytes.NewBuffer(content)

	var line []byte
	var err error
	forwardRetrieval := false

	for {
		if !forwardRetrieval {
			line, err = buffer.ReadBytes(byte('\n'))
		}
		forwardRetrieval = false

		eof := (err == io.EOF)
		if !eof && err != nil {
			log.Fatal(fmt.Sprintf("nestedtext encountered unknown buffer.ReadString error %s", err.Error()))
			break
		}

		firstMeaningfulChar, index := d.readFirstMeaningfulCharacter(line)

		switch firstMeaningfulChar {
		case 0x00: // empty line
		case '#': // comment
		case '>':
			{ // multi line text
				d.Type = DirectiveTypeText
				contentPart := line[index+1:]
				_, contentIndex := d.readFirstMeaningfulCharacter(contentPart)
				d.Text = append(d.Text, string(contentPart[contentIndex:]))
			}
		case '-':
			{ // list
				d.Type = DirectiveTypeList
				elementContent := line[index+1:]

				// TODO: slmost same as map
				for {
					line, err = buffer.ReadBytes(byte('\n'))

					eof := (err == io.EOF)
					if !eof && err != nil {
						log.Fatal(fmt.Sprintf("nestedtext encountered unknown buffer.ReadString error %s", err.Error()))
						break
					}

					char, nextIndex := d.readFirstMeaningfulCharacter(line)

					if nextIndex == index {
						if char != '-' {
							// TODO: irregular case
						}
						// it is next element
						forwardRetrieval = true
						break
					}

					if nextIndex == -1 {
						elementContent = append(elementContent, line[0:]...)
					} else {
						elementContent = append(elementContent, line[nextIndex:]...)
					}

					if eof {
						break
					}
				}

				child := &Directive{
					IndentSize: d.IndentSize,
					Depth: d.Depth + 1,
				}

				_, err = child.Parse(elementContent)
				if err == nil {
					d.List = append(d.List, child)
				} else {
					log.Fatal(fmt.Sprintf("nestedtext encountered unknown default.Parse error %s", err.Error()))
				}
			}
		default:
			{
				sepIndex := d.getCharacterIndex(line[index:], ':')
				if sepIndex == -1 {
					// string
					d.Type = DirectiveTypeString
					// remove trailing line break
					if line[len(line)-1] == '\n' {
						d.String = string(line[index : len(line)-1])
					} else {
						d.String = string(line[index:])
					}
					break
				}

				d.Type = DirectiveTypeMap
				key := line[index:sepIndex]
				elementContent := line[sepIndex+1:]

				for {
					line, err = buffer.ReadBytes(byte('\n'))

					eof := (err == io.EOF)
					if !eof && err != nil {
						log.Fatal(fmt.Sprintf("nestedtext encountered unknown buffer.ReadString error %s", err.Error()))
						break
					}

					char, nextIndex := d.readFirstMeaningfulCharacter(line)

					if nextIndex == index {
						if char == '-' || char == '>' {
							// TODO: irregular case
						}
						// it is next element
						forwardRetrieval = true
						break
					}

					if nextIndex == -1 {
						elementContent = append(elementContent, line[0:]...)
					} else {
						elementContent = append(elementContent, line[nextIndex:]...)
					}

					if eof {
						break
					}
				}

				child := &Directive{
					IndentSize: d.IndentSize,
					Depth: d.Depth + 1,
				}

				_, err = child.Parse(elementContent)
				if err == nil {
					if d.Map == nil {
						d.Map = make(map[string]*Directive)
					}
					d.Map[string(key)] = child
				} else {
					log.Fatal(fmt.Sprintf("nestedtext encountered unknown default.Parse error %s", err.Error()))
				}

			}
		}

		if eof {
			break
		}
	}

	return d, nil
}

func (d *Directive) SetString(data string) {
	d.Type = DirectiveTypeString
	d.String = data
}

func (d *Directive) readFirstMeaningfulCharacter(line []byte) (byte, int) {
	var char byte
	var index int

	for len(line) > index {
		char = line[index]

		if !unicode.IsSpace(rune(char)) {
			return char, index
		}

		index++
	}

	return 0x00, -1
}

func (d *Directive) getCharacterIndex(line []byte, character byte) int {
	var index int
	var char byte

	for len(line) > index {
		char = line[index]

		if character == char {
			return index
		}

		index++
	}

	return -1
}
