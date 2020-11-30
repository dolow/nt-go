package ntgo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ValueIsNotPointerError = errors.New("ntgo: marshaling target must be pointer")
)

func Marshal(content string, v interface{}) error {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr {
		return ValueIsNotPointerError
	}

	value := &Value{}
	value.Parse([]byte(content))

	ref := reflect.ValueOf(v)
	typ = typ.Elem()

	marshal(value, typ, &ref)

	return nil
}

func Unmarshal(v interface{}) string {
	value := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	result, _ := unmarshal(typ, &value, 0, 0)

	return result
}

func marshalSlice(value *Value, elementType reflect.Type, sliceRef *reflect.Value) {
	// type of slice element
	switch elementType.Kind() {
	case reflect.String:
		{
			// multiline string
			switch value.Type {
			case ValueTypeString:
				*sliceRef = reflect.Append(*sliceRef, reflect.ValueOf(value.String))
			case ValueTypeText:
				for _, line := range value.Text {
					*sliceRef = reflect.Append(*sliceRef, reflect.ValueOf(line))
				}
			case ValueTypeList:
				for _, child := range value.List {
					*sliceRef = reflect.Append(*sliceRef, reflect.ValueOf(child.String))
				}
			}
		}
	case reflect.Slice:
		{
			for _, child := range value.List {
				childWork := reflect.MakeSlice(elementType, 0, cap(child.List))
				marshalSlice(child, elementType.Elem(), &childWork)
				*sliceRef = reflect.Append(*sliceRef, childWork)
			}
		}
	case reflect.Struct:
		{
			for _, child := range value.List {
				elementInstance := reflect.New(elementType).Elem()
				marshal(child, elementType, &elementInstance)
				*sliceRef = reflect.Append(*sliceRef, elementInstance)
			}
		}
	case reflect.Ptr:
		{
			switch value.Type {
			case ValueTypeText:
				for i, _ := range value.Text {
					// using value occurs late binding
					*sliceRef = reflect.Append(*sliceRef, reflect.ValueOf(&value.Text[i]))
				}
			case ValueTypeList:
				elementType := elementType.Elem()
				switch elementType.Kind() {
				case reflect.String:
					for _, child := range value.List {
						*sliceRef = reflect.Append(*sliceRef, reflect.ValueOf(&child.String))
					}
				case reflect.Struct:
					for _, child := range value.List {
						elementInstance := reflect.New(elementType)
						marshal(child, elementType, &elementInstance)
						*sliceRef = reflect.Append(*sliceRef, elementInstance)
					}
				}
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
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f%s", ref.Float(), string(LF)), true
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

			sliceElementType := typ.Elem()
			sliceElementPointingType := sliceElementType
			if sliceElementType.Kind() == reflect.Ptr {
				sliceElementPointingType = sliceElementType.Elem()
			}

			switch sliceElementPointingType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				lineBreakAfterKey = string(Space)
			case reflect.String:
				lineBreakAfterKey = string(Space)
				if (tagFlag & MarshallerTagFlagMultilineStrings) == MarshallerTagFlagMultilineStrings {
					valueSymbol = TextSymbol
				}
			default:
				lineBreakAfterKey = string(LF)
			}

			for i := 0; i < ref.Len(); i++ {
				childRef := ref.Index(i)

				childContent, _ := unmarshal(sliceElementType, &childRef, depth+1, tagFlag)
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
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
					lineBreakAfterKey = string(Space)
				case reflect.String:
					lineBreakAfterKey = string(Space)

					lines := strings.Split(fieldRef.String(), string(LF))
					if len(lines) > 1 {
						lineBreakAfterKey = string(LF)
					} else {
						if (childTagFlag & MarshallerTagFlagMultilineStrings) == MarshallerTagFlagMultilineStrings {
							lineBreakAfterKey = string(LF)
						}
					}
				case reflect.Ptr:
					switch fieldRef.Type().Elem().Kind() {
					case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
						lineBreakAfterKey = string(Space)

						lines := strings.Split(fieldRef.String(), string(LF))
						if len(lines) > 1 {
							lineBreakAfterKey = string(LF)
						} else {
							if (childTagFlag & MarshallerTagFlagMultilineStrings) == MarshallerTagFlagMultilineStrings {
								lineBreakAfterKey = string(LF)
							}
						}
					default:
						lineBreakAfterKey = string(LF)
					}
				default:
					lineBreakAfterKey = string(LF)
				}

				marshalizedValue, exists := unmarshal(fieldType, &fieldRef, depth+1, childTagFlag)

				if !exists {
					if (childTagFlag & MarshallerTagFlagOmitEmpty) == MarshallerTagFlagOmitEmpty {
						continue
					}
					marshalizedValue = ""
				}
				result += fmt.Sprintf("%s%s:%s%s", fmt.Sprintf("%*s", depth*UnmarshalDefaultIndentSize, ""), key, lineBreakAfterKey, marshalizedValue)
			}
			return result, true
		}
	case reflect.Ptr:
		{
			if ref.IsNil() {
				return "", false
			}

			elem := ref.Elem()
			return unmarshal(typ.Elem(), &elem, depth, 0)
		}
	}
	return "", false
}

func getTagFlagFromTagValue(tagValues []string) (flag int) {
	for i := 1; i < len(tagValues); i++ {
		switch tagValues[i] {
		case MarshallerTagMultilineStrings:
			flag |= MarshallerTagFlagMultilineStrings
		case MarshallerTagOmitEmpty:
			flag |= MarshallerTagFlagOmitEmpty
		}
	}

	return
}
