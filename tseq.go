package jsonstream

import (
	"bufio"
	"encoding/json"
	"io"
)

// SeqReader reads JSON  documents delimited with a record separator, by default 0x1e
type SeqReader struct {
	sep byte
	scn *bufio.Scanner
}

// NewSeqReader returns a new JSON sequence reader. It uses a bufio.Scanner.
func NewSeqReader(r io.Reader) SeqReader {
	ret := SeqReader{sep: 0x1e, scn: bufio.NewScanner(r)}
	ret.scn.Split(bufio.ScanBytes)
	return ret
}

// NewSeqReaderWithSep returns a new JSON sequence reader with the
// given separator
func NewSeqReaderWithSep(r io.Reader, sep byte) SeqReader {
	ret := NewSeqReader(r)
	ret.sep = sep
	return ret
}

// ReadRaw reads the next raw document. The returned buffer is a copy
// of the read buffer, so subsequent read calls will not overwrite the
// buffer. If the end if the stream is reached, returns io.EOF. This
// does not validate if the returned document is valid JSON. It simply
// returns the input until the next separator is seen
func (r SeqReader) ReadRaw() ([]byte, error) {
	out := make([]byte, 0)
	for r.scn.Scan() {
		c := r.scn.Bytes()[0]
		if c == r.sep {
			return out, nil
		}
		out = append(out, c)
	}
	err := r.scn.Err()
	if err != nil {
		return nil, err
	}
	return nil, io.EOF
}

// Unmarshal the next  document from the input. Returns unmarshal
// error if there is any. If the end of stream is reached, returns
// io.EOF
func (r SeqReader) Unmarshal(out interface{}) error {
	o, err := r.ReadRaw()
	if err != nil {
		return err
	}
	return json.Unmarshal(o, &out)
}

// SeqWriter streams JSON documents delimited with a separator byte
type SeqWriter struct {
	sep      byte
	w        io.Writer
	notfirst bool
}

// NewSeqWriter returns a writer with 0x1e as separator
func NewSeqWriter(w io.Writer) SeqWriter {
	return SeqWriter{w: w, sep: 0x1e}
}

// NewSeqWriterSep returns a writer with a custom separator
func NewSeqWriterSep(w io.Writer, sep byte) SeqWriter {
	return SeqWriter{w: w, sep: sep}
}

// WriteRaw writes data to output. If this is not the first document
// written, then it adds the separator byte.
func (w SeqWriter) WriteRaw(data []byte) error {
	if w.notfirst {
		if _, err := w.w.Write([]byte{w.sep}); err != nil {
			return err
		}
	}
	w.notfirst = true
	_, err := w.w.Write(data)
	return err
}

// Marshal writes data to output.
func (w SeqWriter) Marshal(data interface{}) error {
	x, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.WriteRaw(x)
}
