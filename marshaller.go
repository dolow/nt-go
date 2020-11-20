package ntgo

import (
	"reflect"
)

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

		childDirective, exists := directive.Dictionary[tagValue]
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
				marshal(childDirective, fieldType, &fieldInstance)
				fieldRef.Set(fieldInstance)
			}
		}
	}
}

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
