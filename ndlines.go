package jsonstream

import (
	"bufio"
	"encoding/json"
	"io"
)

// LineReader reads one complete JSON document from each line.
type LineReader struct {
	scn *bufio.Scanner
}

// NewLineReader returns a new  lines reader. Each line must be a
// valid  document. This uses a bufio.Scanner to read lines.
func NewLineReader(r io.Reader) LineReader {
	return LineReader{scn: bufio.NewScanner(r)}
}

// ReadRaw reads the next line of the input. The returned buffer is a
// copy of the read buffer, so subsequent read calls will not
// overwrite the buffer. If the end if the stream is reached, returns
// io.EOF. It does not validate that the line is a valid JSON document.
func (r LineReader) ReadRaw() ([]byte, error) {
	if r.scn.Scan() {
		o := r.scn.Bytes()
		out := make([]byte, len(o))
		copy(out, o)
		return out, nil
	}
	err := r.scn.Err()
	if err == nil {
		return nil, io.EOF
	}
	return nil, err
}

// Unmarshal the next document from the input. Returns unmarshal error
// if there is any. If the end of stream is reached, returns
// io.EOF. This reader can continue reading lines even if the previous
// line had an error.
func (r LineReader) Unmarshal(out interface{}) error {
	if r.scn.Scan() {
		o := r.scn.Bytes()
		return json.Unmarshal(o, &out)
	}
	err := r.scn.Err()
	if err == nil {
		return io.EOF
	}
	return err
}

// LineWriter streams one JSON document every line
type LineWriter struct {
	w io.Writer
}

// NewLineWriter returns a new writer
func NewLineWriter(w io.Writer) LineWriter {
	return LineWriter{w: w}
}

// WriteRaw writes data followed by newline. Any new line characterts
// in data are removed.
func (w LineWriter) WriteRaw(data []byte) error {
	bw := bufio.NewWriter(w.w)
	for _, c := range data {
		if c != '\n' {
			if err := bw.WriteByte(c); err != nil {
				return err
			}
		}
	}
	if err := bw.WriteByte('\n'); err != nil {
		return err
	}
	return bw.Flush()
}

// Marshal writes data as a JSON document followed by newline.
func (w LineWriter) Marshal(data interface{}) error {
	x, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.WriteRaw(x)
}
