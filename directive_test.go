package ntgo

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	var data []byte

	subject := func() (*Directive, error) {
		directive := &Directive{}
		err := directive.Parse(data)
		return directive, err
	}

	t.Run("string", func(t *testing.T) {

		expect := "plain text"
		t.Run("regular string", func(t *testing.T) {
			data = []byte(expect)

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err)
			})
		})

		t.Run("string start with space", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  %s", expect))

			t.Run("should cause RootLevelHasIndentError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootLevelHasIndentError, err)
			})
		})

		t.Run("string start with line break", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n%s", expect))

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err)
			})
		})

		t.Run("string start with line break with forward spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  \n%s", expect))

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err)
			})
		})

		t.Run("string start with line break and second line starts with spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n  %s", expect))

			t.Run("should cause RootLevelHasIndentError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootLevelHasIndentError, err)
			})
		})

		t.Run("string with comment symbol (#)", func(t *testing.T) {
			expect = "plain text # it is not comment"
			data = []byte(expect)

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err)
			})
		})

		t.Run("string ending with line break", func(t *testing.T) {
			expect = "plain text"
			data = []byte(fmt.Sprintf("%s\n", expect))

			t.Run("should cause RootStringError", func(t *testing.T) {
				_, err := subject()
				assert.NotNil(t, err)
				assert.Equal(t, RootStringError, err)
			})
		})
	})

	t.Run("text", func(t *testing.T) {

		expect := MultiLineText{"multiple\n", "line of text"}

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

			t.Run("should return RootLevelHasIndentError", func(t *testing.T) {
				_, err := subject()
				assert.Equal(t, RootLevelHasIndentError, err)
			})
		})

		t.Run("texts start with spaces after blank line", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n  > %s  > %s", expect[0], expect[1]))

			t.Run("should return RootLevelHasIndentError", func(t *testing.T) {
				_, err := subject()
				assert.Equal(t, RootLevelHasIndentError, err)
			})
		})

		t.Run("texts start with spaces after blank line with spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  \n  > %s  > %s", expect[0], expect[1]))

			t.Run("should return RootLevelHasIndentError", func(t *testing.T) {
				_, err := subject()
				assert.Equal(t, RootLevelHasIndentError, err)
			})
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("string elements", func(t *testing.T) {
			expect := MultiLineText{"string", "elements"}
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
			expect := MultiLineText{"   string", "elements"}
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
			expect := []MultiLineText{
				MultiLineText{"aaaa", "bbbb"},
				MultiLineText{"cccc", "dddd"},
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
			expect := []MultiLineText{
				MultiLineText{"aaaa", "bbbb"},
				MultiLineText{"cccc", "dddd"},
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
			expect := []MultiLineText{
				MultiLineText{"aaaa", "bbbb"},
				MultiLineText{"cccc", "dddd"},
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
			expect := []MultiLineText{
				MultiLineText{"aaaa", "bbbb"},
				MultiLineText{"cccc", "dddd"},
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
			expect := [][]MultiLineText{
				[]MultiLineText{
					{"key1", "val1"},
					{"key2", "val2"},
				},
				[]MultiLineText{
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
			expect := [][]MultiLineText{
				[]MultiLineText{
					{"key1", "   val1   "},
					{"key2", "val2"},
				},
				[]MultiLineText{
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

	t.Run("line breaks", func(t *testing.T) {
		t.Run("cr", func(t *testing.T) {
			data = []byte("- elem1\r- elem2")

			t.Run("should parse regulary", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
				assert.Equal(t, "elem1", directive.List[0].String)
				assert.Equal(t, "elem2", directive.List[1].String)
			})
		})
		t.Run("lf", func(t *testing.T) {
			data = []byte("- elem1\n- elem2")

			t.Run("should parse regulary", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
				assert.Equal(t, "elem1", directive.List[0].String)
				assert.Equal(t, "elem2", directive.List[1].String)
			})
		})
		t.Run("crlf", func(t *testing.T) {
			data = []byte("- elem1\r\n- elem2")

			t.Run("should parse regulary", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
				assert.Equal(t, "elem1", directive.List[0].String)
				assert.Equal(t, "elem2", directive.List[1].String)
			})
		})

		t.Run("mixed", func(t *testing.T) {
			data = []byte("- elem1\r\n- elem2\r- elem3\n- elem4")

			t.Run("should parse regulary", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
				assert.Equal(t, "elem1", directive.List[0].String)
				assert.Equal(t, "elem2", directive.List[1].String)
				assert.Equal(t, "elem3", directive.List[2].String)
				assert.Equal(t, "elem4", directive.List[3].String)
			})
		})

		t.Run("mixed as text content", func(t *testing.T) {
			data = []byte("text:\n  > line1\r\n  > line2\r  > line3\n  > line4\r\n  > line5")

			t.Run("should parse regulary", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)

				text := directive.Dictionary["text"].Text
				// TODO: mixed line break code in text content
				assert.Equal(t, "line1\r\n", text[0])
				assert.Equal(t, "line2\r", text[1])
				assert.Equal(t, "line3\n", text[2])
				assert.Equal(t, "line4\r\n", text[3])
				assert.Equal(t, "line5", text[4])
			})
		})
	})
}

func TestToString(t *testing.T) {

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
		directive.Parse(data)
		return directive.ToString()
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

func TestDetectDirectiveType(t *testing.T) {
	var data []byte
	subject := func() (DirectiveType, int, error) {
		return detectDirectiveType(data)
	}

	t.Run("when empty data given", func(t *testing.T) {
		data = []byte("  ")

		t.Run("should return DirectiveTypeUnknown", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, DirectiveTypeUnknown, typ)
		})

		t.Run("should return NotFoundIndex", func(t *testing.T) {
			_, index, _ := subject()
			assert.Equal(t, NotFoundIndex, index)
		})

		t.Run("should return nil Error", func(t *testing.T) {
			_, _, err := subject()
			assert.Nil(t, err)
		})
	})

	t.Run("when comment data given", func(t *testing.T) {
		data = []byte("  # comment")

		t.Run("should return DirectiveTypeComment", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, DirectiveTypeComment, typ)
		})

		t.Run("should return index of # symbol", func(t *testing.T) {
			_, index, _ := subject()
			assert.Equal(t, strings.Index(string(data), "#"), index)
		})

		t.Run("should return nil Error", func(t *testing.T) {
			_, _, err := subject()
			assert.Nil(t, err)
		})
	})

	t.Run("when meaningful data starts with tab", func(t *testing.T) {
		data = []byte("  \t ")

		t.Run("should return DirectiveTypeUnknown", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, DirectiveTypeUnknown, typ)
		})

		t.Run("should return index of tab", func(t *testing.T) {
			_, index, _ := subject()
			assert.Equal(t, strings.Index(string(data), "\t"), index)
		})

		t.Run("should return TabInIndentationError", func(t *testing.T) {
			_, _, err := subject()
			assert.Equal(t, TabInIndentationError, err)
		})
	})

	t.Run("when text data given", func(t *testing.T) {
		data = []byte("  > text")

		t.Run("should return DirectiveTypeText", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, DirectiveTypeText, typ)
		})

		t.Run("should return index of > symbol", func(t *testing.T) {
			_, index, _ := subject()
			assert.Equal(t, strings.Index(string(data), ">"), index)
		})

		t.Run("should return nil error", func(t *testing.T) {
			_, _, err := subject()
			assert.Nil(t, err)
		})
	})

	t.Run("when list data given", func(t *testing.T) {
		data = []byte("  - list")

		t.Run("should return DirectiveTypeList", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, DirectiveTypeList, typ)
		})

		t.Run("should return index of - symbol", func(t *testing.T) {
			_, index, _ := subject()
			assert.Equal(t, strings.Index(string(data), "-"), index)
		})

		t.Run("should return nil error", func(t *testing.T) {
			_, _, err := subject()
			assert.Nil(t, err)
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

	subject := func() ([]byte, bool, error) {
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

		t.Run("should return NextDirectiveAppearedError", func(t *testing.T) {
			prepare()
			_, hasNext, err := subject()
			assert.True(t, hasNext)
			assert.Nil(t, err)
		})
		t.Run("should return first line of next different directive", func(t *testing.T) {
			prepare()
			nextLine, _, _ := subject()
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
			_, _, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _, _ := subject()
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
			_, _, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _, _ := subject()
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
			_, _, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _, _ := subject()
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
			_, _, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _, _ := subject()
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
			_, _, err := subject()
			assert.Nil(t, err)
		})
		t.Run("should return nil for next line", func(t *testing.T) {
			prepare()
			nextLine, _, _ := subject()
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

			t.Run("should return DifferentTypesOnTheSameLevelError", func(t *testing.T) {
				prepare()
				_, _, err := subject()
				assert.Equal(t, DifferentTypesOnTheSameLevelError, err)
			})

			t.Run("should return nil for nextLine", func(t *testing.T) {
				prepare()
				nextLine, _, _ := subject()
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

					t.Run("should return DifferentLevelOnSameChildError", func(t *testing.T) {
						prepare()
						_, _, err := subject()
						assert.Equal(t, DifferentLevelOnSameChildError, err)
					})

					t.Run("should return nil for nextLine", func(t *testing.T) {
						prepare()
						nextLine, _, _ := subject()
						assert.Nil(t, nextLine)
					})
				})

				t.Run("when following line is not text", func(t *testing.T) {
					condition = func() {
						content = []byte(fmt.Sprintf(contentPh, "- list"))
					}

					t.Run("should return TextHasChildError with DifferentLevelOnSameChildError", func(t *testing.T) {
						prepare()
						_, _, err := subject()
						assert.Equal(t, TextHasChildError, err)
					})

					t.Run("should return nil for nextLine", func(t *testing.T) {
						prepare()
						nextLine, _, _ := subject()
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

					t.Run("should return DifferentLevelOnSameChildError", func(t *testing.T) {
						prepare()
						_, _, err := subject()
						assert.Equal(t, DifferentLevelOnSameChildError, err)
					})

					t.Run("should return nil for nextLine", func(t *testing.T) {
						prepare()
						nextLine, _, _ := subject()
						assert.Nil(t, nextLine)
					})
				})

				t.Run("when following line is not text", func(t *testing.T) {

					condition = func() {
						content = []byte(fmt.Sprintf(contentPh, "- list"))
					}

					t.Run("should return NextDirectiveAppearedError", func(t *testing.T) {
						prepare()
						_, hasNext, err := subject()
						assert.True(t, hasNext)
						assert.Nil(t, err)
					})

					t.Run("should return first line of next directive", func(t *testing.T) {
						prepare()
						nextLine, _, _ := subject()
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

			t.Run("should return DifferentLevelOnSameChildError", func(t *testing.T) {
				prepare()
				_, _, err := subject()
				assert.Equal(t, DifferentLevelOnSameChildError, err)
			})

			t.Run("should return nil for next line", func(t *testing.T) {
				prepare()
				nextLine, _, _ := subject()
				assert.Nil(t, nextLine)
			})
		})
	})
}

func TestRemoveStringTrailingLineBreaks(t *testing.T) {
	t.Run("should remove trailing line break", func(t *testing.T) {
		str := "hello world\n"
		removeStringTrailingLineBreaks(&str)
		assert.Equal(t, "hello world", str)
	})

	t.Run("should remove only single line break", func(t *testing.T) {
		t.Run("consequent line breaks", func(t *testing.T) {
			str := "hello world\n\n"
			removeStringTrailingLineBreaks(&str)
			assert.Equal(t, "hello world\n", str)
		})
		t.Run("multilines", func(t *testing.T) {
			str := "hello\nworld\n"
			removeStringTrailingLineBreaks(&str)
			assert.Equal(t, "hello\nworld", str)
		})
	})

	t.Run("should not remove character if last character is not a line break", func(t *testing.T) {
		str := "hello world"
		removeStringTrailingLineBreaks(&str)
		assert.Equal(t, "hello world", str)
	})
}

func TestRemoveBytesTrailingLineBreaks(t *testing.T) {
	t.Run("should remove trailing line break", func(t *testing.T) {
		b := []byte("hello world\n")
		removeBytesTrailingLineBreaks(&b)
		assert.Equal(t, []byte("hello world"), b)
	})

	t.Run("should remove only single line break", func(t *testing.T) {
		t.Run("consequent line breaks", func(t *testing.T) {
			b := []byte("hello world\n\n")
			removeBytesTrailingLineBreaks(&b)
			assert.Equal(t, []byte("hello world\n"), b)
		})
		t.Run("multilines", func(t *testing.T) {
			b := []byte("hello\nworld\n")
			removeBytesTrailingLineBreaks(&b)
			assert.Equal(t, []byte("hello\nworld"), b)
		})
	})

	t.Run("should not remove character if last character is not a line break", func(t *testing.T) {
		b := []byte("hello world")
		removeBytesTrailingLineBreaks(&b)
		assert.Equal(t, []byte("hello world"), b)
	})
}
