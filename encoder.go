package binaryxml

import (
	"encoding/binary"
	"encoding/xml"
	"io"
	"github.com/cevaris/ordered_map"
	"reflect"
)


type BinaryXMLEncoder struct {
	writer io.Writer
}


func (encoder *BinaryXMLEncoder) Encode(v interface{}) error {
	table := ordered_map.NewOrderedMap()
	if err := encoder.populateTable(reflect.ValueOf(v), nil, table); err != nil {return err}
	if err := encoder.writeTable(table); err != nil {return err}
	return encoder.writeSerial(reflect.ValueOf(v), nil, table)
}


func (encoder *BinaryXMLEncoder) populateTable(value reflect.Value, fieldInfo *fieldInfo, table *ordered_map.OrderedMap) error {
	if !value.IsValid() {return nil}
	
	// Drill into interfaces and pointers
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		if value.IsNil() {return nil}
		value = value.Elem()
	}
	
	typeInfo, err := getTypeInfo(value.Type())
	if err != nil {return err}
	
	// Attributes
	for i := range typeInfo.fields {
		fieldInfo := &typeInfo.fields[i]
		fv := fieldInfo.value(value)
		if fieldInfo.flags&fOmitEmpty != 0 && isEmptyValue(fv) {continue}
		if fv.Kind() == reflect.Interface && fv.IsNil() {continue}
		if elementNumber, ok := table.Get(fieldInfo.name); !ok {
			elementNumber = table.Len() +1
			table.Set(fieldInfo.name, elementNumber)
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
		elementNumber := table.Len() +1
		table.Set(name, elementNumber)
	}
	
	return nil
}


func (encoder *BinaryXMLEncoder) marshalValue(value reflect.Value, fieldInfo *fieldInfo, table *ordered_map.OrderedMap) error {
	if !value.IsValid() {return nil}
	
	// Drill into interfaces and pointers
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		if value.IsNil() {return nil}
		value = value.Elem()
	}
	kind := value.Kind()
	type_ := value.Type()
	
	// Slices and arrays iterate over the elements. They do not have an enclosing tag.
	if (kind == reflect.Slice || kind == reflect.Array) && type_.Elem().Kind() != reflect.Uint8 {
		for i, n := 0, value.Len(); i < n; i++ {
			if err := encoder.marshalValue(value.Index(i), fieldInfo, table); err != nil {return err}
		}
		return nil
	}
  	
	typeInfo, err := getTypeInfo(type_)
	if err != nil {return err}
	
	// Create start element.
	// Precedence for the XML element name is:
	// 1. XMLName field in underlying struct;
	// 2. field name/tag in the struct field; and
	// 3. type name
	
	var start xml.StartElement
	if typeInfo.xmlname != nil {
		xmlname := typeInfo.xmlname
		if xmlname.name != "" {
			start.Name.Space, start.Name.Local = xmlname.xmlns, xmlname.name
		} else if v, ok := xmlname.value(value).Interface().(xml.Name); ok && v.Local != "" {
			start.Name = v
		}
	}
	if start.Name.Local == "" && fieldInfo != nil {
		start.Name.Space, start.Name.Local = fieldInfo.xmlns, fieldInfo.name
	}
	if start.Name.Local == "" {
		name := type_.Name()
		if name == "" {
			return &xml.UnsupportedTypeError{type_}
		}
		start.Name.Local = name
	}
	
	// Write open element
	{
		elementNumber, _ := table.Get(start.Name.Local)
		binary.Write(encoder.writer, binary.BigEndian, nodetype)
		binary.Write(encoder.writer, binary.BigEndian, uint16(elementNumber.(int)))
	}
	
	// Attributes
	for i := range typeInfo.fields {
		fieldInfo := &typeInfo.fields[i]
// 		if fieldInfo.flags&fAttr == 0 {
// 			fmt.Printf("*skipping because fAttr==0\n")
// 			continue
// 		}
		fv := fieldInfo.value(value)	
		if fieldInfo.flags&fOmitEmpty != 0 && isEmptyValue(fv) {continue}
		if fv.Kind() == reflect.Interface && fv.IsNil() {continue}
	
		name := xml.Name{Space: fieldInfo.xmlns, Local: fieldInfo.name}
		if err := marshalAttr(&start, name, fv, encoder.writer, table); err != nil {return err}
	}
  	
	// Write close element
	binary.Write(encoder.writer, binary.BigEndian, endtagtype)
	
	return nil
}


func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}


func marshalAttr(start *xml.StartElement, name xml.Name, value reflect.Value, writer io.Writer, table *ordered_map.OrderedMap) error {
	x, _ := table.Get(name.Local); elementNumber := uint16(x.(int))
	switch value.Kind() {
		case reflect.String:
			binary.Write(writer, binary.BigEndian, strtype)
			binary.Write(writer, binary.BigEndian, elementNumber)
			writer.Write([]byte(value.String()))
			writer.Write([]byte("\x00"))
    case reflect.Int8:
      binary.Write(writer, binary.BigEndian, int1btype)
      binary.Write(writer, binary.BigEndian, elementNumber)
      binary.Write(writer, binary.BigEndian, int8(value.Int()))
    case reflect.Int16:
      binary.Write(writer, binary.BigEndian, int2btype)
      binary.Write(writer, binary.BigEndian, elementNumber)
      binary.Write(writer, binary.BigEndian, int16(value.Int()))
    case reflect.Int32:
      binary.Write(writer, binary.BigEndian, int4btype)
      binary.Write(writer, binary.BigEndian, elementNumber)
      binary.Write(writer, binary.BigEndian, int32(value.Int()))
    case reflect.Int64:
      binary.Write(writer, binary.BigEndian, int8btype)
      binary.Write(writer, binary.BigEndian, elementNumber)
      binary.Write(writer, binary.BigEndian, int64(value.Int()))
    case reflect.Uint8:
      binary.Write(writer, binary.BigEndian, uint1btype)
      binary.Write(writer, binary.BigEndian, elementNumber)
      binary.Write(writer, binary.BigEndian, uint8(value.Uint()))
    case reflect.Uint16:
      binary.Write(writer, binary.BigEndian, uint2btype)
      binary.Write(writer, binary.BigEndian, elementNumber)
      binary.Write(writer, binary.BigEndian, uint16(value.Uint()))
    case reflect.Uint32:
      binary.Write(writer, binary.BigEndian, uint4btype)
      binary.Write(writer, binary.BigEndian, elementNumber)
      binary.Write(writer, binary.BigEndian, uint32(value.Uint()))
		case reflect.Uint64:
			binary.Write(writer, binary.BigEndian, uint8btype)
			binary.Write(writer, binary.BigEndian, elementNumber)
			binary.Write(writer, binary.BigEndian, value.Uint())
	}
	
	// Write close element
	binary.Write(writer, binary.BigEndian, endtagtype)
	
	return nil
}


func (encoder *BinaryXMLEncoder) writeTable(table *ordered_map.OrderedMap) error {
	// Write table begin marker
	if err := binary.Write(encoder.writer, binary.BigEndian, tablebegin); err != nil {return err}
	
	// Write table length
	tableLength := uint16(table.Len())
	if err := binary.Write(encoder.writer, binary.BigEndian, tableLength); err != nil {return err}
	
	// Write table, which is already sorted by value
	iter := table.IterFunc()
	for kv, ok := iter(); ok; kv, ok = iter() {
		encoder.writer.Write([]byte(kv.Key.(string)))
		encoder.writer.Write([]byte("\x00"))
	}
	
	// Write table end marker
	if err := binary.Write(encoder.writer, binary.BigEndian, tableend); err != nil {return err}
	
	return nil
}


func (encoder *BinaryXMLEncoder) writeSerial(value reflect.Value, fieldInfo *fieldInfo, table *ordered_map.OrderedMap) error {
	// Write serial begin marker
	if err := binary.Write(encoder.writer, binary.BigEndian, serialbegin); err != nil {return err}
	
	// Write serial section
	if err := encoder.marshalValue(value, nil, table); err != nil {return err}
	
	// Write serial end marker
	if err := binary.Write(encoder.writer, binary.BigEndian, serialend); err != nil {return err}
	
	return nil
}
