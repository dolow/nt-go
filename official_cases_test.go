package ntgo

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestCasePath = "./official_cases/test_cases"

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
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_1/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child, exists := directive.Dictionary["what makes it green"]

		assert.True(t, exists)
		assert.Equal(t, DirectiveTypeString, child.Type)
		assert.Equal(t, "\tgreen\tchilies\t", child.String)
	})

	t.Run("string_2", func(t *testing.T) {
		// string can contain double quote
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_2/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child, exists := directive.Dictionary["key"]

		assert.True(t, exists)
		assert.Equal(t, DirectiveTypeString, child.Type)
		assert.Equal(t, `value " value`, child.String)
	})

	t.Run("string_3", func(t *testing.T) {
		// string can contain single quote
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_3/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child, exists := directive.Dictionary["key"]

		assert.True(t, exists)
		assert.Equal(t, DirectiveTypeString, child.Type)
		assert.Equal(t, `value ' value`, child.String)
	})

	t.Run("string_4", func(t *testing.T) {
		// string can contain mixed quotes
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_4/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

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
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_6/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, StringWithNewLineError, err)

	})
	t.Run("string_7", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_7/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)

		assert.Equal(t, DirectiveTypeText, directive.Type)
		assert.Equal(t, 1, len(directive.Text))
		assert.Equal(t, "what makes it green?", directive.Text[0])
	})
	t.Run("string_8", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_8/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)

		assert.Equal(t, DirectiveTypeText, directive.Type)
		assert.Equal(t, 1, len(directive.Text))
		assert.Equal(t, "what makes it green?", directive.Text[0])
	})
	t.Run("string_9", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_9/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)

		assert.Equal(t, DirectiveTypeText, directive.Type)
		assert.Equal(t, 1, len(directive.Text))
		assert.Equal(t, "what makes it green?", directive.Text[0])
	})
}

func TestText(t *testing.T) {
	t.Run("string_multiline_1", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_01/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeText, directive.Type)
		assert.Equal(t, 2, len(directive.Text))

		assert.Equal(t, "\n", directive.Text[0])
		assert.Equal(t, "", directive.Text[1])
	})

	t.Run("string_multiline_2", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_02/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeText, directive.Type)
		assert.Equal(t, 24, len(directive.Text))

		// line break excluded
		assert.Equal(t, "\n", directive.Text[0])
		assert.Equal(t, "Lorem Ipsum\n", directive.Text[1])
		assert.Equal(t, "\n", directive.Text[2])
		assert.Equal(t, "    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod\n", directive.Text[3])
		assert.Equal(t, "tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, \n", directive.Text[4])
		assert.Equal(t, "quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo \n", directive.Text[5])
		assert.Equal(t, "consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse \n", directive.Text[6])
		assert.Equal(t, "cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat \n", directive.Text[7])
		assert.Equal(t, "non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n", directive.Text[8])
		assert.Equal(t, "\n", directive.Text[9])
		assert.Equal(t, "\n", directive.Text[10])
		assert.Equal(t, "    Sed ut perspiciatis unde omnis iste natus error sit voluptatem\n", directive.Text[11])
		assert.Equal(t, "accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab \n", directive.Text[12])
		assert.Equal(t, "illo inventore veritatis et quasi architecto beatae vitae dicta sunt \n", directive.Text[13])
		assert.Equal(t, "explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit \n", directive.Text[14])
		assert.Equal(t, "aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem \n", directive.Text[15])
		assert.Equal(t, "sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit \n", directive.Text[16])
		assert.Equal(t, "amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora \n", directive.Text[17])
		assert.Equal(t, "incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad \n", directive.Text[18])
		assert.Equal(t, "minima veniam, quis nostrum exercitationem ullam corporis suscipit \n", directive.Text[19])
		assert.Equal(t, "laboriosam, nisi ut aliquid ex ea commodi consequatur? Quis autem vel eum \n", directive.Text[20])
		assert.Equal(t, "iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae \n", directive.Text[21])
		assert.Equal(t, "consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?\"\n", directive.Text[22])
		assert.Equal(t, "", directive.Text[23])
	})
	t.Run("string_multiline_3", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_03/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child := directive.Dictionary["ingredients"]

		assert.Equal(t, DirectiveTypeText, child.Type)
		assert.Equal(t, 1, len(child.Text))
		assert.Equal(t, "green chilies", child.Text[0])
	})
	t.Run("string_multiline_4", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_04/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		assert.Nil(t, err)
		assert.Equal(t, DirectiveTypeDictionary, directive.Type)

		child := directive.Dictionary[`key "' key`]

		assert.Equal(t, DirectiveTypeText, child.Type)
		assert.Equal(t, 1, len(child.Text))
		assert.Equal(t, `value '" value`, child.Text[0])
	})

	t.Run("string_multiline_5", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_05/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentLevelOnSameChildError, err)
	})

	t.Run("string_multiline_6", func(t *testing.T) {
		// string can contain tab characters
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_06/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentLevelOnSameChildError, err)
	})

	t.Run("string_multiline_7", func(t *testing.T) {
		// complex cases
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_07/load_in.nt")
		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, RootLevelHasIndentError, err)
	})
	t.Run("string_multiline_8", func(t *testing.T) {
		// complex cases
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_08/load_in.nt")
		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, TabInIndentationError, err)
	})
	t.Run("string_multiline_9", func(t *testing.T) {
		// complex cases
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_09/load_in.nt")
		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, TabInIndentationError, err)
	})
	t.Run("string_multiline_10", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_10/load_in.nt")
		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentLevelOnSameChildError, err)
	})
	t.Run("string_multiline_11", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_11/load_in.nt")
		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)

		t.Run("string_multiline_11_1", func(t *testing.T) {
			child := directive.Dictionary["no newlines"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("string_multiline_11_2", func(t *testing.T) {
			child := directive.Dictionary["leading newline"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 3, len(child.Text))
			assert.Equal(t, "\n", child.Text[0])
			assert.Equal(t, "line 1\n", child.Text[1])
			assert.Equal(t, "line 2", child.Text[2])
		})
		t.Run("string_multiline_11_3", func(t *testing.T) {
			child := directive.Dictionary["internal newline"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 3, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "\n", child.Text[1])
			assert.Equal(t, "line 2", child.Text[2])
		})
		t.Run("string_multiline_11_4", func(t *testing.T) {
			child := directive.Dictionary["trailing newline"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 3, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2\n", child.Text[1])
			assert.Equal(t, "", child.Text[2])
		})
		t.Run("string_multiline_11_5", func(t *testing.T) {
			child := directive.Dictionary["leading, internal, and trailing newline"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 5, len(child.Text))
			assert.Equal(t, "\n", child.Text[0])
			assert.Equal(t, "line 1\n", child.Text[1])
			assert.Equal(t, "\n", child.Text[2])
			assert.Equal(t, "line 2\n", child.Text[3])
			assert.Equal(t, "", child.Text[4])
		})
		t.Run("string_multiline_11_6", func(t *testing.T) {
			child := directive.Dictionary["leading newlines"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 4, len(child.Text))
			assert.Equal(t, "\n", child.Text[0])
			assert.Equal(t, "\n", child.Text[1])
			assert.Equal(t, "line 1\n", child.Text[2])
			assert.Equal(t, "line 2", child.Text[3])
		})
		t.Run("string_multiline_11_7", func(t *testing.T) {
			child := directive.Dictionary["internal newlines"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 4, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "\n", child.Text[1])
			assert.Equal(t, "\n", child.Text[2])
			assert.Equal(t, "line 2", child.Text[3])
		})
		t.Run("string_multiline_11_8", func(t *testing.T) {
			child := directive.Dictionary["trailing newlines"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 4, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2\n", child.Text[1])
			assert.Equal(t, "\n", child.Text[2])
			assert.Equal(t, "", child.Text[3])
		})
		t.Run("string_multiline_11_9", func(t *testing.T) {
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
		t.Run("string_multiline_11_10", func(t *testing.T) {
			child := directive.Dictionary["leading blank line"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("string_multiline_11_11", func(t *testing.T) {
			child := directive.Dictionary["internal blank line"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("string_multiline_11_12", func(t *testing.T) {
			child := directive.Dictionary["trailing blank line"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("string_multiline_11_13", func(t *testing.T) {
			child := directive.Dictionary["leading comment"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("string_multiline_11_14", func(t *testing.T) {
			child := directive.Dictionary["internal comment"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
		t.Run("string_multiline_11_15", func(t *testing.T) {
			child := directive.Dictionary["trailing comment"]

			assert.Equal(t, DirectiveTypeText, child.Type)
			assert.Equal(t, 2, len(child.Text))
			assert.Equal(t, "line 1\n", child.Text[0])
			assert.Equal(t, "line 2", child.Text[1])
		})
	})
	t.Run("string_multiline_12", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/string_multiline_12/load_in.nt")
		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)

		assert.Equal(t, "\n", directive.Text[0])
		assert.Equal(t, "the BS character \\	Backslash (\\)\n", directive.Text[1])
		assert.Equal(t, "the SQ character '	Single quote (')\n", directive.Text[2])
		assert.Equal(t, "the DQ character \"	Double quote (\")\n", directive.Text[3])
		assert.Equal(t, "the AB character 	ASCII Bell (BEL)\n", directive.Text[4])
		assert.Equal(t, "the BS character 	ASCII Backspace (BS)\n", directive.Text[5])
		assert.Equal(t, "the FF character 	ASCII Formfeed (FF)\n", directive.Text[6])
		assert.Equal(t, "the LF character \n", directive.Text[7])
		assert.Equal(t, "	ASCII Linefeed (LF)\n", directive.Text[8])
		assert.Equal(t, "the CR character \n", directive.Text[9])
		assert.Equal(t, "	ASCII Carriage Return (CR)\n", directive.Text[10])
		assert.Equal(t, "the HT character 		ASCII Horizontal Tab (TAB)\n", directive.Text[11])
		assert.Equal(t, "the VT character 	ASCII Vertical Tab (VT)\n", directive.Text[12])
		assert.Equal(t, "the ES character 	ASCII escape character as octal value\n", directive.Text[13])
		assert.Equal(t, "the ES character 	ASCII escape character as hex value\n", directive.Text[14])
		assert.Equal(t, "", directive.Text[15])
	})
}

func TestList(t *testing.T) {
	t.Run("list_1", func(t *testing.T) {
		// differing types at same level of indentation causes error
		dat, _ := ioutil.ReadFile(TestCasePath + "/list_1/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

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
		dat, _ := ioutil.ReadFile(TestCasePath + "/list_2/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

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
		// json conversion
	})

	t.Run("list_4", func(t *testing.T) {
		// list directive can not have dictionary key on the same level
		dat, _ := ioutil.ReadFile(TestCasePath + "/list_4/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentTypesOnTheSameLevelError, err)
	})

	t.Run("list_5", func(t *testing.T) {
		// list elements levels can not be defered
		dat, _ := ioutil.ReadFile(TestCasePath + "/list_5/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, RootLevelHasIndentError, err)
	})

	t.Run("list_6", func(t *testing.T) {
		// list elements levels can not be defered
		dat, _ := ioutil.ReadFile(TestCasePath + "/list_6/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, StringHasChildError, err)
	})

	t.Run("list_7", func(t *testing.T) {
		// syntax complexed cases
		dat, _ := ioutil.ReadFile(TestCasePath + "/list_7/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, TabInIndentationError, err)
	})

	t.Run("list_8", func(t *testing.T) {
		// syntax complexed cases
		dat, _ := ioutil.ReadFile(TestCasePath + "/list_8/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)
		assert.Equal(t, 11, len(directive.List))

		assert.Equal(t, ":", directive.List[0].String)
		assert.Equal(t, "~!@#$%^&*()_+-1234567890{}[]|\\:;<>?,./", directive.List[1].String)
		assert.Equal(t, "- value 3", directive.List[2].String)
		assert.Equal(t, "' : value 4:'", directive.List[3].String)
		assert.Equal(t, "> value 5", directive.List[4].String)
		assert.Equal(t, "#value 6", directive.List[5].String)
		assert.Equal(t, "key 7' : : value 7", directive.List[6].String)
		assert.Equal(t, "\" value 8 \"", directive.List[7].String)
		assert.Equal(t, "' value 9 '", directive.List[8].String)
		assert.Equal(t, DirectiveTypeText, directive.List[9].Type)
		assert.Equal(t, 1, len(directive.List[9].Text))
		assert.Equal(t, "value '\" 10", directive.List[9].Text[0])
		assert.Equal(t, "And Fred said 'yabba dabba doo!' to Barney.", directive.List[10].String)
	})
}

func TestDictionary(t *testing.T) {
	t.Run("dict_01", func(t *testing.T) {
		expect := &Directive{
			Type: DirectiveTypeDictionary,
			Dictionary: map[string]*Directive{
				"key1": &Directive{
					Type:   DirectiveTypeString,
					String: "",
				},
				"key2": &Directive{
					Type:   DirectiveTypeString,
					String: "",
				},
			},
		}
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_01/load_in.nt")

		directive := &Directive{}
		directive.Parse(dat)
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

		err := directive.Parse([]byte("key\n: value"))
		assert.Equal(t, RootStringError, err)
	})

	t.Run("dict_03", func(t *testing.T) {
		// json input
	})

	t.Run("dict_04", func(t *testing.T) {
		// empty json object
	})
	t.Run("dict_05", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_05/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, RootLevelHasIndentError, err)
	})
	t.Run("dict_06", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_06/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, StringHasChildError, err)
	})
	t.Run("dict_07", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_07/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, TabInIndentationError, err)
	})
	t.Run("dict_08", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_08/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentTypesOnTheSameLevelError, err)
	})
	t.Run("dict_09", func(t *testing.T) {
		// differencing types on the same level
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_09/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, StringHasChildError, err)
	})
	t.Run("dict_10", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_10/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentLevelOnSameChildError, err)
	})
	t.Run("dict_11", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_11/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentLevelOnSameChildError, err)
	})
	t.Run("dict_12", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_12/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DifferentTypesOnTheSameLevelError, err)
	})
	t.Run("dict_13", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_13/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, StringWithNewLineError, err)
	})
	t.Run("dict_14", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_14/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, DictionaryDuplicateKeyError, err)
	})
	t.Run("dict_15", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_15/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Equal(t, TabInIndentationError, err)
	})
	t.Run("dict_16", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_16/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)

		_, exists := directive.Dictionary["key:"]

		assert.Nil(t, err)
		assert.True(t, exists)
	})
	t.Run("dict_17", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_17/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
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
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_18/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)
	})
	t.Run("dict_19", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_19/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)
	})
	t.Run("dict_20", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_20/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
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
	t.Run("dict_23", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/dict_23/load_in.nt")

		directive := &Directive{}

		err := directive.Parse(dat)
		assert.Nil(t, err)

		var value *Directive
		value, _ = directive.Dictionary["key1"]
		assert.Equal(t, "value 1", value.String)
		value, _ = directive.Dictionary["key2"]
		assert.Equal(t, "value 2", value.String)
		value, _ = directive.Dictionary["key 3"]
		assert.Equal(t, "value 3", value.String)
		value, _ = directive.Dictionary["key 4"]
		assert.Equal(t, "value 4", value.String)
		value, _ = directive.Dictionary["key5"]
		assert.Equal(t, "", value.String)
		value, _ = directive.Dictionary["key6"]
		assert.Equal(t, "", value.String)
		value, _ = directive.Dictionary[" key7 "]
		assert.Equal(t, "value 7", value.String)
		value, _ = directive.Dictionary[" key8 "]
		assert.Equal(t, "value 8", value.String)
		value, _ = directive.Dictionary[" ' key9 ' "]
		assert.Equal(t, "value 9", value.String)
		value, _ = directive.Dictionary[" \" key10 \" "]
		assert.Equal(t, "value 10", value.String)
		//value, _ = directive.Dictionary[" \" key11: \" "]
		//assert.Equal(t, "value 11", value.String)
		//value, _ = directive.Dictionary[" \" key12 : \" "]
		//assert.Equal(t, "value 12", value.String)
		//value, _ = directive.Dictionary[" \" key13: "]
		//assert.Equal(t, "value 13", value.String)
		//value, _ = directive.Dictionary[" \" key14 : "]
		//assert.Equal(t, "value 14", value.String)
		//value, _ = directive.Dictionary[" ' key15': "]
		//assert.Equal(t, "value 15", value.String)
		//value, _ = directive.Dictionary[" ' key16' : "]
		//assert.Equal(t, "value 16", value.String)
		value, _ = directive.Dictionary[""]
		assert.Equal(t, "value 17", value.String)
		//value, _ = directive.Dictionary[" ' key18\"' : "]
		//assert.Equal(t, "value 18", value.String)
		//value, _ = directive.Dictionary[" \" key19'\" : "]
		//assert.Equal(t, "value 19", value.String)
	})
	t.Run("dict_24", func(t *testing.T) {
		// json input
	})
}

func TestEmpty(t *testing.T) {
	t.Run("empty_1", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile(TestCasePath + "/empty_1/load_in.nt")

		directive := &Directive{}

		var err error

		err = directive.Parse(dat)
		assert.Equal(t, EmptyDataError, err)

		err = directive.Parse([]byte("\n  \n"))
		assert.Equal(t, EmptyDataError, err)
	})
}

func TestHolistic(t *testing.T) {
	getValueWithAssert := func(t *testing.T, dict map[string]*Directive, key string) *Directive {
		directive, exists := dict[key]
		assert.True(t, exists)
		return directive
	}

	t.Run("holistic_1", func(t *testing.T) {
		dat, _ := ioutil.ReadFile(TestCasePath + "/holistic_1/load_in.nt")

		directive := &Directive{}
		err := directive.Parse(dat)

		t.Run("should Parse successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should Parse with collect structure", func(t *testing.T) {
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
		dat, _ := ioutil.ReadFile(TestCasePath + "/holistic_2/load_in.nt")

		directive := &Directive{}
		err := directive.Parse(dat)

		t.Run("should Parse successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should Parse with collect structure", func(t *testing.T) {
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
		dat, _ := ioutil.ReadFile(TestCasePath + "/holistic_3/load_in.nt")

		directive := &Directive{}
		err := directive.Parse(dat)

		t.Run("should Parse successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should Parse with collect structure", func(t *testing.T) {
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
		dat, _ := ioutil.ReadFile(TestCasePath + "/holistic_4/load_in.nt")

		directive := &Directive{}
		err := directive.Parse(dat)

		t.Run("should Parse successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should Parse with collect structure", func(t *testing.T) {
			assert.Equal(t, DirectiveTypeDictionary, directive.Type)

			rootDict := directive.Dictionary

			keyOutput := getValueWithAssert(t, rootDict, "output current")
			assert.Equal(t, DirectiveTypeString, keyOutput.Type)
			assert.Equal(t, "out", keyOutput.String)

			keyDescription := getValueWithAssert(t, rootDict, "description")
			assert.Equal(t, DirectiveTypeString, keyDescription.Type)
			assert.Equal(t, "Output current", keyDescription.String)

			keyRange := getValueWithAssert(t, rootDict, "range")
			assert.Equal(t, DirectiveTypeString, keyRange.Type)
			assert.Equal(t, "V(gnda) + 0.5V < V < V(vdda) - 0.5V; -500μA <= I <= 500μA", keyRange.String)

			keyBehavior := getValueWithAssert(t, rootDict, "behavior")
			assert.Equal(t, DirectiveTypeText, keyBehavior.Type)
			assert.Equal(t, 4, len(keyBehavior.Text))
			assert.Equal(t, "current:\n", keyBehavior.Text[0])
			assert.Equal(t, "    I = On*Iout;\n", keyBehavior.Text[1])
			assert.Equal(t, "    IoutMeas=I with prail=vdda; nrail=gnda", keyBehavior.Text[3])

			keyNominal := getValueWithAssert(t, rootDict, "nominal")
			assert.Equal(t, DirectiveTypeString, keyNominal.Type)
			assert.Equal(t, "V=1.25V+1Ω*I", keyNominal.String)
		})
	})

	t.Run("holistic_5", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile(TestCasePath + "/holistic_5/load_in.nt")

		directive := &Directive{}
		err := directive.Parse(dat)

		t.Run("should Parse successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should Parse with collect structure", func(t *testing.T) {
			assert.Equal(t, DirectiveTypeDictionary, directive.Type)

			rootDict := directive.Dictionary

			keyPresident := getValueWithAssert(t, rootDict, "president")
			assert.Equal(t, DirectiveTypeDictionary, keyPresident.Type)

			dictPresident := keyPresident.Dictionary
			{
				keyName := getValueWithAssert(t, dictPresident, "name")
				assert.Equal(t, DirectiveTypeString, keyName.Type)
				assert.Equal(t, "Katheryn McDaniel", keyName.String)

				keyAddress := getValueWithAssert(t, dictPresident, "address")
				assert.Equal(t, DirectiveTypeText, keyAddress.Type)
				assert.Equal(t, 2, len(keyAddress.Text))
				assert.Equal(t, "138 Almond Street\n", keyAddress.Text[0])
				assert.Equal(t, "Topeka, Kansas 20697", keyAddress.Text[1])

				keyPhone := getValueWithAssert(t, dictPresident, "phone")
				assert.Equal(t, DirectiveTypeDictionary, keyPhone.Type)

				dictPhone := keyPhone.Dictionary
				{
					keyCell := getValueWithAssert(t, dictPhone, "cell")
					assert.Equal(t, DirectiveTypeString, keyCell.Type)
					assert.Equal(t, "1-210-835-5297", keyCell.String)

					keyHome := getValueWithAssert(t, dictPhone, "home")
					assert.Equal(t, DirectiveTypeString, keyHome.Type)
					assert.Equal(t, "1-210-478-8470", keyHome.String)
				}

				keyEmail := getValueWithAssert(t, dictPresident, "email")
				assert.Equal(t, DirectiveTypeString, keyEmail.Type)
				assert.Equal(t, "KateMcD@aol.com", keyEmail.String)

				keyKids := getValueWithAssert(t, dictPresident, "kids")
				assert.Equal(t, DirectiveTypeList, keyKids.Type)
				assert.Equal(t, 2, len(keyKids.List))

				listKids := keyKids.List
				{
					assert.Equal(t, DirectiveTypeString, listKids[0].Type)
					assert.Equal(t, "Joanie", listKids[0].String)

					assert.Equal(t, DirectiveTypeString, listKids[1].Type)
					assert.Equal(t, "Terrance", listKids[1].String)
				}
			}
		})
	})

	t.Run("holistic_6", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile(TestCasePath + "/holistic_6/load_in.nt")

		directive := &Directive{}
		err := directive.Parse(dat)

		t.Run("should Parse successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should Parse with collect structure", func(t *testing.T) {
			assert.Equal(t, DirectiveTypeDictionary, directive.Type)

			rootDict := directive.Dictionary

			keyVp := getValueWithAssert(t, rootDict, "vice president")
			assert.Equal(t, DirectiveTypeDictionary, keyVp.Type)

			dictVp := keyVp.Dictionary
			{
				keyName := getValueWithAssert(t, dictVp, "name")
				assert.Equal(t, DirectiveTypeString, keyName.Type)
				assert.Equal(t, "Margaret Hodge", keyName.String)

				keyAddress := getValueWithAssert(t, dictVp, "address")
				assert.Equal(t, DirectiveTypeText, keyAddress.Type)
				assert.Equal(t, 2, len(keyAddress.Text))
				assert.Equal(t, "2586 Marigold Lane\n", keyAddress.Text[0])
				assert.Equal(t, "Topeka, Kansas 20682", keyAddress.Text[1])

				keyPhone := getValueWithAssert(t, dictVp, "phone")
				assert.Equal(t, DirectiveTypeString, keyPhone.Type)
				assert.Equal(t, "1-470-974-0398", keyPhone.String)

				keyEmail := getValueWithAssert(t, dictVp, "email")
				assert.Equal(t, DirectiveTypeString, keyEmail.Type)
				assert.Equal(t, "margarett.hodge@ku.edu", keyEmail.String)

				keyKids := getValueWithAssert(t, dictVp, "kids")
				assert.Equal(t, DirectiveTypeList, keyKids.Type)
				assert.Equal(t, 3, len(keyKids.List))

				listKids := keyKids.List
				{
					assert.Equal(t, DirectiveTypeString, listKids[0].Type)
					assert.Equal(t, "Arnie", listKids[0].String)

					assert.Equal(t, DirectiveTypeString, listKids[1].Type)
					assert.Equal(t, "Zach", listKids[1].String)

					assert.Equal(t, DirectiveTypeString, listKids[2].Type)
					assert.Equal(t, "Maggie", listKids[2].String)
				}
			}
		})
	})

	t.Run("holistic_7", func(t *testing.T) {
		// empty content should be null
		dat, _ := ioutil.ReadFile(TestCasePath + "/holistic_7/load_in.nt")

		directive := &Directive{}
		err := directive.Parse(dat)

		t.Run("should Parse successfully", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should Parse with collect structure", func(t *testing.T) {
			assert.Equal(t, DirectiveTypeDictionary, directive.Type)

			rootDict := directive.Dictionary

			keyTreasurer := getValueWithAssert(t, rootDict, "treasurer")
			assert.Equal(t, DirectiveTypeDictionary, keyTreasurer.Type)

			dictTreasurer := keyTreasurer.Dictionary
			{
				keyName := getValueWithAssert(t, dictTreasurer, "name")
				assert.Equal(t, DirectiveTypeString, keyName.Type)
				assert.Equal(t, "	       Fumiko\tPurvis    \t", keyName.String)

				keyAddress := getValueWithAssert(t, dictTreasurer, "address")
				assert.Equal(t, DirectiveTypeText, keyAddress.Type)
				assert.Equal(t, 2, len(keyAddress.Text))
				assert.Equal(t, "\t 3636 Buffalo Ave \t\n", keyAddress.Text[0])
				assert.Equal(t, "\t Topeka, Kansas 20692\t ", keyAddress.Text[1])
			}
		})
	})
}
