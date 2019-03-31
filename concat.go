package jsonstream

import (
	"encoding/json"
	"io"
)

// ConcatReader reads concatenated documents. The input must be valid
// documents concatenated one after the other. This uses a
// json.Decoder to read the stream.
type ConcatReader struct {
	d   *json.Decoder
	err error
}

// NewConcatReader returns a new stream reader
func NewConcatReader(r io.Reader) ConcatReader {
	return ConcatReader{d: json.NewDecoder(r)}
}

// ReadRaw reads the next raw  document. The returned buffer is a
// copy of the read buffer, so subsequent read calls will not
// overwrite the buffer. If the end of the stream is reached, returns
// io.EOF. If another error is detected (such as invalid
// document), the stream is tagged invalid and all subsequent calls
// will fail with the same error.
func (r ConcatReader) ReadRaw() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	var m json.RawMessage
	r.err = r.d.Decode(&m)
	if r.err != nil {
		return nil, r.err
	}
	return []byte(m), nil
}

// Unmarshal the next document from the input. Returns unmarshal error
// if there is any. If the end of stream is reached, returns
// io.EOF. If there is an error during unmarshaling, the rest of the
// stream cannot be recovered.
func (r ConcatReader) Unmarshal(out interface{}) error {
	if r.err != nil {
		return r.err
	}
	r.err = r.d.Decode(out)
	return r.err
}

// ConcatWriter streams JSON documents by concatenating one after the
// other.
type ConcatWriter struct {
	w io.Writer
}

// NewConcatWriter returns a new writer
func NewConcatWriter(w io.Writer) ConcatWriter {
	return ConcatWriter{w: w}
}

// WriteRaw writes a raw JSON document. It does not validate if the
// doc is valid.
func (w ConcatWriter) WriteRaw(out []byte) error {
	_, err := w.w.Write(out)
	return err
}

// Marshal data as a JSON document to the output
func (w ConcatWriter) Marshal(data interface{}) error {
	x, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.WriteRaw(x)
}
