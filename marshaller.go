package ntgo

import (
	"fmt"
	"reflect"
	"strings"
)

func Marshal(content string, v interface{}) {
	directive := &Directive{}
	directive.Parse([]byte(content))

	value := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	var ref reflect.Value
	if typ.Kind() == reflect.Ptr {
		ref = value
		typ = typ.Elem()
	} else {
		ref = reflect.New(typ)
	}

	marshal(directive, typ, &ref)
}

func Unmarshal(v interface{}) string {
	value := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	return unmarshal(typ, &value, 0, nil)
}

func marshalSlice(directive *Directive, typ reflect.Type, ref *reflect.Value) {
	// type of slice element
	switch typ.Kind() {
	case reflect.String:
		{
			// multiline text
			work := *ref
			if directive.Type == DirectiveTypeText {
				for _, line := range directive.Text {
					work = reflect.Append(work, reflect.ValueOf(line))
				}
			} else if directive.Type == DirectiveTypeList {
				for _, child := range directive.List {
					work = reflect.Append(work, reflect.ValueOf(child.String))
				}
			}
			ref.Set(work)
		}
	case reflect.Slice:
		{
			work := *ref
			for _, child := range directive.List {
				elementInstance := reflect.New(typ).Elem()
				marshalSlice(child, typ.Elem(), &elementInstance)
				work = reflect.Append(work, elementInstance)
			}
			ref.Set(work)
		}
	case reflect.Struct:
		{
			work := *ref
			for _, child := range directive.List {
				elementInstance := reflect.New(typ).Elem()
				marshal(child, typ, &elementInstance)
				work = reflect.Append(work, elementInstance)
			}
			ref.Set(work)

		}
	case reflect.Ptr:
		{
			work := *ref
			for _, child := range directive.List {
				elementType := typ.Elem()
				elementInstance := reflect.New(elementType)
				marshal(child, elementType, &elementInstance)
				work = reflect.Append(work, elementInstance)
			}
			ref.Set(work)
		}
	}
}

func marshal(directive *Directive, typ reflect.Type, ref *reflect.Value) {
	substance := *ref
	if ref.Type().Kind() == reflect.Ptr {
		substance = substance.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		fieldInfo := typ.Field(i)
		fieldRef := substance.Field(i)
		tagValue := fieldInfo.Tag.Get(MarshallerTag)

		if tagValue == "" {
			continue
		}

		tagValues := strings.Split(tagValue, MarshallerTagSeparator)
		key := tagValues[0]

		childDirective, exists := directive.Dictionary[key]
		if !exists {
			continue
		}

		fieldType := fieldInfo.Type

		switch fieldType.Kind() {
		case reflect.String:
			{
				fieldRef.SetString(childDirective.String)
			}
		case reflect.Slice:
			{
				// type of slice element
				marshalSlice(childDirective, fieldType.Elem(), &fieldRef)
			}
		case reflect.Struct:
			{
				fieldInstance := reflect.New(fieldType).Elem()
				marshal(childDirective, fieldType, &fieldInstance)
				fieldRef.Set(fieldInstance)
			}
		case reflect.Ptr:
			{
				fieldType := fieldType.Elem()
				fieldInstance := reflect.New(fieldType)
				switch fieldType.Kind() {
				case reflect.Struct:
					marshal(childDirective, fieldType, &fieldInstance)
					fieldRef.Set(fieldInstance)
				case reflect.String:
					fieldRef.Set(reflect.ValueOf(&childDirective.String))
				}
			}
		}
	}
}

func unmarshal(typ reflect.Type, ref *reflect.Value, depth int, tagOption []string) string {

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d%s", ref.Int(), string(LF))
	case reflect.String:
		var value string
		if ref.Kind() == reflect.Ptr {
			value = ref.Elem().String()
		} else {
			value = ref.String()
		}
		lines := strings.Split(value, string(LF))
		if len(lines) == 1 {
			return fmt.Sprintf("%s%s", value, string(LF))
		}

		result := ""
		for i, line := range lines {
			if i == len(lines) - 1 {
				result += fmt.Sprintf("%s%s %s", createIndent(depth), string(TextSymbol), line)
			} else {
				result += fmt.Sprintf("%s%s %s%s", createIndent(depth), string(TextSymbol), line, string(LF))
			}
		}
		return result
	case reflect.Slice:
		{
			var result string
			var lineBreakAfterKey string
			directiveSymbol := ListSymbol

			sliceType := typ.Elem()
			switch sliceType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				lineBreakAfterKey = string(Space)
			case reflect.String:
				lineBreakAfterKey = string(Space)
				if tagOption != nil {
					for i := 1; i < len(tagOption); i++ {
						if tagOption[i] == MarshallerTagMultilineText {
							directiveSymbol = TextSymbol
							break
						}
					}
				}
			default:
				lineBreakAfterKey = string(LF)
			}

			for i := 0; i < ref.Len(); i++ {
				childRef := ref.Index(i)

				childContent := unmarshal(sliceType, &childRef, depth + 1, tagOption)
				result += fmt.Sprintf("%s%s%s%s", createIndent(depth), string(directiveSymbol), lineBreakAfterKey, childContent)
			}
			return result
		}
	case reflect.Struct:
		{
			substance := *ref
			if ref.Type().Kind() == reflect.Ptr {
				substance = ref.Elem()
			}
			var result string
			for i := 0; i < typ.NumField(); i++ {
				fieldInfo := typ.Field(i)
				fieldRef := substance.Field(i)
				tagValue := fieldInfo.Tag.Get(MarshallerTag)
				if tagValue == "" {
					continue
				}

				fieldType := fieldInfo.Type
				tagValues := strings.Split(tagValue, MarshallerTagSeparator)
				key := tagValues[0]

				var lineBreakAfterKey string

				switch fieldType.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					lineBreakAfterKey = string(Space)
				case reflect.String:
					lineBreakAfterKey = string(Space)

					lines := strings.Split(fieldRef.String(), string(LF))
					if len(lines) > 1 {
						lineBreakAfterKey = string(LF)
					} else {
						for i := 1; i < len(tagValues); i++ {
							if tagValues[i] == MarshallerTagMultilineText {
								lineBreakAfterKey = string(LF)
								break
							}
						}
					}
				case reflect.Ptr:
					if fieldRef.Type().Elem().Kind() == reflect.String {
						lineBreakAfterKey = string(Space)

						lines := strings.Split(fieldRef.String(), string(LF))
						if len(lines) > 1 {
							lineBreakAfterKey = string(LF)
						} else {
							for i := 1; i < len(tagValues); i++ {
								if tagValues[i] == MarshallerTagMultilineText {
									lineBreakAfterKey = string(LF)
									break
								}
							}
						}
					} else {
						lineBreakAfterKey = string(LF)
					}
				default:
					lineBreakAfterKey = string(LF)
				}

				marshalizedValue := unmarshal(fieldType, &fieldRef, depth + 1, tagValues)
				result += fmt.Sprintf("%s%s:%s%s", createIndent(depth), key, lineBreakAfterKey, marshalizedValue)
			}
			return result
		}
	case reflect.Ptr:
		{
			return unmarshal(typ.Elem(), ref, depth, nil)
		}
	}
	return ""
}

func createIndent(depth int) string {
	indent := make([]byte, depth * UnmarshalDefaultIndentSize)

	for i := 0; i < depth; i++ {
		for j := 0; j < UnmarshalDefaultIndentSize; j++ {
			indent[i*UnmarshalDefaultIndentSize+j] = Space
		}
	}

	return string(indent)
}