package httpx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
)

// Marshaller interface for json or xml marshaller
type Marshaller interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// XML implement Marshaller interface
type XML struct{}

// Marshal use xml.Marshal
func (XML) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// Unmarshal use xml.Unmarshal
func (XML) Unmarshal(b []byte, v interface{}) error {
	return xml.Unmarshal(bytes.ToValidUTF8(b, []byte("")), v)
}

// JSON implement Marshaller interface
type JSON struct{}

// Marshal use json.Marshal
func (JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal use json.Unmarshal
func (JSON) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(bytes.ToValidUTF8(b, []byte("")), v)
}
