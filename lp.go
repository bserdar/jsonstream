package jsonstream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

// LenPrefixedReader reads  documents delimited with the
// length of the next record. The input looks like this:
//
//  18{"some":"thing\n"}55{"may":{"include":"nested","objects":["and","arrays"]}}
//
// This uses a bufio.Scanner to read the stream.
type LenPrefixedReader struct {
	scn *bufio.Scanner
	err error
}

// NewLenPrefixedReader returns a new reader
func NewLenPrefixedReader(r io.Reader) LenPrefixedReader {
	ret := LenPrefixedReader{scn: bufio.NewScanner(r)}
	ret.scn.Split(bufio.ScanBytes)
	return ret
}

// ReadRaw reads the next raw document. The returned buffer is a copy
// of the read buffer, so subsequent read calls will not overwrite the
// buffer. If the end if the stream is reached, returns io.EOF. If
// another error is detected, stream is put into error state and
// subsequent calls will fail.
func (r LenPrefixedReader) ReadRaw() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	sz := make([]byte, 0, 10)
	out := make([]byte, 0)
	readingLen := true
	len := 0
	var err error
	for r.scn.Scan() {
		if readingLen {
			c := r.scn.Bytes()[0]
			if unicode.IsDigit(rune(c)) {
				sz = append(sz, c)
			} else {
				readingLen = false
				len, err = strconv.Atoi(string(sz))
				if err != nil {
					r.err = err
					return nil, err
				}
				sz = make([]byte, 0, 10)
				out = append(out, c)
				len--
			}
		} else {
			if len > 0 {
				out = append(out, r.scn.Bytes()[0])
				len--
				if len == 0 {
					return out, nil
				}
			}
		}
	}
	err = r.scn.Err()
	if err != nil {
		r.err = err
		return nil, err
	}
	return nil, io.EOF
}

// Unmarshal the next  document from the input. Returns unmarshal
// error if there is any. If the end of stream is reached, returns
// io.EOF
func (r LenPrefixedReader) Unmarshal(out interface{}) error {
	o, err := r.ReadRaw()
	if err != nil {
		return err
	}
	return json.Unmarshal(o, &out)
}

// LenPrefixedWriter streams JSON documents by adding the byte length
// of the JSON document before it.
type LenPrefixedWriter struct {
	w io.Writer
}

// NewLenPrefixedWriter returns a new writer
func NewLenPrefixedWriter(w io.Writer) LenPrefixedWriter {
	return LenPrefixedWriter{w: w}
}

// WriteRaw writes the length of data and then data. It assumes that
// the 'data' is a valid JSON document. If len(data)==0, returns
// error. Don't pass empty documents.
func (w LenPrefixedWriter) WriteRaw(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("Empty doc")
	}
	_, err := w.w.Write([]byte(strconv.Itoa(len(data))))
	if err != nil {
		return err
	}
	_, err = w.w.Write(data)
	return err
}

// Marshal writes the object as length-prefixed JSON to output
func (w LenPrefixedWriter) Marshal(data interface{}) error {
	x, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.WriteRaw(x)
}
