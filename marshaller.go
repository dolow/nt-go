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

	result, _ := unmarshal(typ, &value, 0, 0)

	return result
}

func marshalSlice(directive *Directive, typ reflect.Type, ref *reflect.Value) {
	// type of slice element
	switch typ.Kind() {
	case reflect.String:
		{
			// multiline text
			switch directive.Type {
			case DirectiveTypeString:
				*ref = reflect.Append(*ref, reflect.ValueOf(directive.String))
			case DirectiveTypeText:
				{
					for _, line := range directive.Text {
						*ref = reflect.Append(*ref, reflect.ValueOf(line))
					}
				}
			case DirectiveTypeList:
				{
					for _, child := range directive.List {
						*ref = reflect.Append(*ref, reflect.ValueOf(child.String))
					}
				}
			}
		}
	case reflect.Slice:
		{
			for _, child := range directive.List {
				childWork := reflect.MakeSlice(typ, 0, cap(child.List))
				marshalSlice(child, typ.Elem(), &childWork)
				*ref = reflect.Append(*ref, childWork)
			}
		}
	case reflect.Struct:
		{
			for _, child := range directive.List {
				elementInstance := reflect.New(typ).Elem()
				marshal(child, typ, &elementInstance)
				*ref = reflect.Append(*ref, elementInstance)
			}
		}
	case reflect.Ptr:
		{
			for _, child := range directive.List {
				elementType := typ.Elem()
				elementInstance := reflect.New(elementType)
				marshal(child, elementType, &elementInstance)
				*ref = reflect.Append(*ref, elementInstance)
			}
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
				if childDirective.Type == DirectiveTypeText {
					fieldRef.SetString(strings.Join(childDirective.Text, ""))
				} else {
					fieldRef.SetString(childDirective.String)
				}
			}
		case reflect.Slice:
			{
				work := reflect.MakeSlice(fieldRef.Type(), 0, cap(childDirective.List))
				marshalSlice(childDirective, fieldType.Elem(), &work)
				fieldRef.Set(work)
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

func unmarshal(typ reflect.Type, ref *reflect.Value, depth int, tagFlag int) (string, bool) {

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d%s", ref.Int(), string(LF)), true
	case reflect.String:
		var value string
		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				return "", false
			}
			value = ref.Elem().String()
		} else {
			value = ref.String()
		}
		lines := strings.Split(value, string(LF))
		if len(lines) == 1 {
			return fmt.Sprintf("%s%s", value, string(LF)), value != ""
		}

		result := ""
		for i, line := range lines {
			if i == len(lines) - 1 {
				result += fmt.Sprintf("%s%s %s", fmt.Sprintf("%*s", depth * UnmarshalDefaultIndentSize, ""), string(TextSymbol), line)
			} else {
				result += fmt.Sprintf("%s%s %s%s", fmt.Sprintf("%*s", depth * UnmarshalDefaultIndentSize, ""), string(TextSymbol), line, string(LF))
			}
		}
		return result, true
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
				if (tagFlag & MarshallerTagFlagMultilineText) == MarshallerTagFlagMultilineText {
					directiveSymbol = TextSymbol
				}
			default:
				lineBreakAfterKey = string(LF)
			}

			for i := 0; i < ref.Len(); i++ {
				childRef := ref.Index(i)

				childContent, _ := unmarshal(sliceType, &childRef, depth + 1, tagFlag)
				result += fmt.Sprintf("%s%s%s%s", fmt.Sprintf("%*s", depth * UnmarshalDefaultIndentSize, ""), string(directiveSymbol), lineBreakAfterKey, childContent)
			}
			return result, ref.Len() > 0
		}
	case reflect.Struct:
		{
			substance := *ref
			if ref.Type().Kind() == reflect.Ptr {
				if substance.IsNil() {
					return "", false
				}
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

				childTagFlag := 0
				for i := 1; i < len(tagValues); i++ {
					switch tagValues[i] {
					case MarshallerTagMultilineText:
						childTagFlag |= MarshallerTagFlagMultilineText
					case MarshallerTagOmitEmpty:
						childTagFlag |= MarshallerTagFlagOmitEmpty
					}
				}

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
						if (childTagFlag & MarshallerTagFlagMultilineText) == MarshallerTagFlagMultilineText {
							lineBreakAfterKey = string(LF)
						}
					}
				case reflect.Ptr:
					if fieldRef.Type().Elem().Kind() == reflect.String {
						lineBreakAfterKey = string(Space)

						lines := strings.Split(fieldRef.String(), string(LF))
						if len(lines) > 1 {
							lineBreakAfterKey = string(LF)
						} else {
							if (childTagFlag & MarshallerTagFlagMultilineText) == MarshallerTagFlagMultilineText {
								lineBreakAfterKey = string(LF)
							}
						}
					} else {
						lineBreakAfterKey = string(LF)
					}
				default:
					lineBreakAfterKey = string(LF)
				}

				marshalizedValue, exists := unmarshal(fieldType, &fieldRef, depth + 1, childTagFlag)

				if !exists && ((childTagFlag & MarshallerTagFlagOmitEmpty) == MarshallerTagFlagOmitEmpty) {
					continue
				}
				result += fmt.Sprintf("%s%s:%s%s", fmt.Sprintf("%*s", depth * UnmarshalDefaultIndentSize, ""), key, lineBreakAfterKey, marshalizedValue)
			}
			return result, true
		}
	case reflect.Ptr:
		{
			if ref.IsNil() {
				if (tagFlag & MarshallerTagFlagOmitEmpty) == MarshallerTagFlagOmitEmpty {
					return "", false
				}
			}
			return unmarshal(typ.Elem(), ref, depth, 0)
		}
	}
	return "", false
}
