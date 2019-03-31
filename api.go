// Package jsonstream deals with streaming JSON input/output. Supports
// line delimited (NDJSON) JSON,  length-prefixed JSON,
// record separator delimited JSON, and
// concatenated JSON streams.
//
// https://en.wikipedia.org/wiki/JSON_streaming
//
package jsonstream

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// Reader reads raw JSON data or unmarshals data from a JSON stream
type Reader interface {
	// ReadRaw reads the next raw JSON document. The returned buffer is a
	// copy of the read buffer, so subsequent read calls will not
	// overwrite the buffer. If the end if the stream is reached, returns
	// io.EOF
	ReadRaw() ([]byte, error)
	// Unmarshal the next JSON document from the input. Returns unmarshal
	// error if there is any. If the end of stream is reached, returns
	// io.EOF
	Unmarshal(out interface{}) error
}

// Writer writes raw JSON data or marshals data to a JSON stream. The
// writer keeps the state, so you can call the write functions one
// after the other and the correct delimiter will be inserted between
// records.
type Writer interface {
	// Write raw data to output stream.
	WriteRaw([]byte) error
	// Marshals a JSON document to output stream. This is a single
	// document, if you pass an array, it will be marshaled as a JSON
	// array, not as a series of documents.
	Marshal(interface{}) error
}

// ReadAll reads all lines of the stream. Returns error for read
// errors except io.EOF
func ReadRawAll(r Reader) ([][]byte, error) {
	out := make([][]byte, 0)
	for {
		o, err := r.ReadRaw()
		if err == io.EOF {
			break
		}
		if err != nil {
			return out, err
		}
		out = append(out, o)
	}
	return out, nil
}

// UnmarshalAll unmarshals all JSON documents from the input to the
// pointer to slice 'out'. The 'out' must be a pointer to a slice:
//
//   var out []MyStruct
//   jsonstream.UnmarshalAll(rdr,&out)
//
func UnmarshalAll(r Reader, out interface{}) error {
	ptr := reflect.ValueOf(out)
	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("Pointer to slice required")
	}
	ptrElem := ptr.Type().Elem()
	if ptrElem.Kind() != reflect.Slice {
		return fmt.Errorf("Pointer is not pointing to a slice")
	}
	elemType := ptrElem.Elem()
	for {
		o, err := r.ReadRaw()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		el := reflect.New(elemType)
		err = json.Unmarshal(o, el.Interface())
		if err != nil {
			return err
		}
		ptr.Elem().Set(reflect.Append(ptr.Elem(), el.Elem()))
	}
	return nil
}
