package nestedtext

import (
		"testing"

		"github.com/stretchr/testify/assert"
)


func TestParse(t *testing.T) {
	
	t.Run("string", func(t *testing.T) {
		data := []byte("plain text")

		subject := func() (*Directive, error) {
			directive := &Directive{}
			return directive.Parse(data)
		}

		t.Run("Type should be DirectiveTypeString", func(t *testing.T) {
			directive, err := subject();

			assert.Nil(t, err)
			assert.Equal(t, directive.Type, DirectiveTypeString)
		})

		t.Run("String should be input string", func(t *testing.T) {
			directive, err := subject();
			
			assert.Nil(t, err)
			assert.Equal(t, directive.String, string(data))
		})
	})
}