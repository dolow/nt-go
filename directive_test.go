package nestedtext

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	var data []byte

	subject := func() (*Directive, error) {
		directive := &Directive{}
		return directive.Parse(data)
	}

	t.Run("string", func(t *testing.T) {

		expect := "plain text"

		t.Run("regular string", func(t *testing.T) {
			data = []byte(expect)

			t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeString, directive.Type)
			})
			t.Run("String should be input string", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.String)
			})
		})

		t.Run("string start with space", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  %s", expect))

			t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeString, directive.Type)
			})
			t.Run("String should be input string without initial spaces", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.String)
			})
		})

		t.Run("string start with line break", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n%s", expect))

			t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeString, directive.Type)
			})
			t.Run("Empty line should be ignored", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.String)
			})
		})

		t.Run("string start with line break with forward spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("  \n%s", expect))

			t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeString, directive.Type)
			})
			t.Run("Empty line should be ignored", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.String)
			})
		})

		t.Run("string start with line break and second line starts with spaces", func(t *testing.T) {
			data = []byte(fmt.Sprintf("\n  %s", expect))

			t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeString, directive.Type)
			})
			t.Run("Empty line and forwarding spaces should be ignored", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, directive.String)
			})
		})

		t.Run("string with comment symbol (#)", func(t *testing.T) {
			expect = "plain text # it is not comment"
			data = []byte(expect)

			t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeString, directive.Type)
			})
			t.Run("After comment symbol is also treated as string content", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, directive.String, expect)
			})
		})

		t.Run("string ending with line break", func(t *testing.T) {
			expect = "plain text"
			data = []byte(fmt.Sprintf("%s\n", expect))

			t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeString, directive.Type)
			})
			t.Run("After comment symbol is also treated as string content", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, directive.String, expect)
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
			expect := []string{"string", "elements"}
			data = []byte(fmt.Sprintf(
`-    %s
- %s`,
				expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List should contain trimmed string", func(t *testing.T) {
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
				[]string { "aaaa", "bbbb" },
				[]string { "cccc", "dddd" },
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
						if j != len(element.Text) - 1 {
							e = fmt.Sprintf("%s\n", e)
						}
						assert.Equal(t, e, element.Text[j])
					}
				}
			})
		})

		t.Run("text elements with unbalanced indentations", func(t *testing.T) {
			expect := [][]string{
				[]string { "aaaa", "bbbb" },
				[]string { "cccc", "dddd" },
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
					assert.Equal(t, directive.Depth + 1, element.Depth)
				}
			})
		})

		t.Run("list elements", func(t *testing.T) {
			expect := [][]string{
				[]string { "aaaa", "bbbb" },
				[]string { "cccc", "dddd" },
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
				[]string { "aaaa", "bbbb" },
				[]string { "cccc", "dddd" },
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
					assert.Equal(t, directive.Depth + 1, element.Depth)
					for j := 0; j < len(element.List); j++ {
						assert.Equal(t, element.Depth + 1, element.List[j].Depth)
					}
				}
			})
		})

		t.Run("map string elements", func(t *testing.T) {
			expect := [][][]string{
				[][]string {
					{ "key1", "val1" },
					{ "key2", "val2" },
				},
				[][]string {
					{ "key3", "val3" },
					{ "key4", "val4" },
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
			t.Run("List should contain directives with DirectiveTypeMap", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(directive.List))

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]

					assert.Equal(t, DirectiveTypeMap, element.Type)
					assert.Equal(t, expect[i][0][1], element.Map[expect[i][0][0]].String)
					assert.Equal(t, expect[i][1][1], element.Map[expect[i][1][0]].String)
				}
			})
		})

		t.Run("map string elements with unbalanced spaces", func(t *testing.T) {
			expect := [][][]string{
				[][]string {
					{ "key1", "val1" },
					{ "key2", "val2" },
				},
				[][]string {
					{ "key3", "val3" },
					{ "key4", "val4" },
				},
			}
			data = []byte(fmt.Sprintf(
`-
  %s:    %s
  %s: %s
-
     %s: %s
     %s:   %s`,
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
			t.Run("Map elements should be in the same depth", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]

					assert.Equal(t, element.Depth + 1, element.Map[expect[i][0][0]].Depth)
				}
			})

			t.Run("Map elements value string should be trimmed", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)

				for i := 0; i < len(directive.List); i++ {
					element := directive.List[i]

					assert.Equal(t, expect[i][0][1], element.Map[expect[i][0][0]].String)
					assert.Equal(t, expect[i][1][1], element.Map[expect[i][1][0]].String)
				}
			})
		})
	})


	t.Run("map", func(t *testing.T) {
		t.Run("string elements", func(t *testing.T) {
			expectKey   := []string{"key1",   "key2"}
			expectValue := []string{"value1", "value2"}

			data = []byte(fmt.Sprintf(
`%s: %s
%s: %s`,
				expectKey[0], expectValue[0],
				expectKey[1], expectValue[1],
			))

			t.Run("Type should be DirectiveTypeMap", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeMap, directive.Type)
			})
			t.Run("Map should contain directives with DirectiveTypeString and certain keys", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expectKey), len(directive.Map))

				assert.Equal(t, expectValue[0], directive.Map[expectKey[0]].String)
				assert.Equal(t, expectValue[1], directive.Map[expectKey[1]].String)
			})
		})
	})
}


func TestToString(t *testing.T) {

	var data []byte
	var indentSize int
	var depth      int

	expect := func() string { return "" }

	resetCondition := func() {
		indentSize = 2
		depth = 0
	}

	subject := func() string {
		directive := &Directive{
			IndentSize: indentSize,
			Depth: depth,
		}
		directive.Parse(data)
		return directive.ToString()
	}

	t.Run("string", func(t *testing.T) {
		data = []byte("stringvalue")

		t.Run("Depth is 0", func(t *testing.T) {
			defer resetCondition()

			depth = 0

			expect = func() string { return string(data) }

			t.Run("should return same value as input", func(t *testing.T) {
				assert.Equal(t, expect(), subject())
			})
		})
		t.Run("Depth is larger than 0", func(t *testing.T) {
			defer resetCondition()

			indentSize = 2
			depth = 4

			expect = func() string { return string(data) }

			t.Run("should not indent even if Depth and IndentSize is larger than 0", func(t *testing.T) {
				assert.Equal(t, expect(), subject())
			})
		})
	})

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
				indent := fmt.Sprintf("%*s", depth * indentSize, "")
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
				indent := fmt.Sprintf("%*s", depth * indentSize, "")
				return fmt.Sprintf("%s%s%s%s", indent, string(line1), indent, string(line2))
			}

			t.Run("should return text with indentation", func(t *testing.T) {
				assert.Equal(t, expect(), subject())
			})	
		})
	})

	t.Run("map", func(t *testing.T) {
		line1 := []byte("key1: value1")
		line2 := []byte("key2: value2")
		data = []byte(fmt.Sprintf("%s\n%s", string(line1), string(line2)))

		t.Run("Depth is 0", func(t *testing.T) {
			defer resetCondition()

			depth = 0

			expect1 := func() string { return fmt.Sprintf("%s\n%s", string(line1), string(line2)) }
			expect2 := func() string { return fmt.Sprintf("%s\n%s", string(line2), string(line1)) }

			t.Run("should return text with no indentation", func(t *testing.T) {
				// map is unordered
				assert.True(t, expect1() == subject() || expect2() == subject())
			})
		})

		t.Run("Depth and IndentSize is larger than 0", func(t *testing.T) {
			defer resetCondition()

			depth = 2
			indentSize = 4

			expect1 := func() string {
				indent := fmt.Sprintf("%*s", depth * indentSize, "")
				return fmt.Sprintf("%s%s\n%s%s", indent, string(line1), indent, string(line2))
			}

			expect2 := func() string {
				indent := fmt.Sprintf("%*s", depth * indentSize, "")
				return fmt.Sprintf("%s%s\n%s%s", indent, string(line2), indent, string(line1))
			}

			t.Run("should return text with indentation", func(t *testing.T) {
				// map is unordered
				assert.True(t, expect1() == subject() || expect2() == subject())
			})	
		})
	})
}