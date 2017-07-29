package binaryxml

import (
	"bytes"
	"encoding/xml"
	"reflect"

	"github.com/cevaris/ordered_map"
)

func generateElementNameDictionaryForValue(value reflect.Value, table *ordered_map.OrderedMap) error {
	if !value.IsValid() {
		return nil
	}

	// Drill into interfaces and pointers
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	typeInfo, err := getTypeInfo(value.Type())
	if err != nil {
		return err
	}

	// Attributes
	for i := range typeInfo.fields {
		fieldInfo := &typeInfo.fields[i]
		fieldValue := fieldInfo.value(value)
		if fieldInfo.flags&fOmitEmpty != 0 && isEmptyValue(fieldValue) {
			continue
		}
		if fieldValue.CanInterface() == fieldValue.Type().Implements(marshalerType) {
			xmlBytes, err := xml.Marshal(fieldValue.Interface())
			if err != nil {
				return err
			}

			// Traverse generated XML snippet and populate dictionary with element names
			buffer := bytes.NewBuffer(xmlBytes)
			decoder := xml.NewDecoder(buffer)
			var node xmlTraversalNode
			if err := decoder.Decode(&node); err != nil {
				return err
			}
			walk([]xmlTraversalNode{node}, func(node xmlTraversalNode) bool {
				addIfNeeded(node.XMLName.Local, table)
				return true
			})
		}
		if fieldValue.Kind() == reflect.Interface && fieldValue.IsNil() {
			continue
		}

		addIfNeeded(fieldInfo.name, table)

		// Drill into nested structs
		if fieldValue.Kind() == reflect.Struct {
			generateElementNameDictionaryForValue(fieldValue, table)
		}

		// Drill into nested slices
		if fieldValue.Kind() == reflect.Slice {
			for i, n := 0, fieldValue.Len(); i < n; i++ {
				if err := generateElementNameDictionaryForValue(fieldValue.Index(i), table); err != nil {
					return err
				}
			}
		}
	}

	// Element name
	if typeInfo.xmlname != nil {
		xmlName := typeInfo.xmlname
		name := xmlName.name
		if name == "" {
			if v, ok := xmlName.value(value).Interface().(xml.Name); ok && v.Local != "" {
				name = v.Local
			}
		}
		addIfNeeded(name, table)
	}

	return nil
}

func addIfNeeded(elementName string, table *ordered_map.OrderedMap) {
	if elementNumber, ok := table.Get(elementName); !ok {
		elementNumber = table.Len() + 1
		table.Set(elementName, elementNumber)
	}
}
