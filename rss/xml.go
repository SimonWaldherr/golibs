package rss

import (
	"encoding/xml"
	"reflect"
	"strings"
)

var (
	mainXMLFields = feedXMLFields(reflect.TypeOf(Main{}))
	itemXMLFields = feedXMLFields(reflect.TypeOf(Item{}))
)

// UnmarshalXML decodes RSS channel fields alongside namespaced podcast fields.
// Go 1.18's default decoder rejects the image and itunes:image tags as a conflict.
func (feed *Main) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return decodeFeedFields(decoder, start, reflect.ValueOf(feed).Elem(), mainXMLFields)
}

// UnmarshalXML keeps RSS title and itunes:title distinct on all supported Go versions.
func (item *Item) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return decodeFeedFields(decoder, start, reflect.ValueOf(item).Elem(), itemXMLFields)
}

// The feed structs use simple element tags, optionally preceded by a namespace.
// Cache their field indexes once rather than inspecting tags for every item.
func feedXMLFields(typ reflect.Type) map[xml.Name]int {
	fields := make(map[xml.Name]int, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("xml")
		name := xml.Name{Local: tag}
		if separator := strings.IndexByte(tag, ' '); separator >= 0 {
			name.Space, name.Local = tag[:separator], tag[separator+1:]
		}
		fields[name] = i
	}
	return fields
}

func decodeFeedFields(decoder *xml.Decoder, start xml.StartElement, value reflect.Value, fields map[xml.Name]int) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch element := token.(type) {
		case xml.StartElement:
			index, found := fields[element.Name]
			if !found {
				// Like encoding/xml, an unqualified tag matches any namespace,
				// but a namespace-specific field takes precedence when present.
				index, found = fields[xml.Name{Local: element.Name.Local}]
			}
			if !found {
				if err := decoder.Skip(); err != nil {
					return err
				}
				continue
			}
			field := value.Field(index)
			if field.Kind() == reflect.Slice {
				entry := reflect.New(field.Type().Elem())
				if err := decoder.DecodeElement(entry.Interface(), &element); err != nil {
					return err
				}
				field.Set(reflect.Append(field, entry.Elem()))
			} else if err := decoder.DecodeElement(field.Addr().Interface(), &element); err != nil {
				return err
			}
		case xml.EndElement:
			if element.Name == start.Name {
				return nil
			}
		}
	}
}
