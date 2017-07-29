// Binary XML marshaller derived from https://golang.org/src/encoding/xml/marshal.go

package binaryxml

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"github.com/cevaris/ordered_map"
	"io"
	"reflect"
	"strconv"
)

type BinaryXMLEncoder struct {
	writer io.Writer
}

func NewEncoder(writer io.Writer) *BinaryXMLEncoder {
	return &BinaryXMLEncoder{writer: writer}
}

func (encoder *BinaryXMLEncoder) Encode(v interface{}) error {
	table := ordered_map.NewOrderedMap()
	if err := generateElementNameDictionaryForValue(reflect.ValueOf(v), table); err != nil {
		return err
	}
	if err := encoder.writeTable(table); err != nil {
		return err
	}
	return encoder.writeSerial(reflect.ValueOf(v), nil, table)
}

var (
	marshalerType = reflect.TypeOf((*xml.Marshaler)(nil)).Elem()
)

func (encoder *BinaryXMLEncoder) marshalValue(val reflect.Value, finfo *fieldInfo, startTemplate *xml.StartElement, table *ordered_map.OrderedMap) error {
	if startTemplate != nil && startTemplate.Name.Local == "" {
		return fmt.Errorf("binaryxml: Encoding is missing name for StartElement")
	}
	if !val.IsValid() {
		return nil
	}
	if finfo != nil && finfo.flags&fOmitEmpty != 0 && isEmptyValue(val) {
		return nil
	}

	// Drill into interfaces and pointers
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	kind := val.Kind()
	typ := val.Type()

	// Slices and arrays iterate over the elements. They do not have an enclosing tag.
	if (kind == reflect.Slice || kind == reflect.Array) && typ.Elem().Kind() != reflect.Uint8 {
		for i, n := 0, val.Len(); i < n; i++ {
			if err := encoder.marshalValue(val.Index(i), finfo, startTemplate, table); err != nil {
				return err
			}
		}
		return nil
	}

	tinfo, err := getTypeInfo(typ)
	if err != nil {
		return err
	}

	// Create start element.
	// Precedence for the XML element name is:
	// 1. XMLName field in underlying struct;
	// 2. field name/tag in the struct field; and
	// 3. type name

	var start xml.StartElement
	if startTemplate != nil {
		start.Name = startTemplate.Name
	} else if tinfo.xmlname != nil {
		xmlname := tinfo.xmlname
		if xmlname.name != "" {
			start.Name.Space, start.Name.Local = xmlname.xmlns, xmlname.name
		} else if v, ok := xmlname.value(val).Interface().(xml.Name); ok && v.Local != "" {
			start.Name = v
		}
	}
	if start.Name.Local == "" && finfo != nil {
		start.Name.Space, start.Name.Local = finfo.xmlns, finfo.name
	}
	if start.Name.Local == "" {
		name := typ.Name()
		if name == "" {
			return &xml.UnsupportedTypeError{typ}
		}
		start.Name.Local = name
	}

	// Write open element
	{
		if x, ok := table.Get(start.Name.Local); ok {
			elementNumber := uint16(x.(int))
			binary.Write(encoder.writer, binary.BigEndian, nodetype)
			binary.Write(encoder.writer, binary.BigEndian, elementNumber)
		} else {
			return fmt.Errorf("binaryxml: failed looking up elementNumber for %s", start.Name.Local)
		}
	}

	// Attributes
	for i := range tinfo.fields {
		finfo := &tinfo.fields[i]
		// 		if finfo.flags&fAttr == 0 {
		// 			fmt.Printf("*skipping because fAttr==0\n")
		// 			continue
		// 		}
		fv := finfo.value(val)
		if finfo.flags&fOmitEmpty != 0 && isEmptyValue(fv) {
			continue
		}
		if fv.Kind() == reflect.Interface && fv.IsNil() {
			continue
		}

		name := xml.Name{Space: finfo.xmlns, Local: finfo.name}
		if err := encoder.marshalAttr(&start, name, finfo, fv, table); err != nil {
			return err
		}
	}

	// Write close element
	binary.Write(encoder.writer, binary.BigEndian, endtagtype)

	return nil
}

func (encoder *BinaryXMLEncoder) marshalAttr(start *xml.StartElement, name xml.Name, finfo *fieldInfo, val reflect.Value, table *ordered_map.OrderedMap) error {
	writer := encoder.writer

	var elementNumber uint16
	if x, ok := table.Get(name.Local); ok {
		elementNumber = uint16(x.(int))
	} else {
		return fmt.Errorf("No table entry for element %s", name.Local)
	}

	switch val.Kind() {
	case reflect.Bool:
		binary.Write(writer, binary.BigEndian, strtype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		writer.Write([]byte(strconv.FormatBool(val.Bool())))
		writer.Write([]byte("\x00"))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Float32:
		binary.Write(writer, binary.BigEndian, float4type)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, float32(val.Float()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Int8:
		binary.Write(writer, binary.BigEndian, int1btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, int8(val.Int()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Int16:
		binary.Write(writer, binary.BigEndian, int2btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, int16(val.Int()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Int32:
		binary.Write(writer, binary.BigEndian, int4btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, int32(val.Int()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Int64:
		binary.Write(writer, binary.BigEndian, int8btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, int64(val.Int()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Uint8:
		binary.Write(writer, binary.BigEndian, uint1btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, uint8(val.Uint()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Uint16:
		binary.Write(writer, binary.BigEndian, uint2btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, uint16(val.Uint()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Uint32:
		binary.Write(writer, binary.BigEndian, uint4btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, uint32(val.Uint()))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.String:
		binary.Write(writer, binary.BigEndian, strtype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		writer.Write([]byte(val.String()))
		writer.Write([]byte("\x00"))
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Uint64:
		binary.Write(writer, binary.BigEndian, uint8btype)
		binary.Write(writer, binary.BigEndian, elementNumber)
		binary.Write(writer, binary.BigEndian, val.Uint())
		binary.Write(writer, binary.BigEndian, endtagtype)
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 { // binary
			binary.Write(writer, binary.BigEndian, binarytype)
			binary.Write(writer, binary.BigEndian, elementNumber)
			binary.Write(writer, binary.BigEndian, uint32(val.Len()))
			binary.Write(writer, binary.BigEndian, val.Bytes())
			binary.Write(writer, binary.BigEndian, endtagtype)
		} else { // Walk slices of nested elements
			for i, n := 0, val.Len(); i < n; i++ {
				var start xml.StartElement
				start.Name.Local = finfo.name
				if err := encoder.marshalValue(val.Index(i), finfo, &start, table); err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		var startElement xml.StartElement
		startElement.Name.Local = finfo.name
		if err := encoder.marshalValue(val, finfo, &startElement, table); err != nil {
			return err
		}
	}

	return nil
}

func (encoder *BinaryXMLEncoder) writeTable(table *ordered_map.OrderedMap) error {
	// Write table begin marker
	if err := binary.Write(encoder.writer, binary.BigEndian, tablebegin); err != nil {
		return err
	}

	// Write table length
	tableLength := uint16(table.Len())
	if err := binary.Write(encoder.writer, binary.BigEndian, tableLength); err != nil {
		return err
	}

	// Write table, which is already sorted by value
	iter := table.IterFunc()
	for kv, ok := iter(); ok; kv, ok = iter() {
		encoder.writer.Write([]byte(kv.Key.(string)))
		encoder.writer.Write([]byte("\x00"))
	}

	// Write table end marker
	if err := binary.Write(encoder.writer, binary.BigEndian, tableend); err != nil {
		return err
	}

	return nil
}

func (encoder *BinaryXMLEncoder) writeSerial(val reflect.Value, finfo *fieldInfo, table *ordered_map.OrderedMap) error {
	// Write serial begin marker
	if err := binary.Write(encoder.writer, binary.BigEndian, serialbegin); err != nil {
		return err
	}

	// Write serial section
	if err := encoder.marshalValue(val, nil, nil, table); err != nil {
		return err
	}

	// Write serial end marker
	if err := binary.Write(encoder.writer, binary.BigEndian, serialend); err != nil {
		return err
	}

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
