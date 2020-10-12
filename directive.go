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
	List   []Directive
	Map    map[string]Directive
}

func (d *Directive) Parse(content []byte) (*Directive, error) {
	d.Type = DirectiveTypeUnknown

	buffer := bytes.NewBuffer(content)
	line, err := buffer.ReadBytes(byte('\n'))

	eof := (err == io.EOF)
	if !eof && err != nil {
		log.Fatal(fmt.Sprintf("nestedtext encountered unknown buffer.ReadString error %s", err.Error()))
		return nil, err
	}

	// initial line
	firstMeaningfulChar, index := d.readFirstMeaningfulCharacter(line)
	if firstMeaningfulChar != 0x00 {
		d.SetString(string(line[index:]))
		return d, nil
	}

	for {
		line, err = buffer.ReadBytes(byte('\n'))

		eof := (err == io.EOF)
		if !eof && err != nil {
			log.Fatal(fmt.Sprintf("nestedtext encountered unknown buffer.ReadString error %s", err.Error()))
			break
		}

		firstMeaningfulChar, index = d.readFirstMeaningfulCharacter(line)

		switch firstMeaningfulChar {
			case 0x00: // empty line
			case '#': // comment
			case '>': { // multi line text
				d.Type = DirectiveTypeText
				firstMeaningfulChar, index = d.readFirstMeaningfulCharacter(line[index:])
				d.Text = append(d.Text, string(line[index:]))
			}
			case '-': // list

			}
		}
	}

	return d, nil
}

func (d *Directive) SetString(data string) {
	d.Type   = DirectiveTypeString
	d.String = data
}

func (d *Directive) readFirstMeaningfulCharacter(line []byte) (byte, int) {
	var char  byte
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