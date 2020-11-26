package ntgo

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	var data []byte

	subject := func() (*Value, error) {
		value := &Value{}
		err := value.Parse(data)
		return value, err
	}

	t.Run("back and forth", func(t *testing.T) {
		data, _ = ioutil.ReadFile("./sample/sample.nt")

		t.Run("should keep same schema", func(t *testing.T) {
			d, err := subject()

			assert.Nil(t, err)
			str := d.ToNestedText()

			another := &Value{}
			err = another.Parse([]byte(str))

			var deepEqual func (*testing.T, *Value, *Value)

			deepEqual = func (t *testing.T, d1 *Value, d2 *Value) {
				switch d1.Type {
				case ValueTypeString:
					assert.Equal(t, d1.String, d2.String)
				case ValueTypeText:
					assert.Equal(t, len(d1.Text), len(d2.Text))
					for i, _ := range d1.Text {
						assert.Equal(t, d1.Text[i], d2.Text[i])
					}
				case ValueTypeList:
					assert.Equal(t, len(d1.List), len(d2.List))
					for i, _ := range d1.List {
						deepEqual(t, d1.List[i], d2.List[i])
					}
				case ValueTypeDictionary:
					assert.Equal(t, len(d1.Dictionary), len(d2.Dictionary))
					for k, _ := range d1.Dictionary {
						deepEqual(t, d1.Dictionary[k], d2.Dictionary[k])
					}
				}
			}

			deepEqual(t, d, another)
		})
	})

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

			t.Run("Type should be ValueTypeText", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeText, value.Type)
			})
			t.Run("Text should be input string array", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, expect, value.Text)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("List should contain values with ValueTypeString", func(t *testing.T) {
				value, err := subject()
				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(value.List))

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]
					assert.Equal(t, ValueTypeString, element.Type)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("List should contain string with leading space", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(value.List))

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]
					assert.Equal(t, ValueTypeString, element.Type)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("List should contain values with ValueTypeText", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(value.List))

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]
					assert.Equal(t, ValueTypeText, element.Type)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("Text values should be in the same depth", func(t *testing.T) {
				value, err := subject()
				assert.Nil(t, err)

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]
					assert.Equal(t, value.Depth+1, element.Depth)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("List should contain values with ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(value.List))

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]
					assert.Equal(t, ValueTypeList, element.Type)
					for j := 0; j < len(element.List); j++ {
						nestedElement := element.List[j]
						assert.Equal(t, ValueTypeString, nestedElement.Type)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("List values should be in the same depth", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(value.List))

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]
					assert.Equal(t, value.Depth+1, element.Depth)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("List should contain values with ValueTypeDictionary", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expect), len(value.List))

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]

					assert.Equal(t, ValueTypeDictionary, element.Type)
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

			t.Run("Type should be ValueTypeList", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
			})
			t.Run("Dictionary elements should be in the same depth", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]

					assert.Equal(t, element.Depth+1, element.Dictionary[expect[i][0][0]].Depth)
				}
			})

			t.Run("Dictionary elements value string should contain leading and trailing spaces and tabs", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)

				for i := 0; i < len(value.List); i++ {
					element := value.List[i]

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

			t.Run("Type should be ValueTypeDictionary", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeDictionary, value.Type)
			})
			t.Run("Dictionary should contain values with ValueTypeString and certain keys", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, len(expectKey), len(value.Dictionary))

				assert.Equal(t, expectValue[0], value.Dictionary[expectKey[0]].String)
				assert.Equal(t, expectValue[1], value.Dictionary[expectKey[1]].String)
			})
		})

		t.Run("empty string", func(t *testing.T) {
			t.Run("space after delimiter", func(t *testing.T) {
				data = []byte("key1: \nkey2: ")
				t.Run("should treat value as empty string", func(t *testing.T) {
					value, err := subject()
					assert.Nil(t, err)
					assert.Equal(t, "", value.Dictionary["key1"].String)
					assert.Equal(t, "", value.Dictionary["key2"].String)
				})
			})
			t.Run("no space after delimiter", func(t *testing.T) {
				data = []byte("key1:\nkey2:")
				t.Run("should treat value as empty string", func(t *testing.T) {
					value, err := subject()
					assert.Nil(t, err)
					assert.Equal(t, "", value.Dictionary["key1"].String)
					assert.Equal(t, "", value.Dictionary["key2"].String)
				})
			})
		})
	})

	t.Run("line breaks", func(t *testing.T) {
		t.Run("cr", func(t *testing.T) {
			data = []byte("- elem1\r- elem2")

			t.Run("should parse regulary", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
				assert.Equal(t, "elem1", value.List[0].String)
				assert.Equal(t, "elem2", value.List[1].String)
			})
		})
		t.Run("lf", func(t *testing.T) {
			data = []byte("- elem1\n- elem2")

			t.Run("should parse regulary", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
				assert.Equal(t, "elem1", value.List[0].String)
				assert.Equal(t, "elem2", value.List[1].String)
			})
		})
		t.Run("crlf", func(t *testing.T) {
			data = []byte("- elem1\r\n- elem2")

			t.Run("should parse regulary", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
				assert.Equal(t, "elem1", value.List[0].String)
				assert.Equal(t, "elem2", value.List[1].String)
			})
		})

		t.Run("mixed", func(t *testing.T) {
			data = []byte("- elem1\r\n- elem2\r- elem3\n- elem4")

			t.Run("should parse regulary", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, ValueTypeList, value.Type)
				assert.Equal(t, "elem1", value.List[0].String)
				assert.Equal(t, "elem2", value.List[1].String)
				assert.Equal(t, "elem3", value.List[2].String)
				assert.Equal(t, "elem4", value.List[3].String)
			})
		})

		t.Run("mixed as text content", func(t *testing.T) {
			data = []byte("text:\n  > line1\r\n  > line2\r  > line3\n  > line4\r\n  > line5")

			t.Run("should parse regulary", func(t *testing.T) {
				value, err := subject()

				assert.Nil(t, err)

				text := value.Dictionary["text"].Text
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

func TestToNestedText(t *testing.T) {

	var data []byte
	var indentSize int
	var depth int

	expect := func() string { return "" }

	resetCondition := func() {
		indentSize = 2
		depth = 0
	}

	subject := func() string {
		value := &Value{
			IndentSize: indentSize,
			Depth:      depth,
		}
		value.Parse(data)
		return value.ToNestedText()
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
		line2 := []byte("- bbbb\n")
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

			expect1 := func() string { return fmt.Sprintf("%s\n%s\n", string(line1), string(line2)) }
			expect2 := func() string { return fmt.Sprintf("%s\n%s\n", string(line2), string(line1)) }

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
				return fmt.Sprintf("%s%s\n%s%s\n", indent, string(line1), indent, string(line2))
			}

			expect2 := func() string {
				indent := fmt.Sprintf("%*s", depth*indentSize, "")
				return fmt.Sprintf("%s%s\n%s%s\n", indent, string(line2), indent, string(line1))
			}

			t.Run("should return text with indentation", func(t *testing.T) {
				// dictionary is unordered
				result := subject()
				assert.True(t, expect1() == result || expect2() == result)
			})
		})
	})
}

func TestDetectValueType(t *testing.T) {
	var data []byte
	subject := func() (ValueType, int, error) {
		return detectValueType(data)
	}

	t.Run("when empty data given", func(t *testing.T) {
		data = []byte("  ")

		t.Run("should return ValueTypeUnknown", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, ValueTypeUnknown, typ)
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

		t.Run("should return ValueTypeComment", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, ValueTypeComment, typ)
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

		t.Run("should return ValueTypeUnknown", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, ValueTypeUnknown, typ)
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

		t.Run("should return ValueTypeText", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, ValueTypeText, typ)
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

		t.Run("should return ValueTypeList", func(t *testing.T) {
			typ, _, _ := subject()
			assert.Equal(t, ValueTypeList, typ)
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

func TestReadTextValue(t *testing.T) {

	var value *Value

	var index int
	var firstLine []byte
	var buffer *bytes.Buffer

	var content []byte

	var bufferInitializer func() *bytes.Buffer

	prepare := func() {
		value = &Value{}
		bufferInitializer = func() *bytes.Buffer {
			return bytes.NewBuffer(content)
		}
	}

	condition := func() {}

	subject := func() ([]byte, bool, error) {
		condition()
		buffer = bufferInitializer()
		return value.readTextValue(index, firstLine, buffer)
	}

	t.Run("when next value appeared", func(t *testing.T) {
		condition = func() {
			index = 2
			firstLine = []byte("  > first line\n")
			content = []byte(
				`  > second line
  > third line
- list`)
		}

		t.Run("should return NextValueAppearedError", func(t *testing.T) {
			prepare()
			_, hasNext, err := subject()
			assert.True(t, hasNext)
			assert.Nil(t, err)
		})
		t.Run("should return first line of next different value", func(t *testing.T) {
			prepare()
			nextLine, _, _ := subject()
			assert.Equal(t, []byte("- list"), nextLine)
		})
		t.Run("should Text slice ends with last line of text value", func(t *testing.T) {
			prepare()
			subject()
			assert.Equal(t, 3, len(value.Text))
			assert.Equal(t, "third line", value.Text[2])
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
			assert.Equal(t, 3, len(value.Text))
			assert.Equal(t, "third line", value.Text[2])
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
			assert.Equal(t, 2, len(value.Text))
			assert.Equal(t, "", value.Text[1])
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
			assert.Equal(t, 2, len(value.Text))
			assert.Equal(t, "first line\n", value.Text[0])
			assert.Equal(t, "second line", value.Text[1])
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
			assert.Equal(t, 2, len(value.Text))
			assert.Equal(t, "first line\n", value.Text[0])
			assert.Equal(t, "second line", value.Text[1])
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
			assert.Equal(t, 2, len(value.Text))
			assert.Equal(t, "first line\n", value.Text[0])
			assert.Equal(t, "second line", value.Text[1])
		})
	})

	t.Run("irregulars", func(t *testing.T) {
		t.Run("when Type is already defined", func(t *testing.T) {
			condition = func() {
				index = 0
				firstLine = []byte("> first line\n")
				content = []byte("> second line")

				value.Type = ValueTypeText
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

					t.Run("should return NextValueAppearedError", func(t *testing.T) {
						prepare()
						_, hasNext, err := subject()
						assert.True(t, hasNext)
						assert.Nil(t, err)
					})

					t.Run("should return first line of next value", func(t *testing.T) {
						prepare()
						nextLine, _, _ := subject()
						assert.Equal(t, []byte("- list"), nextLine)
					})
				})
			})
		})

		t.Run("when content consists of lines with different values", func(t *testing.T) {
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

func TestValueType(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		t.Run("shold retrun alias", func(t *testing.T) {
			assert.Equal(t, "string", ValueTypeString.String())
			assert.Equal(t, "text", ValueTypeText.String())
			assert.Equal(t, "list", ValueTypeList.String())
			assert.Equal(t, "dictionary", ValueTypeDictionary.String())
		})
	})
}

func TestMultiLineText(t *testing.T) {
	m := MultiLineText{
		"line 1\n",
		"line 2",
	}

	t.Run("String", func(t *testing.T) {
		t.Run("shold joins elements with empty glue", func(t *testing.T) {
			assert.Equal(t, strings.Join(m, ""), m.String())
		})
	})
}