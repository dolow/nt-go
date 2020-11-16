package ntgo

import (
	"testing"
)

func Benchmark_Directive(b *testing.B) {
	b.Run("Marshal", func(b *testing.B) {
		b.Run("String", func(b *testing.B) {
			content := []byte(`key: string`)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Marshal(content)
			}
		})
		b.Run("Text", func(b *testing.B) {
			content := []byte(`key:
  > multiline
  > text`)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Marshal(content)
			}
		})

		b.Run("List", func(b *testing.B) {
			content := []byte(`key:
  - list
  - element`)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Marshal(content)
			}
		})

		b.Run("Dictionary", func(b *testing.B) {
			content := []byte(`key:
  child1: val1
  child2: val2`)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive := &Directive{}
				directive.Marshal(content)
			}
		})
	})


	b.Run("Unmarshal", func(b *testing.B) {
		b.Run("String", func(b *testing.B) {
			content := []byte(`key: string`)
			directive := &Directive{}
			directive.Marshal(content)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.Unmarshal()
			}
		})
		b.Run("Text", func(b *testing.B) {
			content := []byte(`key:
  > multiline
  > text`)
  			directive := &Directive{}
			directive.Marshal(content)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.Unmarshal()
			}
		})

		b.Run("List", func(b *testing.B) {
			content := []byte(`key:
  - list
  - element`)
			directive := &Directive{}
			directive.Marshal(content)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.Unmarshal()
			}
		})

		b.Run("Dictionary", func(b *testing.B) {
			content := []byte(`key:
  child1: val1
  child2: val2`)
			directive := &Directive{}
			directive.Marshal(content)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				directive.Unmarshal()
			}
		})
	})
}
