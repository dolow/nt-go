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
		t.Run("string elements string", func(t *testing.T) {
			expect := []string{"string", "elements"}
			data = []byte(fmt.Sprintf("- %s\n- %s", expect[0], expect[1]))

			t.Run("Type should be DirectiveTypeList", func(t *testing.T) {
				directive, err := subject()

				assert.Nil(t, err)
				assert.Equal(t, DirectiveTypeList, directive.Type)
			})
			t.Run("List should directives with DirectiveTypeString", func(t *testing.T) {
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
	})
}
