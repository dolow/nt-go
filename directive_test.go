package ntgo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {

	var data []byte

	subject := func() (*Directive, *DirectiveMarshalError) {
		directive := &Directive{}
		err := directive.Marshal(data)
		return directive, err
	}

	t.Run("string", func(t *testing.T) {

		expect := "plain text"
		t.Run("regular string", func(t *testing.T) {
			data = []byte(expect)

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err.error)
			})
		})

		t.Run("string start with space", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  %s", expect))

			t.Run("should cause RootLevelHasIndentError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootLevelHasIndentError, err.error)
			})
		})

		t.Run("string start with line break", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n%s", expect))

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err.error)
			})
		})

		t.Run("string start with line break with forward spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  \n%s", expect))

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err.error)
			})
		})

		t.Run("string start with line break and second line starts with spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n  %s", expect))

			t.Run("should cause RootLevelHasIndentError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootLevelHasIndentError, err.error)
			})
		})

		t.Run("string with comment symbol (#)", func(t *testing.T) {
			expect = "plain text # it is not comment"
			data = []byte(expect)

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err.error)
			})
		})

		t.Run("string ending with line break", func(t *testing.T) {
			expect = "plain text"
			data = []byte(fmt.Sprintf("%s\n", expect))

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err.error)
			})
		})
	})

	t.Run("text", func(t *testing.T) {

		expect := []string{"multiple\n", "line of text"}

		t.Run("regular text", func(t *testing.T) {
			data = []byte(fmt.Sprintf("> %s> %s", expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeText", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeText, directive.Type)
			})
			t.Run("Text should be input string array", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.Text)
			})
		})

		t.Run("texts start with spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  > %s  > %s", expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeText", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeText, directive.Type)
			})
			t.Run("Text should be input string array", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.Text)
			})
		})

		t.Run("texts start with blank line", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n  > %s  > %s", expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeText", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeText, directive.Type)
			})
			t.Run("Text should be input string array", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.Text)
			})
		})

		t.Run("texts start with blank line with forwarding spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  \n  > %s  > %s", expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeText", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeText, directive.Type)
			})
			t.Run("Text should be input string array", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.Text)
			})
		})

		// irregular case
		/*
			t.Run("texts start with different indentations", func(t *testing.T) {
			})
		*/
	})

	t.Run("list", func(t *testing.T) {
		t.Run("string elements", func(t *testing.T) {
			expect := []string{"string", "elements"}
			data = []byte(fmt.Sprintf(
				`- %s
- %s`,
				expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List should contain directives with DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()
				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(directive.List))

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]
					assert.Equal(t, DirectiveTypeString, element.Type)
					assert.Equal(t, expect[i], element.String)
				}
			})
		})

		t.Run("string elements with unbalanced forwarding spaces", func(t *testing.T) {
			expect := []string{"   string", "elements"}
			data = []byte(fmt.Sprintf(
				`- %s
- %s`,
				expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List should contain string with leading space", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(directive.List))

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]
					assert.Equal(t, DirectiveTypeString, element.Type)
					assert.Equal(t, expect[i], element.String)
				}
			})
		})

		t.Run("text elements", func(t *testing.T) {
			expect := [][]string{
				[]string{"aaaa", "bbbb"},
				[]string{"cccc", "dddd"},
			}
			data = []byte(fmt.Sprintf(
				`-
  > %s
  > %s
-
  > %s
  > %s`,
				expect[0][0], expect[0][1],
				expect[1][0], expect[1][1],
			))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List should contain directives with DirectiveTypeText", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(directive.List))

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]
					assert.Equal(t, DirectiveTypeText, element.Type)
					for j := 0; j < len(element.Text); j++ {
						e := expect[i][j]
						if j != len(element.Text)-1 {
							e = fmt.Sprintf("%s\n", e)
						}
						assert.Equal(t, e, element.Text[j])
					}
				}
			})
		})

		t.Run("text elements with unbalanced indentations", func(t *testing.T) {
			expect := [][]string{
				[]string{"aaaa", "bbbb"},
				[]string{"cccc", "dddd"},
			}
			data = []byte(fmt.Sprintf(
				`-
  > %s
  > %s
-
     > %s
     > %s`,
				expect[0][0], expect[0][1],
				expect[1][0], expect[1][1],
			))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("Text directives should be in the same depth", func(t *testing.T) {
				directive, err := subject()
				assert.Nil(t, err)

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]
					assert.Equal(t, directive.Depth+1, element.Depth)
				}
			})
		})

		t.Run("list elements", func(t *testing.T) {
			expect := [][]string{
				[]string{"aaaa", "bbbb"},
				[]string{"cccc", "dddd"},
			}
			data = []byte(fmt.Sprintf(
				`-
  - %s
  - %s
-
  - %s
  - %s`,
				expect[0][0], expect[0][1],
				expect[1][0], expect[1][1],
			))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List should contain directives with DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(directive.List))

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]
					assert.Equal(t, DirectiveTypeList, element.Type)
					for j := 0; j < len(element.List); j++ {
						nestedElement := element.List[j]
						assert.Equal(t, DirectiveTypeString, nestedElement.Type)
						assert.Equal(t, expect[i][j], nestedElement.String)
					}
				}
			})
		})

		t.Run("list elements with unbalanced indentations", func(t *testing.T) {
			expect := [][]string{
				[]string{"aaaa", "bbbb"},
				[]string{"cccc", "dddd"},
			}
			data = []byte(fmt.Sprintf(
				`-
     - %s
     - %s
-
  - %s
  - %s`,
				expect[0][0], expect[0][1],
				expect[1][0], expect[1][1],
			))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List directives should be in the same depth", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(directive.List))

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]
					assert.Equal(t, directive.Depth+1, element.Depth)
					for j := 0; j < len(element.List); j++ {
						assert.Equal(t, element.Depth+1, element.List[j].Depth)
					}
				}
			})
		})

		t.Run("dictionary string elements", func(t *testing.T) {
			expect := [][][]string{
				[][]string{
					{"key1", "val1"},
					{"key2", "val2"},
				},
				[][]string{
					{"key3", "val3"},
					{"key4", "val4"},
				},
			}
			data = []byte(fmt.Sprintf(
				`-
  %s: %s
  %s: %s
-
  %s: %s
  %s: %s`,
				expect[0][0][0], expect[0][0][1],
				expect[0][1][0], expect[0][1][1],
				expect[1][0][0], expect[1][0][1],
				expect[1][1][0], expect[1][1][1],
			))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List should contain directives with DirectiveTypeDictionary", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(directive.List))

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]

					assert.Equal(t, DirectiveTypeDictionary, element.Type)
					assert.Equal(t, expect[i][0][1], element.Dictionary[expect[i][0][0]].String)
					assert.Equal(t, expect[i][1][1], element.Dictionary[expect[i][1][0]].String)
				}
			})
		})

		t.Run("dictionary string elements with unbalanced spaces", func(t *testing.T) {
			expect := [][][]string{
				[][]string{
					{"key1", "   val1   "},
					{"key2", "val2"},
				},
				[][]string{
					{"key3", "val3"},
					{"key4", "\tval4\t"},
				},
			}
			data = []byte(fmt.Sprintf(
				`-
  %s: %s
  %s: %s
-
     %s: %s
     %s: %s`,
				expect[0][0][0], expect[0][0][1],
				expect[0][1][0], expect[0][1][1],
				expect[1][0][0], expect[1][0][1],
				expect[1][1][0], expect[1][1][1],
			))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("Dictionary elements should be in the same depth", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]

					assert.Equal(t, element.Depth+1, element.Dictionary[expect[i][0][0]].Depth)
				}
			})

			t.Run("Dictionary elements value string should contain leading and trailing spaces and tabs", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]

					assert.Equal(t, expect[i][0][1], element.Dictionary[expect[i][0][0]].String)
					assert.Equal(t, expect[i][1][1], element.Dictionary[expect[i][1][0]].String)
				}
			})
		})
	})

	t.Run("dictionary", func(t *testing.T) {
		t.Run("string elements", func(t *testing.T) {
			expectKey := []string{"key1", "key2"}
			expectValue := []string{"value1", "value2"}

			data = []byte(fmt.Sprintf(
				`%s: %s
%s: %s`,
				expectKey[0], expectValue[0],
				expectKey[1], expectValue[1],
			))

			t.Run("Type should be DirectiveTypeDictionary", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeDictionary, directive.Type)
			})
			t.Run("Dictionary should contain directives with DirectiveTypeString and certain keys", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expectKey), len(directive.Dictionary))

				assert.Equal(t, expectValue[0], directive.Dictionary[expectKey[0]].String)
				assert.Equal(t, expectValue[1], directive.Dictionary[expectKey[1]].String)
			})
		})
	})
}

func TestUnmarshal(t *testing.T) {

	var data []byte
	var indentSize int
	var depth int

	expect := func() string { return "" }

	resetCondition := func() {
		indentSize = 2
		depth = 0
	}

	subject := func() string {
		directive := &Directive{
			IndentSize: indentSize,
			Depth:      depth,
		}
		directive.Marshal(data)
		return directive.Unmarshal()
	}

	t.Run("text", func(t *testing.T) {
		line1 := []byte("> aaaa\n")
		line2 := []byte("> bbbb")
		data = []byte(fmt.Sprintf("%s%s", string(line1), string(line2)))

		t.Run("Depth is 0", func(t *testing.T) {
			defer resetCondition()

			depth = 0

			expect = func() string { return fmt.Sprintf("%s%s", string(line1), string(line2)) }

			t.Run("should return text with no indentation", func(t *testing.T) {
				assert.Equal(t, expect(), subject())
			})
		})

		t.Run("Depth and IndentSize is larger than 0", func(t *testing.T) {
			defer resetCondition()

			depth = 2
			indentSize = 4

			expect = func() string {
				indent := fmt.Sprintf("%*s", depth*indentSize, "")
				return fmt.Sprintf("%s%s%s%s", indent, string(line1), indent, string(line2))
			}

			t.Run("should return text with indentation", func(t *testing.T) {
				assert.Equal(t, expect(), subject())
			})
		})
	})

	t.Run("list", func(t *testing.T) {
		line1 := []byte("- aaaa\n")
		line2 := []byte("- bbbb")
		data = []byte(fmt.Sprintf("%s%s", string(line1), string(line2)))

		t.Run("Depth is 0", func(t *testing.T) {
			defer resetCondition()

			depth = 0

			expect = func() string { return fmt.Sprintf("%s%s", string(line1), string(line2)) }

			t.Run("should return text with no indentation", func(t *testing.T) {
				assert.Equal(t, expect(), subject())
			})
		})

		t.Run("Depth and IndentSize is larger than 0", func(t *testing.T) {
			defer resetCondition()

			depth = 2
			indentSize = 4

			expect = func() string {
				indent := fmt.Sprintf("%*s", depth*indentSize, "")
				return fmt.Sprintf("%s%s%s%s", indent, string(line1), indent, string(line2))
			}

			t.Run("should return text with indentation", func(t *testing.T) {
				assert.Equal(t, expect(), subject())
			})
		})
	})

	t.Run("dictionary", func(t *testing.T) {
		line1 := []byte("key1: value1")
		line2 := []byte("key2: value2")
		data = []byte(fmt.Sprintf("%s\n%s", string(line1), string(line2)))

		t.Run("Depth is 0", func(t *testing.T) {
			defer resetCondition()

			depth = 0

			expect1 := func() string { return fmt.Sprintf("%s\n%s", string(line1), string(line2)) }
			expect2 := func() string { return fmt.Sprintf("%s\n%s", string(line2), string(line1)) }

			t.Run("should return text with no indentation", func(t *testing.T) {
				// dictionary is unordered
				result := subject()
				assert.True(t, expect1() == result || expect2() == result)
			})
		})

		t.Run("Depth and IndentSize is larger than 0", func(t *testing.T) {
			defer resetCondition()

			depth = 2
			indentSize = 4

			expect1 := func() string {
				indent := fmt.Sprintf("%*s", depth*indentSize, "")
				return fmt.Sprintf("%s%s\n%s%s", indent, string(line1), indent, string(line2))
			}

			expect2 := func() string {
				indent := fmt.Sprintf("%*s", depth*indentSize, "")
				return fmt.Sprintf("%s%s\n%s%s", indent, string(line2), indent, string(line1))
			}

			t.Run("should return text with indentation", func(t *testing.T) {
				// dictionary is unordered
				result := subject()
				assert.True(t, expect1() == result || expect2() == result)
			})
		})
	})
}

func TestDetectKeyBytes(t *testing.T) {

	cases := [][]string{
		[]string{`-#:'>: -#:">:`, `-#:'>`},
		[]string{`-#:">: -#:'>:`, `-#:">`},
		[]string{`-#'\'>:: -#"\">::`, `-#'\'>:`},
		[]string{`-#"\">:: -#'\'>::`, `-#"\">:`},
		[]string{`:-#:'>: :-#:">:`, `:-#:'>`},
		[]string{`:-#:">: :-#:'>:`, `:-#:">`},
		[]string{`:-#'\'>:: :-#"\">::`, `:-#'\'>:`},
		[]string{`:-#"\">:: :-#'\'>::`, `:-#"\">:`},
		[]string{`>:-#:'>: >:-#:">:`, `>:-#:'>`},
		[]string{`>:-#:">: >:-#:'>:`, `>:-#:">`},
		[]string{`>:-#'\'>:: >:-#"\">::`, `>:-#'\'>:`},
		[]string{`>:-#"\">:: >:-#'\'>::`, `>:-#"\">:`},
		[]string{`:`, ``},
		[]string{`~!@#$%^&*()_+-1234567890{}[]|\;<>?,./: ~!@#$%^&*()_+-1234567890{}[]|\:;<>?,./`, `~!@#$%^&*()_+-1234567890{}[]|\;<>?,./`},
		[]string{`'- key 3': - value 3`, `'- key 3'`},       // not sanitize
		[]string{`'key 4: ': value 4:`, `'key 4: '`},        // not sanitize
		[]string{`'> key 5': > value 5`, `'> key 5'`},       // not sanitize
		[]string{`'# key 6': #value 6`, `'# key 6'`},        // not sanitize
		[]string{`': key 7': : value 7`, `': key 7'`},       // not sanitize
		[]string{`'" key 8 "': " value 8 "`, `'" key 8 "'`}, // not sanitize
		[]string{`"' key 9 '": ' value 9 '`, `"' key 9 '"`}, // not sanitize
		[]string{`key 10: value '" 10`, `key 10`},
		[]string{`key 11: And Fred said 'yabba dabba doo!' to Barney.`, `key 11`},
		[]string{`key " 12: value ' 12`, `key " 12`},
		[]string{`$€¥£₩₺₽₹ɃΞȄ: $€¥£₩₺₽₹ɃΞȄ`, `$€¥£₩₺₽₹ɃΞȄ`},
		[]string{`YZEPTGMKk_cmuµμnpfazy: YZEPTGMKk_cmuµμnpfazy`, `YZEPTGMKk_cmuµμnpfazy`},
		[]string{`a-zA-Z%√{us}{cur}][-^/()\w·⁻⁰¹²³⁴⁵⁶⁷⁸⁹°ÅΩƱΩ℧: a-zA-Z%√{us}{cur}][-^/()\w·⁻⁰¹²³⁴⁵⁶⁷⁸⁹°ÅΩƱΩ℧`, `a-zA-Z%√{us}{cur}][-^/()\w·⁻⁰¹²³⁴⁵⁶⁷⁸⁹°ÅΩƱΩ℧`},
	}

	for _, item := range cases {
		it := item[0]
		expect := []byte(item[1])

		t.Run(fmt.Sprintf(`key of %s should be %s`, it, expect), func(t *testing.T) {
			key, _ := detectKeyBytes([]byte(it))
			assert.Equal(t, expect, key)
		})
	}
}

func TestReadTextDirective(t *testing.T) {

	var directive *Directive

	var index int
	var firstLine []byte
	var buffer *bytes.Buffer

	var content []byte

	var bufferInitializer func() *bytes.Buffer

	prepare := func() {
		directive = &Directive{}
		bufferInitializer = func() *bytes.Buffer {
			return bytes.NewBuffer(content)
		}
	}

	condition := func() {}

	subject := func() ([]byte, *DirectiveMarshalError) {
		condition()
		buffer = bufferInitializer()
		return directive.readTextDirective(index, firstLine, buffer)
	}

	t.Run("when next directive appeared", func(t *testing.T) {
		condition = func() {
			index = 2
			firstLine = []byte("  > first line\n")
			content = []byte(
				`  > second line
  > third line
- list`)
		}

		t.Run("should return DirectiveMarshalError with NextDirectiveAppearedError", func(t *testing.T) {
			prepare()
			_, err := subject()
			assert.Equal(t, NextDirectiveAppearedError, err.error)
		})
		t.Run("should return first line of next different directive", func(t *testing.T) {
			prepare()
			nextLine, _ := subject()
			assert.Equal(t, []byte("- list"), nextLine)
		})
		t.Run("should Text slice ends with last line of text directive", func(t *testing.T) {
			prepare()
			subject()
			assert.Equal(t, 3, len(directive.Text))
			assert.Equal(t, "third line", directive.Text[2])
		})
	})

	t.Run("when eof occured", func(t *testing.T) {
		condition = func() {
			index = 0
			firstLine = []byte("> first line\n")
			content = []byte("> second line\n> third line")
		}

		t.Run("should return nil for error", func(t *testing.T) {
			prepare()
			_, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _ := subject()
			assert.Nil(t, nextLine)
		})
		t.Run("should Text slice ends with last line", func(t *testing.T) {
			prepare()
			subject()
			assert.Equal(t, 3, len(directive.Text))
			assert.Equal(t, "third line", directive.Text[2])
		})
	})
	t.Run("when text ends with text symbol(>)", func(t *testing.T) {
		condition = func() {
			index = 0
			firstLine = []byte("> first line\n")
			content = []byte(">")
		}
		t.Run("should return nil for error", func(t *testing.T) {
			prepare()
			_, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _ := subject()
			assert.Nil(t, nextLine)
		})
		t.Run("should add empty line to Text slice", func(t *testing.T) {
			prepare()
			subject()
			assert.Equal(t, 2, len(directive.Text))
			assert.Equal(t, "", directive.Text[1])
		})
	})
	t.Run("when text contains blank line on the head", func(t *testing.T) {
		condition = func() {
			index = 0
			firstLine = []byte("\n")
			content = []byte("> first line\n> second line")
		}
		t.Run("should return nil for error", func(t *testing.T) {
			prepare()
			_, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _ := subject()
			assert.Nil(t, nextLine)
		})
		t.Run("should add only meaningful lines to Text slice", func(t *testing.T) {
			prepare()
			subject()
			assert.Equal(t, 2, len(directive.Text))
			assert.Equal(t, "first line\n", directive.Text[0])
			assert.Equal(t, "second line", directive.Text[1])
		})
	})
	t.Run("when text contains blank line on the middle", func(t *testing.T) {
		condition = func() {
			index = 0
			firstLine = []byte("> first line\n")
			content = []byte("\n")
			content = []byte("> second line")
		}
		t.Run("should return nil for error", func(t *testing.T) {
			prepare()
			_, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _ := subject()
			assert.Nil(t, nextLine)
		})
		t.Run("should add only meaningful lines to Text slice", func(t *testing.T) {
			prepare()
			subject()
			assert.Equal(t, 2, len(directive.Text))
			assert.Equal(t, "first line\n", directive.Text[0])
			assert.Equal(t, "second line", directive.Text[1])
		})
	})
	t.Run("when text ends with blank line", func(t *testing.T) {
		condition = func() {
			index = 0
			firstLine = []byte("> first line\n")
			content = []byte("> second line\n")
		}
		t.Run("should return nil for error", func(t *testing.T) {
			prepare()
			_, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _ := subject()
			assert.Nil(t, nextLine)
		})
		t.Run("should add only meaningful lines to Text slice", func(t *testing.T) {
			prepare()
			subject()
			assert.Equal(t, 2, len(directive.Text))
			assert.Equal(t, "first line\n", directive.Text[0])
			assert.Equal(t, "second line", directive.Text[1])
		})
	})

	t.Run("irregulars", func(t *testing.T) {
		t.Run("when Type is already defined", func(t *testing.T) {
			condition = func() {
				index = 0
				firstLine = []byte("> first line\n")
				content = []byte("> second line")

				directive.Type = DirectiveTypeText
			}

			t.Run("should return DirectiveMarshalError with DifferentTypesOnTheSameLevelError", func(t *testing.T) {
				prepare()
				_, err := subject()
				assert.Equal(t, DifferentTypesOnTheSameLevelError, err.error)
			})

			t.Run("should return nil for nextLine", func(t *testing.T) {
				prepare()
				nextLine, _ := subject()
				assert.Nil(t, nextLine)
			})
		})

		t.Run("when content consists of lines with different indentations", func(t *testing.T) {

			index = 2
			firstLine = []byte("  > first line\n")

			t.Run("when following line is deeper", func(t *testing.T) {

				contentPh := "    %s"

				t.Run("when following line is text", func(t *testing.T) {

					condition = func() {
						content = []byte(fmt.Sprintf(contentPh, "> text"))
					}

					t.Run("should return DirectiveMarshalError with DifferentLevelOnSameChildError", func(t *testing.T) {
						prepare()
						_, err := subject()
						assert.Equal(t, DifferentLevelOnSameChildError, err.error)
					})

					t.Run("should return nil for nextLine", func(t *testing.T) {
						prepare()
						nextLine, _ := subject()
						assert.Nil(t, nextLine)
					})
				})

				t.Run("when following line is not text", func(t *testing.T) {
					condition = func() {
						content = []byte(fmt.Sprintf(contentPh, "- list"))
					}

					t.Run("should return TextHasChildError with DifferentLevelOnSameChildError", func(t *testing.T) {
						prepare()
						_, err := subject()
						assert.Equal(t, TextHasChildError, err.error)
					})

					t.Run("should return nil for nextLine", func(t *testing.T) {
						prepare()
						nextLine, _ := subject()
						assert.Nil(t, nextLine)
					})
				})
			})

			t.Run("when following line is shallower", func(t *testing.T) {

				contentPh := "%s"

				t.Run("when following line is text", func(t *testing.T) {

					condition = func() {
						content = []byte(fmt.Sprintf(contentPh, "> text"))
					}

					t.Run("should return DirectiveMarshalError with DifferentLevelOnSameChildError", func(t *testing.T) {
						prepare()
						_, err := subject()
						assert.Equal(t, DifferentLevelOnSameChildError, err.error)
					})

					t.Run("should return nil for nextLine", func(t *testing.T) {
						prepare()
						nextLine, _ := subject()
						assert.Nil(t, nextLine)
					})
				})

				t.Run("when following line is not text", func(t *testing.T) {

					condition = func() {
						content = []byte(fmt.Sprintf(contentPh, "- list"))
					}

					t.Run("should return DirectiveMarshalError with NextDirectiveAppearedError", func(t *testing.T) {
						prepare()
						_, err := subject()
						assert.Equal(t, NextDirectiveAppearedError, err.error)
					})

					t.Run("should return first line of next directive", func(t *testing.T) {
						prepare()
						nextLine, _ := subject()
						assert.Equal(t, []byte("- list"), nextLine)
					})
				})
			})
		})

		t.Run("when content consists of lines with different directives", func(t *testing.T) {
			condition = func() {
				firstLine = []byte("> first line\n")
				content = []byte("- list")
			}

			t.Run("should return DirectiveMarshalError with DifferentLevelOnSameChildError", func(t *testing.T) {
				prepare()
				_, err := subject()
				assert.Equal(t, DifferentLevelOnSameChildError, err.error)
			})

			t.Run("should return nil for next line", func(t *testing.T) {
				prepare()
				nextLine, _ := subject()
				assert.Nil(t, nextLine)
			})
		})
	})
}

func TestRemoveTrailingLineBreak(t *testing.T) {
	t.Run("should remove trailing line break", func(t *testing.T) {
		str := "hello world\n"
		removeTrailingLineBreak(&str)
		assert.Equal(t, "hello world", str)
	})

	t.Run("should remove only single line break", func(t *testing.T) {
		t.Run("consequent line breaks", func(t *testing.T) {
			str := "hello world\n\n"
			removeTrailingLineBreak(&str)
			assert.Equal(t, "hello world\n", str)
		})
		t.Run("multilines", func(t *testing.T) {
			str := "hello\nworld\n"
			removeTrailingLineBreak(&str)
			assert.Equal(t, "hello\nworld", str)
		})
	})

	t.Run("should not remove character if last character is not a line break", func(t *testing.T) {
		str := "hello world"
		removeTrailingLineBreak(&str)
		assert.Equal(t, "hello world", str)
	})
}
