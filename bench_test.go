package ntgo

import (
	"testing"
)

const StringSample = `key: string`
const TextSample = `key:
  > multiline
  > text`
const ListSample = `key:
  - list
  - element`
const DictionarySample = `key:
  child1: val1
  child2: val2`

func Benchmark_Directive(b *testing.B) {
	b.Run("Parse", func(b *testing.B) {
		b.Run("String", func(b *testing.B) {
			content := []byte(StringSample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Parse(content)
			}
		})
		b.Run("Text", func(b *testing.B) {
			content := []byte(TextSample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Parse(content)
			}
		})

		b.Run("List", func(b *testing.B) {
			content := []byte(ListSample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Parse(content)
			}
		})

		b.Run("Dictionary", func(b *testing.B) {
			content := []byte(DictionarySample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Parse(content)
			}
		})

		b.Run("Large data", func(b *testing.B) {
			str := `root:
`

			for i := 0; i < 100; i++ {
				str += `  -
    -
      -
        - a
        - b
      -
        - c
        - d
    -
      -
        - e
        - f
      -
        - g
        - h
`
			}

			content := []byte(str)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Parse(content)
			}
		})
	})

	b.Run("ToString", func(b *testing.B) {
		prepare := func(str string) *Directive {
			content := []byte(str)
			directive := &Directive{}
			directive.Parse(content)
			return directive
		}

		b.Run("String", func(b *testing.B) {
			directive := prepare(StringSample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.ToString()
			}
		})
		b.Run("Text", func(b *testing.B) {
			directive := prepare(TextSample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.ToString()
			}
		})

		b.Run("List", func(b *testing.B) {
			directive := prepare(ListSample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.ToString()
			}
		})

		b.Run("Dictionary", func(b *testing.B) {
			directive := prepare(DictionarySample)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.ToString()
			}
		})

		b.Run("Mixed large data", func(b *testing.B) {
			content := `root:
`

			for i := 0; i < 100; i++ {
				content += `  -
    -
      -
        - a
        - b
      -
        - c
        - d
    -
      -
        - e
        - f
      -
        - g
        - h
`
			}
			directive := prepare(string(content))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.ToString()
			}
		})
	})
}
