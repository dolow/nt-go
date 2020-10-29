package nestedtext

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}

func TestString(t *testing.T) {
	t.Run("string_1", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile("./test/cases/string_1/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child, exists := directive.Dictionary["what makes it green"]
		
		assert.True(t, exists)
		assert.Equal(t, DirectiveTypeString, child.Type)
		assert.Equal(t, "\tgreen\tchilies\t", child.String)
	})

	t.Run("string_2", func(t *testing.T) {
		// string can contain double quote
		dat, _ := ioutil.ReadFile("./test/cases/string_2/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child, exists := directive.Dictionary["key"]
		
		assert.True(t, exists)
		assert.Equal(t, DirectiveTypeString, child.Type)
		assert.Equal(t, `value " value`, child.String)
	})

	t.Run("string_3", func(t *testing.T) {
		// string can contain single quote
		dat, _ := ioutil.ReadFile("./test/cases/string_3/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child, exists := directive.Dictionary["key"]
		
		assert.True(t, exists)
		assert.Equal(t, DirectiveTypeString, child.Type)
		assert.Equal(t, `value ' value`, child.String)
	})

	t.Run("string_4", func(t *testing.T) {
		// string can contain mixed quotes
		dat, _ := ioutil.ReadFile("./test/cases/string_4/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		var child *Directive

		child, _ = directive.Dictionary["key1"]
		assert.Equal(t, `'And Fred said "yabba dabba doo!" to Barney.'`, child.String)
		child, _ = directive.Dictionary["key2"]
		assert.Equal(t, `"And Fred said 'yabba dabba doo!' to Barney."`, child.String)
		child, _ = directive.Dictionary["key3"]
		assert.Equal(t, `"And Fred said "yabba dabba doo!" to Barney."`, child.String)
		child, _ = directive.Dictionary["key4"]
		assert.Equal(t, `'And Fred said 'yabba dabba doo!' to Barney.'`, child.String)
		child, _ = directive.Dictionary["key5"]
		assert.Equal(t, `And Fred said "yabba dabba doo!" to Barney.`, child.String)
		child, _ = directive.Dictionary["key6"]
		assert.Equal(t, `And Fred said 'yabba dabba doo!' to Barney.`, child.String)
	})
	t.Run("string_5", func(t *testing.T) {
		// json with only empty string is converted to empty multiline text
		// TODO: json conversion
	})
	t.Run("string_6", func(t *testing.T) {
		// string can not be begin with new line
		dat, _ := ioutil.ReadFile("./test/cases/string_6/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, StringWithNewLineError, err.error)
		
	})
}

func TestText(t *testing.T) {
	t.Run("text_1", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile("./test/cases/string_multiline_1/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeText, directive.Type)
		assert.Equal(t, 2, len(directive.Text))

		assert.Equal(t, "\n", directive.Text[0])
		assert.Equal(t, "", directive.Text[1])
	})

	t.Run("text_2", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile("./test/cases/string_multiline_2/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeText, directive.Type)
		assert.Equal(t, 2, len(directive.Text))

		// line break excluded
		assert.Equal(t, "ingredients\n", directive.Text[0])
		assert.Equal(t, "green chilies", directive.Text[1])
	})
	t.Run("text_3", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile("./test/cases/string_multiline_3/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child := directive.Dictionary["ingredients"]

		assert.Equal(t, DirectiveTypeText, child.Type)
		assert.Equal(t, 1, len(child.Text))
		assert.Equal(t, "green chilies", child.Text[0])
	})
	t.Run("text_4", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile("./test/cases/string_multiline_4/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child := directive.Dictionary[`key "' key`]

		assert.Equal(t, DirectiveTypeText, child.Type)
		assert.Equal(t, 1, len(child.Text))
		assert.Equal(t, `value '" value`, child.Text[0])
	})

	t.Run("text_5", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile("./test/cases/string_multiline_5/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DifferentLevelOnSameChildError, err.error)
	})

	t.Run("text_6", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile("./test/cases/string_multiline_6/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DifferentLevelOnSameChildError, err.error)
	})

	t.Run("text_7", func(t *testing.T) {
		// complex cases
		dat, _ := ioutil.ReadFile("./test/cases/string_multiline_7/load_in.nt")
		directive := &Directive{}
		err := directive.Marshal(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		t.Run("text_7_1", func(t *testing.T) {
			child := directive.Dictionary["no newlines"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("text_7_2", func(t *testing.T) {
			child := directive.Dictionary["leading newline"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 3, len(child.Text))
			assert.Equal(t, "\n", child.Text[0])
			assert.Equal(t, "line 1\n", child.Text[1])
			assert.Equal(t, "line 2", child.Text[2])
		})
		t.Run("text_7_3", func(t *testing.T) {
			child := directive.Dictionary["internal newline"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 3, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "\n", child.Text[1])
			assert.Equal(t, "line 2", child.Text[2])
		})
		t.Run("text_7_4", func(t *testing.T) {
			child := directive.Dictionary["trailing newline"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 3, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2\n", child.Text[1])
			assert.Equal(t, "", child.Text[2])
		})
		t.Run("text_7_5", func(t *testing.T) {
			child := directive.Dictionary["leading, internal, and trailing newline"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 5, len(child.Text))
			assert.Equal(t, "\n", child.Text[0])
			assert.Equal(t, "line 1\n", child.Text[1])
			assert.Equal(t, "\n", child.Text[2])
			assert.Equal(t, "line 2\n", child.Text[3])
			assert.Equal(t, "", child.Text[4])
		})
		t.Run("text_7_6", func(t *testing.T) {
			child := directive.Dictionary["leading newlines"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 4, len(child.Text))
			assert.Equal(t, "\n", child.Text[0])
			assert.Equal(t, "\n", child.Text[1])
			assert.Equal(t, "line 1\n", child.Text[2])
			assert.Equal(t, "line 2", child.Text[3])
		})
		t.Run("text_7_7", func(t *testing.T) {
			child := directive.Dictionary["internal newlines"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 4, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "\n", child.Text[1])
			assert.Equal(t, "\n", child.Text[2])
			assert.Equal(t, "line 2", child.Text[3])
		})
		t.Run("text_7_8", func(t *testing.T) {
			child := directive.Dictionary["trailing newlines"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 4, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2\n", child.Text[1])
			assert.Equal(t, "\n", child.Text[2])
			assert.Equal(t, "", child.Text[3])
		})
		t.Run("text_7_9", func(t *testing.T) {
			child := directive.Dictionary["leading, internal, and trailing newlines"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 8, len(child.Text))
			assert.Equal(t, "\n", child.Text[0])
			assert.Equal(t, "\n", child.Text[1])
			assert.Equal(t, "line 1\n", child.Text[2])
			assert.Equal(t, "\n", child.Text[3])
			assert.Equal(t, "\n", child.Text[4])
			assert.Equal(t, "line 2\n", child.Text[5])
			assert.Equal(t, "\n", child.Text[6])
			assert.Equal(t, "", child.Text[7])
		})
		t.Run("text_7_10", func(t *testing.T) {
			child := directive.Dictionary["leading blank line"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("text_7_11", func(t *testing.T) {
			child := directive.Dictionary["internal blank line"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("text_7_12", func(t *testing.T) {
			child := directive.Dictionary["trailing blank line"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("text_7_13", func(t *testing.T) {
			child := directive.Dictionary["leading comment"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("text_7_14", func(t *testing.T) {
			child := directive.Dictionary["internal comment"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 3\n", child.Text[0])
			assert.Equal(t, "line 4", child.Text[1])
		})
		t.Run("text_7_15", func(t *testing.T) {
			child := directive.Dictionary["trailing comment"]
			
			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})                       
	})
}

func TestList(t *testing.T) {
	t.Run("list_1", func(t *testing.T) {
		// differing types at same level of indentation causes error
		dat, _ := ioutil.ReadFile("./test/cases/list_1/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		
		assert.Equal(t, DirectiveTypeList, directive.Type)
		assert.Equal(t, 2, len(directive.List))
		for _, entity := range directive.List {
			assert.Equal(t, DirectiveTypeString, entity.Type)
			assert.Equal(t, "", entity.String)
		}
	})

	t.Run("list_2", func(t *testing.T) {
		// differing types at same level of indentation causes error
		dat, _ := ioutil.ReadFile("./test/cases/list_2/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		
		assert.Equal(t, DirectiveTypeList, directive.Type)
		assert.Equal(t, 5, len(directive.List))

		for i, entity := range directive.List {
			if i == 3 {
				assert.Equal(t, DirectiveTypeList, entity.Type)
				assert.Equal(t, 2, len(entity.List))

				for _, childEntity := range entity.List {
					assert.Equal(t, DirectiveTypeString, childEntity.Type)
					assert.True(t, len(childEntity.String) > 0)
				}
			} else {
				assert.Equal(t, DirectiveTypeString, entity.Type)
				assert.True(t, len(entity.String) > 0)
			}
		}
	})

	t.Run("list_3", func(t *testing.T) {
		// empty array should cause error
		// TODO: json conversion
	})

	t.Run("list_4", func(t *testing.T) {
		// list directive can not have dictionary key on the same level
		dat, _ := ioutil.ReadFile("./test/cases/list_4/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DifferentTypesOnTheSameLevelError, err.error)
	})

	t.Run("list_5", func(t *testing.T) {
		// list elements levels can not be defered
		dat, _ := ioutil.ReadFile("./test/cases/list_5/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DifferentLevelOnSameChildError, err.error)
	})

	t.Run("list_6", func(t *testing.T) {
		// list elements levels can not be defered
		dat, _ := ioutil.ReadFile("./test/cases/list_6/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, StringHasChildError, err.error)
	})

	t.Run("list_7", func(t *testing.T) {
		// syntax complexed cases
		dat, _ := ioutil.ReadFile("./test/cases/list_7/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
		
		for i, entity := range directive.List {
			if i == 9 {
				assert.Equal(t, DirectiveTypeText, entity.Type)
			} else {
				assert.Equal(t, DirectiveTypeString, entity.Type)
			}
		}

		t.Run("list_7_1", func(t *testing.T) {
			assert.Equal(t, ":", directive.List[0].String)
		})
		t.Run("list_7_2", func(t *testing.T) {
			assert.Equal(t, `~!@#$%^&*()_+-1234567890{}[]|\:;<>?,./`, directive.List[1].String)
		})
		t.Run("list_7_3", func(t *testing.T) {
			assert.Equal(t, `- value 3`, directive.List[2].String)
		})
		t.Run("list_7_4", func(t *testing.T) {
			assert.Equal(t, `' : value 4:'`, directive.List[3].String)
		})
		t.Run("list_7_5", func(t *testing.T) {
			assert.Equal(t, `> value 5`, directive.List[4].String)
		})
		t.Run("list_7_6", func(t *testing.T) {
			assert.Equal(t, `#value 6`, directive.List[5].String)
		})
		t.Run("list_7_7", func(t *testing.T) {
			assert.Equal(t, `key 7' : : value 7`, directive.List[6].String)
		})
		t.Run("list_7_8", func(t *testing.T) {
			assert.Equal(t, `" value 8 "`, directive.List[7].String)
		})
		t.Run("list_7_9", func(t *testing.T) {
			assert.Equal(t, `' value 9 '`, directive.List[8].String)
		})
		t.Run("list_7_10", func(t *testing.T) {
			assert.Equal(t, 1, len(directive.List[9].Text))
		})
		t.Run("list_7_11", func(t *testing.T) {
			assert.Equal(t, `value '" 10`, directive.List[9].Text[0])
		})
		t.Run("list_7_12", func(t *testing.T) {
			assert.Equal(t, `And Fred said 'yabba dabba doo!' to Barney.`, directive.List[10].String)
		})
	})
}

func TestDictionary(t *testing.T) {	
	t.Run("dict_01", func(t *testing.T) {
		expect := &Directive{
			Type: DirectiveTypeDictionary,
			Dictionary: map[string]*Directive{
				"key1": &Directive{
					Type: DirectiveTypeString,
					String: "",
				},
				"key2": &Directive{
					Type: DirectiveTypeString,
					String: "",
				},
			},
		}
		dat, _ := ioutil.ReadFile("./test/cases/dict_01/load_in.nt")
		
		directive := &Directive{}
		directive.Marshal(dat)
		assert.Equal(t, len(expect.Dictionary), len(directive.Dictionary))
		it := 0
		for k, v := range expect.Dictionary {
			assert.Equal(t, v.Type, directive.Dictionary[k].Type)
			assert.Equal(t, v.String, directive.Dictionary[k].String)
			it++
		}
		assert.Equal(t, len(expect.Dictionary), it)
	})

	t.Run("dict_02", func(t *testing.T) {
		directive := &Directive{}

		err := directive.Marshal([]byte("key\n: value"))
		assert.NotNil(t, err)
		// DifferentTypesOnTheSameLevelError
		assert.Equal(t, RootStringError, err.error)
	})

	t.Run("dict_03", func(t *testing.T) {
		directive := &Directive{}

		var err *DirectiveMarshalError

		err = directive.Marshal([]byte("'ke'y': value"))
		assert.NotNil(t, err)
		assert.Equal(t, DictionaryKeyNestedQuotesError, err.error)

		err = directive.Marshal([]byte("'ke\"y': value"))
		assert.NotNil(t, err)
		assert.Equal(t, DictionaryKeyNestedQuotesError, err.error)
	})

	t.Run("dict_04", func(t *testing.T) {
		directive := &Directive{}

		var err *DirectiveMarshalError

		err = directive.Marshal([]byte("\"ke\"y\": value"))
		assert.NotNil(t, err)
		assert.Equal(t, DictionaryKeyNestedQuotesError, err.error)

		err = directive.Marshal([]byte("\"ke'y\": value"))
		assert.NotNil(t, err)
		assert.Equal(t, DictionaryKeyNestedQuotesError, err.error)
	})
	t.Run("dict_05", func(t *testing.T) {
		// empty json root object should cause error
		// TODO: json conversion
	})
	t.Run("dict_06", func(t *testing.T) {
		dat, _ := ioutil.ReadFile("./test/cases/dict_06/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)
		assert.NotNil(t, err)
		assert.Equal(t, RootLevelHasIndentError, err.error)
	})
	t.Run("dict_07", func(t *testing.T) {
		// causes error when the dictionary elements has different indentation
		dat, _ := ioutil.ReadFile("./test/cases/dict_07/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)
		
		assert.NotNil(t, err)
		// TODO: identify invalid indentation
		assert.Equal(t, StringHasChildError, err.error)
	})
	t.Run("dict_08", func(t *testing.T) {
		// causes error when the indentation contains tab
		dat, _ := ioutil.ReadFile("./test/cases/dict_08/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, TabInIndentationError, err.error)
	})
	t.Run("dict_09", func(t *testing.T) {
		// differencing types on the same level
		dat, _ := ioutil.ReadFile("./test/cases/dict_09/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DifferentTypesOnTheSameLevelError, err.error)
	})
	t.Run("dict_10", func(t *testing.T) {
		// list elements with different indentation causes error, case of following element is deeper
		dat, _ := ioutil.ReadFile("./test/cases/dict_10/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, StringHasChildError, err.error)
	})
	t.Run("dict_11", func(t *testing.T) {
		// list elements with different indentation causes error, case of following element is shallower
		dat, _ := ioutil.ReadFile("./test/cases/dict_11/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DifferentLevelOnSameChildError, err.error)
	})
	t.Run("dict_12", func(t *testing.T) {
		// differing types at same level of indentation causes error
		dat, _ := ioutil.ReadFile("./test/cases/dict_12/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DifferentTypesOnTheSameLevelError, err.error)
	})
	t.Run("dict_13", func(t *testing.T) {
		// string elements starts with new line causes error
		dat, _ := ioutil.ReadFile("./test/cases/dict_13/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, StringWithNewLineError, err.error)
	})
	t.Run("dict_14", func(t *testing.T) {
		// dictionary elements with same key causes error
		dat, _ := ioutil.ReadFile("./test/cases/dict_14/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, DictionaryDuplicateKeyError, err.error)
	})
	t.Run("dict_15", func(t *testing.T) {
		// causes error when the dictionary text elements has tab indentation
		dat, _ := ioutil.ReadFile("./test/cases/dict_15/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.NotNil(t, err)
		assert.Equal(t, TabInIndentationError, err.error)
	})
	t.Run("dict_16", func(t *testing.T) {
		// key can contain ":"
		dat, _ := ioutil.ReadFile("./test/cases/dict_16/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		_, exists := directive.Dictionary["key:"]

		assert.Nil(t, err)
		assert.True(t, exists)
	})
	t.Run("dict_17", func(t *testing.T) {
		// empty key, key sorrounded by quetes and starts with ">", "-", "#", and special characters are ok
		dat, _ := ioutil.ReadFile("./test/cases/dict_17/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)
		assert.Nil(t, err)

		var v *Directive
		var ok bool

		t.Run("dict_17_1", func(t *testing.T) {
			v, ok = directive.Dictionary[""]
			assert.Equal(t, "", v.String)
			assert.True(t, ok)
		})

		t.Run("dict_17_2", func(t *testing.T) {
			v, ok = directive.Dictionary["~!@#$%^&*()_+-1234567890{}[]|\\;<>?,./"]
			assert.Equal(t, "~!@#$%^&*()_+-1234567890{}[]|\\:;<>?,./", v.String)
			assert.True(t, ok)
		})

		t.Run("dict_17_3", func(t *testing.T) {
			v, ok = directive.Dictionary["- key 3"]
			assert.Equal(t, "- value 3", v.String)
			assert.True(t, ok)
		})

		t.Run("dict_17_4", func(t *testing.T) {
			v, ok = directive.Dictionary["key 4: "]
			assert.Equal(t, "value 4: ", v.String)
			assert.True(t, ok)
		})

		t.Run("dict_17_5", func(t *testing.T) {
			v, ok = directive.Dictionary["> key 5"]
			assert.Equal(t, "> value 5", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_6", func(t *testing.T) {
			v, ok = directive.Dictionary["# key 6"]
			assert.Equal(t, "#value 6", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_7", func(t *testing.T) {
			v, ok = directive.Dictionary[": key 7"]
			assert.Equal(t, ": value 7", v.String)
			assert.True(t, ok)
		})

		t.Run("dict_17_8", func(t *testing.T) {
			v, ok = directive.Dictionary["\" key 8 \""]
			assert.Equal(t, "\" value 8 \"", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_9", func(t *testing.T) {
			v, ok = directive.Dictionary["' key 9 '"]
			assert.Equal(t, "' value 9 '", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_10", func(t *testing.T) {
			v, ok = directive.Dictionary["key 10"]
			assert.Equal(t, "value '\" 10", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_11", func(t *testing.T) {
			v, ok = directive.Dictionary["key 11"]
			assert.Equal(t, "And Fred said 'yabba dabba doo!' to Barney.", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_12", func(t *testing.T) {
			v, ok = directive.Dictionary["key \" 12"]
			assert.Equal(t, "value ' 12", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_13", func(t *testing.T) {
			v, ok = directive.Dictionary["$€¥£₩₺₽₹ɃΞȄ"]
			assert.Equal(t, "$€¥£₩₺₽₹ɃΞȄ", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_14", func(t *testing.T) {
			v, ok = directive.Dictionary["YZEPTGMKk_cmuµμnpfazy"]
			assert.Equal(t, "YZEPTGMKk_cmuµμnpfazy", v.String)
			assert.True(t, ok)
		})
		
		t.Run("dict_17_15", func(t *testing.T) {
			v, ok = directive.Dictionary["a-zA-Z%√{us}{cur}][-^/()\\w·⁻⁰¹²³⁴⁵⁶⁷⁸⁹°ÅΩƱΩ℧"]
			assert.Equal(t, "a-zA-Z%√{us}{cur}][-^/()\\w·⁻⁰¹²³⁴⁵⁶⁷⁸⁹°ÅΩƱΩ℧", v.String)
			assert.True(t, ok)
		})
	})
	t.Run("dict_18", func(t *testing.T) {
		// key with quotes without sorrounding quates are ok
		dat, _ := ioutil.ReadFile("./test/cases/dict_18/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
	})
	t.Run("dict_19", func(t *testing.T) {
		// key with trailing white spaces are ignored
		dat, _ := ioutil.ReadFile("./test/cases/dict_19/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
	})
	t.Run("dict_20", func(t *testing.T) {
		// allowed mixed syntaxes cases
		dat, _ := ioutil.ReadFile("./test/cases/dict_20/load_in.nt")

		directive := &Directive{}

		err := directive.Marshal(dat)

		assert.Nil(t, err)
	})
	t.Run("dict_21", func(t *testing.T) {
		// json number for value is not allowed
		// TODO: json conversion
	})
	t.Run("dict_22", func(t *testing.T) {
		// json number for key is not allowed
		// TODO: json conversion
	})
}

func TestEmpty(t *testing.T) {	
	t.Run("empty_1", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile("./test/cases/empty_1/load_in.nt")

		directive := &Directive{}

		var err *DirectiveMarshalError

		err = directive.Marshal(dat)
		assert.NotNil(t, err)
		// TODO: should be parsed and treat as null/nil by converter
		assert.Equal(t, EmptyDataError, err.error)

		err = directive.Marshal([]byte("\n  \n"))
		assert.NotNil(t, err)
		assert.Equal(t, EmptyDataError, err.error)
	})
}

func TestHolistic(t *testing.T) {	
	getValueWithAssert := func (t *testing.T, dict map[string]*Directive, key string) *Directive {
		directive, exists := dict[key]
		assert.True(t, exists)
		return directive
	}

	t.Run("holistic_1", func(t *testing.T) {
		dat, _ := ioutil.ReadFile("./test/cases/holistic_1/load_in.nt")

		directive := &Directive{}
		err := directive.Marshal(dat)

		t.Run("should marshal successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should marshal with collect structure", func(t *testing.T) {
			assert.Equal(t, DirectiveTypeDictionary, directive.Type)

			rootDict := directive.Dictionary
			key1 := getValueWithAssert(t, rootDict, "key 1")
			assert.Equal(t, DirectiveTypeString, key1.Type)
			assert.Equal(t, "value 1", key1.String)

			key2 := getValueWithAssert(t, rootDict, "- key2:")
			assert.Equal(t, DirectiveTypeString, key2.Type)
			assert.Equal(t, "value2:", key2.String)

			key3 := getValueWithAssert(t, rootDict, "  #key3  ")
			assert.Equal(t, DirectiveTypeString, key3.Type)
			assert.Equal(t, "  #value3  ", key3.String)

			key4 := getValueWithAssert(t, rootDict, "key 4")
			assert.Equal(t, DirectiveTypeDictionary, key4.Type)

			dict4 := key4.Dictionary
			{
				key4_1 := getValueWithAssert(t, dict4, "key 4.1")
				assert.Equal(t, DirectiveTypeString, key4_1.Type)
				assert.Equal(t, "value 4.1", key4_1.String)

				key4_2 := getValueWithAssert(t, dict4, "key 4.2")
				assert.Equal(t, DirectiveTypeString, key4_2.Type)
				assert.Equal(t, "value 4.2", key4_2.String)

				key4_3 := getValueWithAssert(t, dict4, "key 4.3")
				assert.Equal(t, DirectiveTypeDictionary, key4_3.Type)

				dict4_3 := key4_3.Dictionary
				{
					key4_3_1 := getValueWithAssert(t, dict4_3, "key 4.3.1")
					assert.Equal(t, DirectiveTypeString, key4_3_1.Type)
					assert.Equal(t, "value 4.3.1", key4_3_1.String)

					key4_3_2 := getValueWithAssert(t, dict4_3, "key 4.3.2")
					assert.Equal(t, DirectiveTypeString, key4_3_2.Type)
					assert.Equal(t, "value 4.3.2", key4_3_2.String)
				}

				key4_4 := getValueWithAssert(t, dict4, "key 4.4")
				assert.Equal(t, DirectiveTypeList, key4_4.Type)
				assert.Equal(t, 3, len(key4_4.List))

				list4_4 := key4_4.List
				{
					assert.Equal(t, DirectiveTypeString, list4_4[0].Type)
					assert.Equal(t, "value 4.4.1", list4_4[0].String)

					assert.Equal(t, DirectiveTypeString, list4_4[1].Type)
					assert.Equal(t, "value 4.4.2", list4_4[1].String)

					assert.Equal(t, DirectiveTypeList, list4_4[2].Type)

					list4_4_2 := list4_4[2].List
					assert.Equal(t, 2, len(list4_4_2))
					{
						assert.Equal(t, DirectiveTypeString, list4_4_2[0].Type)
						assert.Equal(t, "value 4.4.3.1", list4_4_2[0].String)

						assert.Equal(t, DirectiveTypeString, list4_4_2[1].Type)
						assert.Equal(t, "value 4.4.3.2", list4_4_2[1].String)
					}
				}
			}

			key5 := getValueWithAssert(t, rootDict, "key 5")
			assert.Equal(t, DirectiveTypeText, key5.Type)
			assert.Equal(t, 1, len(key5.Text))
			assert.Equal(t, "value 5 part 1", key5.Text[0])

			key6 := getValueWithAssert(t, rootDict, "key 6")
			assert.Equal(t, DirectiveTypeText, key6.Type)
			assert.Equal(t, 2, len(key6.Text))
			assert.Equal(t, "value 6 part 1\n", key6.Text[0])
			assert.Equal(t, "value 6 part 2", key6.Text[1])

			key7 := getValueWithAssert(t, rootDict, "key 7")
			assert.Equal(t, DirectiveTypeText, key7.Type)
			assert.Equal(t, 4, len(key7.Text))
			assert.Equal(t, "value 7 part 1\n", key7.Text[0])
			assert.Equal(t, "\n", key7.Text[1])
			assert.Equal(t, "value 7 part 3\n", key7.Text[2])
			assert.Equal(t, "", key7.Text[3])

			key8 := getValueWithAssert(t, rootDict, "key 8")
			assert.Equal(t, DirectiveTypeList, key8.Type)

			list8 := key8.List
			assert.Equal(t, 2, len(list8))
			{
				assert.Equal(t, DirectiveTypeString, list8[0].Type)
				assert.Equal(t, "value 8.1", list8[0].String)

				assert.Equal(t, DirectiveTypeString, list8[1].Type)
				assert.Equal(t, "value 8.2", list8[1].String)
			}

			key9 := getValueWithAssert(t, rootDict, "key 9")
			assert.Equal(t, DirectiveTypeList, key9.Type)

			list9 := key9.List
			assert.Equal(t, 2, len(list9))
			{
				assert.Equal(t, DirectiveTypeString, list9[0].Type)
				assert.Equal(t, "value 9.1", list9[0].String)

				assert.Equal(t, DirectiveTypeString, list9[1].Type)
				assert.Equal(t, "value 9.2", list9[1].String)
			}

			key10 := getValueWithAssert(t, rootDict, "key 10")
			assert.Equal(t, DirectiveTypeText, key10.Type)
			assert.Equal(t, 1, len(key10.Text))
			assert.Equal(t, "This is a multiline string.  It should end without a newline.", key10.Text[0])

			key11 := getValueWithAssert(t, rootDict, "key 11")
			assert.Equal(t, DirectiveTypeText, key11.Type)
			assert.Equal(t, 2, len(key11.Text))
			assert.Equal(t, "This is a multiline string.  It should end with a newline.\n", key11.Text[0])
			assert.Equal(t, "", key11.Text[1])

			key12 := getValueWithAssert(t, rootDict, "key 12")
			assert.Equal(t, DirectiveTypeText, key12.Type)
			assert.Equal(t, 7, len(key12.Text))
			assert.Equal(t, "\n", key12.Text[0])
			assert.Equal(t, "This is another\n", key12.Text[1])
			assert.Equal(t, "multiline string.\n", key12.Text[2])
			assert.Equal(t, "\n", key12.Text[3])
			assert.Equal(t, "This continues the same string.\n", key12.Text[4])
			assert.Equal(t, "\n", key12.Text[5])
			assert.Equal(t, "", key12.Text[6])

			key13 := getValueWithAssert(t, rootDict, "key 13")
			assert.Equal(t, DirectiveTypeString, key13.Type)
			assert.Equal(t, "", key13.String)
		})
	})

	t.Run("holistic_2", func(t *testing.T) {
		dat, _ := ioutil.ReadFile("./test/cases/holistic_2/load_in.nt")

		directive := &Directive{}
		err := directive.Marshal(dat)

		t.Run("should marshal successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should marshal with collect structure", func(t *testing.T) {
			assert.Equal(t, DirectiveTypeDictionary, directive.Type)

			rootDict := directive.Dictionary
			
			keySrcDir := getValueWithAssert(t, rootDict, "src_dir")
			assert.Equal(t, DirectiveTypeString, keySrcDir.Type)
			assert.Equal(t, "/", keySrcDir.String)

			keyExcludes := getValueWithAssert(t, rootDict, "excludes")
			assert.Equal(t, DirectiveTypeList, keyExcludes.Type)

			listExists := keyExcludes.List
			assert.Equal(t, 10, len(listExists))
			{
				assert.Equal(t, DirectiveTypeString, listExists[0].Type)
				assert.Equal(t, "/dev", listExists[0].String)

				assert.Equal(t, DirectiveTypeString, listExists[1].Type)
				assert.Equal(t, "/home/*/.cache", listExists[1].String)

				assert.Equal(t, DirectiveTypeString, listExists[2].Type)
				assert.Equal(t, "/root/*/.cache", listExists[2].String)

				assert.Equal(t, DirectiveTypeString, listExists[3].Type)
				assert.Equal(t, "/proc", listExists[3].String)

				assert.Equal(t, DirectiveTypeString, listExists[4].Type)
				assert.Equal(t, "/sys", listExists[4].String)

				assert.Equal(t, DirectiveTypeString, listExists[5].Type)
				assert.Equal(t, "/tmp", listExists[5].String)

				assert.Equal(t, DirectiveTypeString, listExists[6].Type)
				assert.Equal(t, "/var/cache", listExists[6].String)

				assert.Equal(t, DirectiveTypeString, listExists[7].Type)
				assert.Equal(t, "/var/lock", listExists[7].String)

				assert.Equal(t, DirectiveTypeString, listExists[8].Type)
				assert.Equal(t, "/var/run", listExists[8].String)

				assert.Equal(t, DirectiveTypeString, listExists[9].Type)
				assert.Equal(t, "/var/tmp", listExists[9].String)
			}

			keyKeep := getValueWithAssert(t, rootDict, "keep")
			assert.Equal(t, DirectiveTypeDictionary, keyKeep.Type)

			dictKeep := keyKeep.Dictionary
			{
				keyHourly := getValueWithAssert(t, dictKeep, "hourly")
				assert.Equal(t, DirectiveTypeString, keyHourly.Type)
				assert.Equal(t, "24", keyHourly.String)

				keyDaily := getValueWithAssert(t, dictKeep, "daily")
				assert.Equal(t, DirectiveTypeString, keyDaily.Type)
				assert.Equal(t, "7", keyDaily.String)

				keyWeekly := getValueWithAssert(t, dictKeep, "weekly")
				assert.Equal(t, DirectiveTypeString, keyWeekly.Type)
				assert.Equal(t, "4", keyWeekly.String)

				keyMonthly := getValueWithAssert(t, dictKeep, "monthly")
				assert.Equal(t, DirectiveTypeString, keyMonthly.Type)
				assert.Equal(t, "12", keyMonthly.String)

				keyYearly := getValueWithAssert(t, dictKeep, "yearly")
				assert.Equal(t, DirectiveTypeString, keyMonthly.Type)
				assert.Equal(t, "5", keyYearly.String)
			}

			keyPassphrase := getValueWithAssert(t, rootDict, "passphrase")
			assert.Equal(t, DirectiveTypeText, keyPassphrase.Type)
			assert.Equal(t, 4, len(keyPassphrase.Text))

		    assert.Equal(t, "trouper segregate militia airway pricey sweetmeat tartan bookstall\n", keyPassphrase.Text[0])
		    assert.Equal(t, "obsession charlady twosome silky puffball grubby ranger notation\n", keyPassphrase.Text[1])
		    assert.Equal(t, "rosebud replicate freshen javelin abbot autocue beater byway\n", keyPassphrase.Text[2])
		    assert.Equal(t, "", keyPassphrase.Text[3])
		})
	})

	t.Run("holistic_3", func(t *testing.T) {
		dat, _ := ioutil.ReadFile("./test/cases/holistic_3/load_in.nt")

		directive := &Directive{}
		err := directive.Marshal(dat)

		t.Run("should marshal successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should marshal with collect structure", func(t *testing.T) {
			assert.Equal(t, DirectiveTypeDictionary, directive.Type)

			rootDict := directive.Dictionary
			
			keyTux := getValueWithAssert(t, rootDict, "tux")
			assert.Equal(t, DirectiveTypeString, keyTux.Type)
			assert.Equal(t, "", keyTux.String)

			keyJux := getValueWithAssert(t, rootDict, "jux")
			assert.Equal(t, DirectiveTypeString, keyJux.Type)
			assert.Equal(t, "lux", keyJux.String)

			keyDux := getValueWithAssert(t, rootDict, "dux")
			assert.Equal(t, DirectiveTypeList, keyDux.Type)

			listDux := keyDux.List
			assert.Equal(t, 6, len(listDux))

			assert.Equal(t, DirectiveTypeString, listDux[0].Type)
			assert.Equal(t, "bux", listDux[0].String)

			assert.Equal(t, DirectiveTypeString, listDux[1].Type)
			assert.Equal(t, "", listDux[1].String)

			assert.Equal(t, DirectiveTypeText, listDux[2].Type)
			assert.Equal(t, 2, len(listDux[2].Text))
			assert.Equal(t, "\n", listDux[2].Text[0])
			assert.Equal(t, "", listDux[2].Text[1])

			assert.Equal(t, DirectiveTypeString, listDux[3].Type)
			assert.Equal(t, "crux", listDux[3].String)

			assert.Equal(t, DirectiveTypeString, listDux[4].Type)
			assert.Equal(t, "", listDux[4].String)

			assert.Equal(t, DirectiveTypeString, listDux[5].Type)
			assert.Equal(t, " — ", listDux[5].String)
		})
	})

	t.Run("holistic_4", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile("./test/cases/holistic_4/load_in.nt")

		directive := &Directive{}
		err := directive.Marshal(dat)

		assert.Nil(t, err)
	})

	t.Run("holistic_5", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile("./test/cases/holistic_5/load_in.nt")

		directive := &Directive{}
		err := directive.Marshal(dat)

		assert.Nil(t, err)
	})

	t.Run("holistic_6", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile("./test/cases/holistic_6/load_in.nt")

		directive := &Directive{}
		err := directive.Marshal(dat)

		assert.Nil(t, err)
	})

	t.Run("holistic_7", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile("./test/cases/holistic_7/load_in.nt")

		directive := &Directive{}
		err := directive.Marshal(dat)

		assert.Nil(t, err)
	})
}