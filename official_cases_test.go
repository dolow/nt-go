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

		assert.NotNil(t, err)
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
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
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

		assert.Equal(t, ":", directive.List[0].String)
		assert.Equal(t, `~!@#$%^&*()_+-1234567890{}[]|\:;<>?,./`, directive.List[1].String)
		assert.Equal(t, `- value 3`, directive.List[2].String)
		assert.Equal(t, `' : value 4:'`, directive.List[3].String)
		assert.Equal(t, `> value 5`, directive.List[4].String)
		assert.Equal(t, `#value 6`, directive.List[5].String)
		assert.Equal(t, `key 7' : : value 7`, directive.List[6].String)
		assert.Equal(t, `" value 8 "`, directive.List[7].String)
		assert.Equal(t, `' value 9 '`, directive.List[8].String)
		assert.Equal(t, 1, len(directive.List[9].Text))
		assert.Equal(t, `value '" 10`, directive.List[9].Text[0])
		assert.Equal(t, `And Fred said 'yabba dabba doo!' to Barney.`, directive.List[10].String)
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

		v, ok = directive.Dictionary[""]

		assert.Equal(t, "", v.String)
		assert.True(t, ok)

		v, ok = directive.Dictionary["~!@#$%^&*()_+-1234567890{}[]|\\;<>?,./"]
		assert.Equal(t, "~!@#$%^&*()_+-1234567890{}[]|\\:;<>?,./", v.String)
		assert.True(t, ok)

		v, ok = directive.Dictionary["- key 3"]
		assert.Equal(t, "- value 3", v.String)
		assert.True(t, ok)

		v, ok = directive.Dictionary["key 4: "]
		assert.Equal(t, "value 4: ", v.String)
		assert.True(t, ok)

		v, ok = directive.Dictionary["> key 5"]
		assert.Equal(t, "> value 5", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["# key 6"]
		assert.Equal(t, "#value 6", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary[": key 7"]
		assert.Equal(t, ": value 7", v.String)
		assert.True(t, ok)

		v, ok = directive.Dictionary["\" key 8 \""]
		assert.Equal(t, "\" value 8 \"", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["' key 9 '"]
		assert.Equal(t, "' value 9 '", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["key 10"]
		assert.Equal(t, "value '\" 10", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["key 11"]
		assert.Equal(t, "And Fred said 'yabba dabba doo!' to Barney.", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["key \" 12"]
		assert.Equal(t, "value ' 12", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["$€¥£₩₺₽₹ɃΞȄ"]
		assert.Equal(t, "$€¥£₩₺₽₹ɃΞȄ", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["YZEPTGMKk_cmuµμnpfazy"]
		assert.Equal(t, "YZEPTGMKk_cmuµμnpfazy", v.String)
		assert.True(t, ok)
		
		v, ok = directive.Dictionary["a-zA-Z%√{us}{cur}][-^/()\\w·⁻⁰¹²³⁴⁵⁶⁷⁸⁹°ÅΩƱΩ℧"]
		assert.Equal(t, "a-zA-Z%√{us}{cur}][-^/()\\w·⁻⁰¹²³⁴⁵⁶⁷⁸⁹°ÅΩƱΩ℧", v.String)
		assert.True(t, ok)
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

