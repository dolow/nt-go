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

func TestJSON(t *testing.T) {
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
		directive := &Directive{}

		var err *DirectiveMarshalError

		err = directive.Marshal([]byte(""))
		assert.NotNil(t, err)
		assert.Equal(t, EmptyDataError, err.error)

		err = directive.Marshal([]byte("\n  \n"))
		assert.NotNil(t, err)
		assert.Equal(t, EmptyDataError, err.error)
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



}

