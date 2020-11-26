package ntgo

import (
	"fmt"
	"reflect"
	"strings"
)

func Marshal(content string, v interface{}) {
	value := &Value{}
	value.Parse([]byte(content))

	typ := reflect.TypeOf(v)

	var ref reflect.Value
	if typ.Kind() == reflect.Ptr {
		ref = reflect.ValueOf(v)
		typ = typ.Elem()
	} else {
		ref = reflect.New(typ)
	}

	marshal(value, typ, &ref)
}

func Unmarshal(v interface{}) string {
	value := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	result, _ := unmarshal(typ, &value, 0, 0)

	return result
}

func marshalSlice(value *Value, typ reflect.Type, ref *reflect.Value) {
	// type of slice element
	switch typ.Kind() {
	case reflect.String:
		{
			// multiline text
			switch value.Type {
			case ValueTypeString:
				*ref = reflect.Append(*ref, reflect.ValueOf(value.String))
			case ValueTypeText:
				{
					for _, line := range value.Text {
						*ref = reflect.Append(*ref, reflect.ValueOf(line))
					}
				}
			case ValueTypeList:
				{
					for _, child := range value.List {
						*ref = reflect.Append(*ref, reflect.ValueOf(child.String))
					}
				}
			}
		}
	case reflect.Slice:
		{
			for _, child := range value.List {
				childWork := reflect.MakeSlice(typ, 0, cap(child.List))
				marshalSlice(child, typ.Elem(), &childWork)
				*ref = reflect.Append(*ref, childWork)
			}
		}
	case reflect.Struct:
		{
			for _, child := range value.List {
				elementInstance := reflect.New(typ).Elem()
				marshal(child, typ, &elementInstance)
				*ref = reflect.Append(*ref, elementInstance)
			}
		}
	case reflect.Ptr:
		{
			for _, child := range value.List {
				elementType := typ.Elem()
				elementInstance := reflect.New(elementType)
				marshal(child, elementType, &elementInstance)
				*ref = reflect.Append(*ref, elementInstance)
			}
		}
	}
}

func marshal(value *Value, typ reflect.Type, ref *reflect.Value) {
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

		childValue, exists := value.Dictionary[key]
		if !exists {
			continue
		}

		fieldType := fieldInfo.Type

		switch fieldType.Kind() {
		case reflect.String:
			{
				if childValue.Type == ValueTypeText {
					fieldRef.SetString(strings.Join(childValue.Text, ""))
				} else {
					fieldRef.SetString(childValue.String)
				}
			}
		case reflect.Slice:
			{
				work := reflect.MakeSlice(fieldRef.Type(), 0, cap(childValue.List))
				marshalSlice(childValue, fieldType.Elem(), &work)
				fieldRef.Set(work)
			}
		case reflect.Struct:
			{
				fieldInstance := reflect.New(fieldType).Elem()
				marshal(childValue, fieldType, &fieldInstance)
				fieldRef.Set(fieldInstance)
			}
		case reflect.Ptr:
			{
				fieldType := fieldType.Elem()
				fieldInstance := reflect.New(fieldType)
				switch fieldType.Kind() {
				case reflect.Struct:
					marshal(childValue, fieldType, &fieldInstance)
					fieldRef.Set(fieldInstance)
				case reflect.String:
					fieldRef.Set(reflect.ValueOf(&childValue.String))
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
			if i == len(lines)-1 {
				result += fmt.Sprintf("%s%s %s", fmt.Sprintf("%*s", depth*UnmarshalDefaultIndentSize, ""), string(TextSymbol), line)
			} else {
				result += fmt.Sprintf("%s%s %s%s", fmt.Sprintf("%*s", depth*UnmarshalDefaultIndentSize, ""), string(TextSymbol), line, string(LF))
			}
		}
		return result, true
	case reflect.Slice:
		{
			var result string
			var lineBreakAfterKey string
			valueSymbol := ListSymbol

			sliceType := typ.Elem()
			switch sliceType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				lineBreakAfterKey = string(Space)
			case reflect.String:
				lineBreakAfterKey = string(Space)
				if (tagFlag & MarshallerTagFlagMultilineText) == MarshallerTagFlagMultilineText {
					valueSymbol = TextSymbol
				}
			default:
				lineBreakAfterKey = string(LF)
			}

			for i := 0; i < ref.Len(); i++ {
				childRef := ref.Index(i)

				childContent, _ := unmarshal(sliceType, &childRef, depth+1, tagFlag)
				result += fmt.Sprintf("%s%s%s%s", fmt.Sprintf("%*s", depth*UnmarshalDefaultIndentSize, ""), string(valueSymbol), lineBreakAfterKey, childContent)
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

				childTagFlag := getTagFlagFromTagValue(tagValues)

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

				marshalizedValue, exists := unmarshal(fieldType, &fieldRef, depth+1, childTagFlag)

				if !exists && ((childTagFlag & MarshallerTagFlagOmitEmpty) == MarshallerTagFlagOmitEmpty) {
					continue
				}
				result += fmt.Sprintf("%s%s:%s%s", fmt.Sprintf("%*s", depth*UnmarshalDefaultIndentSize, ""), key, lineBreakAfterKey, marshalizedValue)
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

func getTagFlagFromTagValue(tagValues []string) (flag int) {
	for i := 1; i < len(tagValues); i++ {
		switch tagValues[i] {
		case MarshallerTagMultilineText:
			flag |= MarshallerTagFlagMultilineText
		case MarshallerTagOmitEmpty:
			flag |= MarshallerTagFlagOmitEmpty
		}
	}

	return
}
