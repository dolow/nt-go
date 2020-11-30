package ntgo

import (
	"reflect"
	"strings"
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
list_string_pointer:
  - list string pointer aaaa
  - list string pointer bbbb`

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
	ListOfStringPointer       []*string              `nt:"list_string_pointer"`

	OmitEmptyString    string `nt:"omit_string,omitempty"`
	NotOmitEmptyString string `nt:"not_omit_string"`

	NoTag string
}

type StringStruct struct {
	Str string `nt:"key"`
}
type RefStruct struct {
	RefString1 *string `nt:"key1"`
	RefString2 *string `nt:"key2,omitempty"`
}

type NumberStruct struct {
	Int        int      `nt:"int"`
	Float32    float32  `nt:"float"`
	IntPtr     *int     `nt:"int_ptr"`
	Float32Ptr *float32 `nt:"float_ptr"`

	IntSlice        []int      `nt:"int_slice"`
	Float32Slice    []float32  `nt:"float_slice"`
	IntPtrSlice     []*int     `nt:"int_ptr_slice"`
	Float32PtrSlice []*float32 `nt:"float_ptr_slice"`
}

func TestMarshal(t *testing.T) {

	subject := func() (*SampleStruct, error) {
		s := &SampleStruct{}
		err := Marshal(Sample, s)
		return s, err
	}

	t.Run("when argument is not pointer", func(t *testing.T) {

		arg := SampleStruct{}

		t.Run("should return ValueIsNotPointerError", func(t *testing.T) {
			err := Marshal(Sample, arg)
			assert.Equal(t, ValueIsNotPointerError, err)
		})
	})

	t.Run("string", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "hello", s.String)
	})

	t.Run("string slice", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "aaaa\n", s.Text[0])
		assert.Equal(t, "bbbb", s.Text[1])
	})

	t.Run("string slice alias", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "aaaa alias\n", s.TextAlias[0])
		assert.Equal(t, "bbbb alias", s.TextAlias[1])
	})

	t.Run("struct", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "world", s.Dict.DictString)
		assert.Equal(t, "dict text aaaa\n", s.Dict.DictText[0])
		assert.Equal(t, "dict text bbbb", s.Dict.DictText[1])
	})

	t.Run("reference of struct", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "world pointer", s.DictOfPointer.DictString)
		assert.Equal(t, "dict text aaaa pointer\n", s.DictOfPointer.DictText[0])
		assert.Equal(t, "dict text bbbb pointer", s.DictOfPointer.DictText[1])
	})

	t.Run("slice of struct", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "aaaa", s.ListOfStruct[0].ListString)
		assert.Equal(t, "bbbb", s.ListOfStruct[1].ListString)
	})

	t.Run("slice of slice of struct", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "aaaa nested", s.ListOfListOfStruct[0][0].ListString)
		assert.Equal(t, "bbbb nested", s.ListOfListOfStruct[0][1].ListString)
		assert.Equal(t, "cccc nested", s.ListOfListOfStruct[1][0].ListString)
		assert.Equal(t, "dddd nested", s.ListOfListOfStruct[1][1].ListString)
	})

	t.Run("slice of slice of struct pointer", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "aaaa nested pointer", s.ListOfListOfStructPointer[0][0].ListString)
		assert.Equal(t, "bbbb nested pointer", s.ListOfListOfStructPointer[0][1].ListString)
		assert.Equal(t, "cccc nested pointer", s.ListOfListOfStructPointer[1][0].ListString)
		assert.Equal(t, "dddd nested pointer", s.ListOfListOfStructPointer[1][1].ListString)
	})

	t.Run("slice of text", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "list text aaaa\n", s.ListOfText[0][0])
		assert.Equal(t, "list text bbbb", s.ListOfText[0][1])
		assert.Equal(t, "list text cccc\n", s.ListOfText[1][0])
		assert.Equal(t, "list text dddd", s.ListOfText[1][1])
	})

	t.Run("slice of string", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "list string aaaa", s.ListOfString[0])
		assert.Equal(t, "list string bbbb", s.ListOfString[1])
	})

	t.Run("slice of string pointer", func(t *testing.T) {
		s, _ := subject()

		assert.Equal(t, "list string pointer aaaa", *s.ListOfStringPointer[0])
		assert.Equal(t, "list string pointer bbbb", *s.ListOfStringPointer[1])
	})

	t.Run("holistic", func(t *testing.T) {
		HolisticSample := `
string:
  > str with
  > line break
text: one liner text
`
		subject = func() (*SampleStruct, error) {
			s := &SampleStruct{}
			err := Marshal(HolisticSample, s)
			return s, err
		}

		t.Run("should join multi line text into string", func(t *testing.T) {
			s, _ := subject()
			assert.Equal(t, "str with\nline break", s.String)
		})
		t.Run("should split string into string slice", func(t *testing.T) {
			s, _ := subject()
			assert.Equal(t, 1, len(s.Text))
			assert.Equal(t, "one liner text", s.Text[0])
		})
	})
}

func TestUnmarshal(t *testing.T) {

	subject := func() string {
		ptr := "str pointer value"
		listStrPtr1 := "list of str ptr 1"
		listStrPtr2 := "list of str ptr 2"
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
			ListOfStringPointer: []*string{
				&listStrPtr1,
				&listStrPtr2,
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
list_string_pointer:
  - list of str ptr 1
  - list of str ptr 2
not_omit_string: `

	t.Run("holistic", func(t *testing.T) {
		s := subject()

		assert.Equal(t, expect, s)
	})

	t.Run("struct with number", func(t *testing.T) {
		var i int = 123456
		var f float32 = 1.23456
		s := NumberStruct{
			-123456,
			-1.23456,
			&i,
			&f,
			[]int{123, 456},
			[]float32{1.23, 4.56},
			[]*int{&i}, []*float32{&f},
		}
		ret := Unmarshal(s)

		t.Run("should unmarshaled to string", func(t *testing.T) {
			lines := strings.Split(ret, "\n")
			assert.Equal(t, "int: -123456", lines[0])
			assert.Equal(t, "float: -1.2345", lines[1][0:14])
			assert.Equal(t, "int_ptr: 123456", lines[2])
			assert.Equal(t, "float_ptr: 1.2345", lines[3][0:17])
			assert.Equal(t, "int_slice:", lines[4])
			assert.Equal(t, "  - 123", lines[5])
			assert.Equal(t, "  - 456", lines[6])
			assert.Equal(t, "float_slice:", lines[7])
			assert.Equal(t, "  - 1.2", lines[8][0:7])
			assert.Equal(t, "  - 4.5", lines[9][0:7])
			assert.Equal(t, "int_ptr_slice:", lines[10])
			assert.Equal(t, "  - 123456", lines[11])
			assert.Equal(t, "float_ptr_slice:", lines[12])
			assert.Equal(t, "  - 1.2345", lines[13][0:10])
			assert.Equal(t, "", lines[14])
		})
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

func TestMarshalSlice(t *testing.T) {
	var value *Value
	var sliceType reflect.Type
	var elementType reflect.Type
	var elementRef *reflect.Value

	condition := func() {}

	subject := func() interface{} {
		marshalSlice(value, elementType, elementRef)
		return elementRef.Interface()
	}

	t.Run("string", func(t *testing.T) {
		t.Run("from multiline strings", func(t *testing.T) {
			condition = func() {
				value = &Value{
					Type: ValueTypeText,
					Text: MultilineStrings{
						"line 1\n",
						"line 2",
					},
				}
				elementType = sliceType.Elem()

				instance := reflect.MakeSlice(sliceType, 0, cap(value.Text))
				elementRef = &instance
			}

			t.Run("to string array", func(t *testing.T) {
				sliceType = reflect.TypeOf([]string{})

				t.Run("should pass", func(t *testing.T) {
					condition()

					ret := subject().([]string)
					assert.Equal(t, len(value.Text), len(ret))
					for i, v := range ret {
						assert.Equal(t, value.Text[i], v)
					}
				})
			})

			t.Run("to string pointer slice", func(t *testing.T) {
				sliceType = reflect.TypeOf([]*string{})

				t.Run("should pass", func(t *testing.T) {
					condition()

					ret := subject().([]*string)
					assert.Equal(t, len(value.Text), len(ret))
					for i, v := range ret {
						assert.Equal(t, value.Text[i], *v)
						assert.Equal(t, &value.Text[i], v)
					}
				})
			})
		})

		t.Run("from list of strings", func(t *testing.T) {
			condition = func() {
				value = &Value{
					Type: ValueTypeList,
					List: []*Value{
						&Value{
							Type:   ValueTypeString,
							String: "line 1",
						},
						&Value{
							Type:   ValueTypeString,
							String: "line 2",
						},
					},
				}

				instance := reflect.MakeSlice(sliceType, 0, cap(value.List))

				elementType = sliceType.Elem()
				elementRef = &instance
			}

			t.Run("to string slice", func(t *testing.T) {

				sliceType = reflect.TypeOf([]string{})

				t.Run("should pass", func(t *testing.T) {
					condition()

					ret := subject().([]string)
					assert.Equal(t, len(value.List), len(ret))
					for i, v := range ret {
						assert.Equal(t, value.List[i].String, v)
					}
				})
			})

			t.Run("to string pointer slice", func(t *testing.T) {

				sliceType = reflect.TypeOf([]*string{})

				t.Run("should pass", func(t *testing.T) {
					condition()

					ret := subject().([]*string)
					assert.Equal(t, len(value.List), len(ret))
					for i, v := range ret {
						assert.Equal(t, value.List[i].String, *v)
						assert.Equal(t, &value.List[i].String, v)
					}
				})
			})
		})

		t.Run("from nested list of string to slice of string slice", func(t *testing.T) {
			condition = func() {
				value = &Value{
					Type: ValueTypeList,
					List: []*Value{
						&Value{
							Type: ValueTypeList,
							List: []*Value{
								&Value{
									Type:   ValueTypeString,
									String: "elem 1 of 1",
								},
								&Value{
									Type:   ValueTypeString,
									String: "elem 2 of 1",
								},
							},
						},
						&Value{
							Type: ValueTypeList,
							List: []*Value{
								&Value{
									Type:   ValueTypeString,
									String: "elem 1 of 2",
								},
								&Value{
									Type:   ValueTypeString,
									String: "elem 2 of 2",
								},
							},
						},
					},
				}
				instance := reflect.MakeSlice(sliceType, 0, cap(value.List))

				elementType = sliceType.Elem()
				elementRef = &instance
			}

			t.Run("to slice of string slice", func(t *testing.T) {

				sliceType = reflect.TypeOf([][]string{})

				t.Run("should pass", func(t *testing.T) {
					condition()

					ret := subject().([][]string)
					assert.Equal(t, len(value.List), len(ret))
					for i, v := range ret {
						childList := value.List[i].List

						assert.Equal(t, len(childList), len(v))

						for ci, cv := range v {
							assert.Equal(t, childList[ci].String, cv)
						}
					}
				})
			})
		})
	})

	t.Run("struct", func(t *testing.T) {
		t.Run("from list of structure", func(t *testing.T) {
			condition = func() {
				value = &Value{
					Type: ValueTypeList,
					List: []*Value{
						&Value{
							Type: ValueTypeDictionary,
							Dictionary: map[string]*Value{
								"key": &Value{
									Type:   ValueTypeString,
									String: "item 1",
								},
							},
						},
						&Value{
							Type: ValueTypeDictionary,
							Dictionary: map[string]*Value{
								"key": &Value{
									Type:   ValueTypeString,
									String: "item 2",
								},
							},
						},
					},
				}

				instance := reflect.MakeSlice(sliceType, 0, cap(value.List))

				elementType = sliceType.Elem()
				elementRef = &instance
			}
			t.Run("to slice of struct", func(t *testing.T) {

				sliceType = reflect.TypeOf([]StringStruct{})

				t.Run("should pass", func(t *testing.T) {
					condition()
					ret := subject().([]StringStruct)
					assert.Equal(t, len(value.List), len(ret))
					for i, v := range ret {
						assert.Equal(t, value.List[i].Dictionary["key"].String, v.Str)
					}
				})
			})

			t.Run("to slice of struct pointer", func(t *testing.T) {

				sliceType = reflect.TypeOf([]*StringStruct{})

				t.Run("should pass", func(t *testing.T) {
					condition()
					ret := subject().([]*StringStruct)
					assert.Equal(t, len(value.List), len(ret))
					for i, v := range ret {
						assert.Equal(t, value.List[i].Dictionary["key"].String, v.Str)
					}
				})
			})
		})
	})
}
