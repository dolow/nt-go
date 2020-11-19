package ntgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const Sample = `
string: hello
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
`

// TODO: list with different elements that have differed schema

type SampleDict struct {
	DictString string `nt:"dict_string"`
	DictText []string `nt:"dict_text"`
}

type SampleListElement struct {
	ListString string `nt:"list_string"`
}

type StringSlice []string

type SampleStruct struct {
	String string `nt:"string"`

	Text []string `nt:"text"`
	TextAlias  StringSlice `nt:"text_alias"`

	Dict SampleDict `nt:"dict"`
	DictOfPointer *SampleDict `nt:"dict_ptr"`

	ListOfStruct []SampleListElement `nt:"list_struct"`
	ListOfStructPointer []*SampleListElement `nt:"list_ptr"`
	ListOfListOfStruct [][]SampleListElement `nt:"list_of_list_struct"`
	ListOfListOfStructPointer [][]SampleListElement `nt:"list_of_list_struct_pointer"`
	ListOfText [][]string `nt:"list_text"`

	NoTag string
}

func TestMarshal(t *testing.T) {

	subject := func() *SampleStruct {
		m := &Marshaller{}
		s := &SampleStruct{}
		m.Marshal(Sample, s)
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
}
