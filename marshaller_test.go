package ntgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const Sample = `
string: hello
string_ptr: world
text:
  > aaaa
  > bbbb
text_alias:
  > aaaa alias
  > bbbb alias
dict:
  dict_string: world
  dict_text:
    > dict text aaaa
    > dict text bbbb
dict_ptr:
  dict_string: world pointer
  dict_text:
    > dict text aaaa pointer
    > dict text bbbb pointer
list_struct:
  -
    list_string: aaaa
  -
    list_string: bbbb
list_ptr:
  -
    list_string: aaaa pointer
  -
    list_string: bbbb pointer
list_of_list_struct:
  -
    -
      list_string: aaaa nested
    -
      list_string: bbbb nested
  -
    -
      list_string: cccc nested
    -
      list_string: dddd nested
list_of_list_struct_pointer:
  -
    -
      list_string: aaaa nested pointer
    -
      list_string: bbbb nested pointer
  -
    -
      list_string: cccc nested pointer
    -
      list_string: dddd nested pointer
list_text:
  -
    > list text aaaa
    > list text bbbb
  -
    > list text cccc
    > list text dddd
list_string:
  - list string aaaa
  - list string bbbb
`

// TODO: list with different elements that have differed schema

type SampleDict struct {
	DictString string           `nt:"dict_string"`
	DictText   MultilineStrings `nt:"dict_text,multilinestrings"`
}

type SampleListElement struct {
	ListString string `nt:"list_string"`
}

type SampleStruct struct {
	String        string  `nt:"string"`
	StringPointer *string `nt:"string_ptr"`

	Text      []string         `nt:"text,multilinestrings"`
	TextAlias MultilineStrings `nt:"text_alias,multilinestrings"`

	Dict          SampleDict  `nt:"dict"`
	DictOfPointer *SampleDict `nt:"dict_ptr"`

	ListOfStruct              []SampleListElement    `nt:"list_struct"`
	ListOfStructPointer       []*SampleListElement   `nt:"list_ptr"`
	ListOfListOfStruct        [][]SampleListElement  `nt:"list_of_list_struct"`
	ListOfListOfStructPointer [][]*SampleListElement `nt:"list_of_list_struct_pointer"`
	ListOfText                []MultilineStrings     `nt:"list_text,multilinestrings"`
	ListOfString              []string               `nt:"list_string"`

	OmitEmptyString    string `nt:"omit_string,omitempty"`
	NotOmitEmptyString string `nt:"not_omit_string"`

	NoTag string
}

func TestMarshal(t *testing.T) {

	subject := func() *SampleStruct {
		s := &SampleStruct{}
		Marshal(Sample, s)
		return s
	}

	t.Run("string", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "hello", s.String)
	})

	t.Run("string slice", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "aaaa\n", s.Text[0])
		assert.Equal(t, "bbbb", s.Text[1])
	})

	t.Run("string slice alias", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "aaaa alias\n", s.TextAlias[0])
		assert.Equal(t, "bbbb alias", s.TextAlias[1])
	})

	t.Run("struct", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "world", s.Dict.DictString)
		assert.Equal(t, "dict text aaaa\n", s.Dict.DictText[0])
		assert.Equal(t, "dict text bbbb", s.Dict.DictText[1])
	})

	t.Run("reference of struct", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "world pointer", s.DictOfPointer.DictString)
		assert.Equal(t, "dict text aaaa pointer\n", s.DictOfPointer.DictText[0])
		assert.Equal(t, "dict text bbbb pointer", s.DictOfPointer.DictText[1])
	})

	t.Run("slice of struct", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "aaaa", s.ListOfStruct[0].ListString)
		assert.Equal(t, "bbbb", s.ListOfStruct[1].ListString)
	})

	t.Run("slice of slice of struct", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "aaaa nested", s.ListOfListOfStruct[0][0].ListString)
		assert.Equal(t, "bbbb nested", s.ListOfListOfStruct[0][1].ListString)
		assert.Equal(t, "cccc nested", s.ListOfListOfStruct[1][0].ListString)
		assert.Equal(t, "dddd nested", s.ListOfListOfStruct[1][1].ListString)
	})

	t.Run("slice of slice of struct pointer", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "aaaa nested pointer", s.ListOfListOfStructPointer[0][0].ListString)
		assert.Equal(t, "bbbb nested pointer", s.ListOfListOfStructPointer[0][1].ListString)
		assert.Equal(t, "cccc nested pointer", s.ListOfListOfStructPointer[1][0].ListString)
		assert.Equal(t, "dddd nested pointer", s.ListOfListOfStructPointer[1][1].ListString)
	})

	t.Run("slice of text", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "list text aaaa\n", s.ListOfText[0][0])
		assert.Equal(t, "list text bbbb", s.ListOfText[0][1])
		assert.Equal(t, "list text cccc\n", s.ListOfText[1][0])
		assert.Equal(t, "list text dddd", s.ListOfText[1][1])
	})

	t.Run("slice of string", func(t *testing.T) {
		s := subject()

		assert.Equal(t, "list string aaaa", s.ListOfString[0])
		assert.Equal(t, "list string bbbb", s.ListOfString[1])
	})

	t.Run("holistic", func(t *testing.T) {
		HolisticSample := `
string:
  > str with
  > line break
text: one liner text
`
		subject = func() *SampleStruct {
			s := &SampleStruct{}
			Marshal(HolisticSample, s)
			return s
		}

		t.Run("should join multi line text into string", func(t *testing.T) {
			s := subject()
			assert.Equal(t, "str with\nline break", s.String)
		})
		t.Run("should split string into string slice", func(t *testing.T) {
			s := subject()
			assert.Equal(t, 1, len(s.Text))
			assert.Equal(t, "one liner text", s.Text[0])
		})
	})
}

type StringStruct struct {
	Lines string `nt:"key"`
}
type RefStruct struct {
	RefString1 *string `nt:"key1"`
	RefString2 *string `nt:"key2,omitempty"`
}

func TestUnmarshal(t *testing.T) {

	subject := func() string {
		ptr := "str pointer value"
		s := SampleStruct{
			String:        "str value",
			StringPointer: &ptr,
			Text: []string{
				"text value 1",
				"text value 2",
			},
			TextAlias: []string{
				"text alias value 1",
				"text alias value 2",
			},
			Dict: SampleDict{
				DictString: "dict string value",
				DictText: MultilineStrings{
					"dict text value 1",
					"dict text value 2",
				},
			},
			DictOfPointer: &SampleDict{
				DictString: "dict pointer string value",
				DictText: MultilineStrings{
					"dict pointer text value 1",
					"dict pointer text value 2",
				},
			},
			ListOfStruct: []SampleListElement{
				SampleListElement{
					ListString: "list str 1",
				},
				SampleListElement{
					ListString: "list str 2",
				},
			},
			ListOfStructPointer: []*SampleListElement{
				&SampleListElement{
					ListString: "list pointer str 1",
				},
				&SampleListElement{
					ListString: "list pointer str 2",
				},
			},
			ListOfListOfStruct: [][]SampleListElement{
				[]SampleListElement{
					SampleListElement{
						ListString: "list of list str 1",
					},
					SampleListElement{
						ListString: "list of list str 2",
					},
				},
				[]SampleListElement{
					SampleListElement{
						ListString: "list of list str 3",
					},
					SampleListElement{
						ListString: "list of list str 4",
					},
				},
			},
			ListOfListOfStructPointer: [][]*SampleListElement{
				[]*SampleListElement{
					&SampleListElement{
						ListString: "list of list pointer str 1",
					},
					&SampleListElement{
						ListString: "list of list pointer str 2",
					},
				},
				[]*SampleListElement{
					&SampleListElement{
						ListString: "list of list pointer str 3",
					},
					&SampleListElement{
						ListString: "list of list pointer str 4",
					},
				},
			},
			ListOfText: []MultilineStrings{
				MultilineStrings{
					"list of text 1",
					"list of text 2",
				},
				MultilineStrings{
					"list of text 3",
					"list of text 4",
				},
			},
			ListOfString: []string{
				"list of str 1",
				"list of str 2",
			},
		}
		return Unmarshal(s)
	}

	expect := `string: str value
string_ptr: str pointer value
text:
  > text value 1
  > text value 2
text_alias:
  > text alias value 1
  > text alias value 2
dict:
  dict_string: dict string value
  dict_text:
    > dict text value 1
    > dict text value 2
dict_ptr:
  dict_string: dict pointer string value
  dict_text:
    > dict pointer text value 1
    > dict pointer text value 2
list_struct:
  -
    list_string: list str 1
  -
    list_string: list str 2
list_ptr:
  -
    list_string: list pointer str 1
  -
    list_string: list pointer str 2
list_of_list_struct:
  -
    -
      list_string: list of list str 1
    -
      list_string: list of list str 2
  -
    -
      list_string: list of list str 3
    -
      list_string: list of list str 4
list_of_list_struct_pointer:
  -
    -
      list_string: list of list pointer str 1
    -
      list_string: list of list pointer str 2
  -
    -
      list_string: list of list pointer str 3
    -
      list_string: list of list pointer str 4
list_text:
  -
    > list of text 1
    > list of text 2
  -
    > list of text 3
    > list of text 4
list_string:
  - list of str 1
  - list of str 2
not_omit_string: 
`

	t.Run("holistic", func(t *testing.T) {
		s := subject()

		assert.Equal(t, expect, s)
	})

	t.Run("string contains linebreak", func(t *testing.T) {
		s := StringStruct{"line1\nline2"}

		t.Run("should unmarshaled to multi line text", func(t *testing.T) {
			ret := Unmarshal(s)
			assert.Equal(t, "key:\n  > line1\n  > line2", ret)
		})
	})

	t.Run("pointer value is nil", func(t *testing.T) {
		s := RefStruct{nil, nil}

		t.Run("should unmarshaled to multi line text", func(t *testing.T) {
			ret := Unmarshal(s)
			assert.Equal(t, "key1: ", ret)
		})
	})
}
