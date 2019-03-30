package jsonstream

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"testing"
)

type nestedType struct {
	Value string `json:"value"`
}

type marshalTestType struct {
	Str    string     `json:"str"`
	Int    int        `json:"int"`
	Struct nestedType `json:"struct"`
}

var testData1 = []marshalTestType{{Str: "str1", Int: 1, Struct: nestedType{Value: "v1"}},
	{Str: "str2", Int: 1, Struct: nestedType{Value: "v2"}},
	{Str: "str3", Int: 1, Struct: nestedType{Value: "v3"}},
	{Str: "str4", Int: 1, Struct: nestedType{Value: "v4"}}}

func testRaw(reader Reader, raw [][]byte, t *testing.T) {
	for _, x := range raw {
		b, err := reader.ReadRaw()
		if err != nil {
			t.Errorf("Cannot unmarshal: %v", err)
		}

		if bytes.Compare(b, x) != 0 {
			t.Errorf("Expected %s got %s", string(b), string(x))
		}
	}
	// Must get eof
	_, err := reader.ReadRaw()
	if err != io.EOF {
		t.Errorf("No eof")
	}
}

func TestND(t *testing.T) {
	str := bytes.Buffer{}
	raw := make([][]byte, 0)
	for _, d := range testData1 {
		x, _ := json.Marshal(d)
		raw = append(raw, x)
		str.Write(x)
		str.WriteRune('\n')
	}

	reader := NewJSONLReader(bytes.NewReader(str.Bytes()))
	var x marshalTestType
	for _, d := range testData1 {
		err := reader.Unmarshal(&x)
		if err != nil {
			t.Errorf("Cannot unmarshal: %v", err)
		}
		if x != d {
			t.Errorf("Expected %+v got %+v", d, x)
		}
	}
	// Must get eof
	err := reader.Unmarshal(&x)
	if err != io.EOF {
		t.Errorf("No eof")
	}

	testRaw(NewJSONLReader(bytes.NewReader(str.Bytes())), raw, t)
}

func TestConcat(t *testing.T) {
	str := bytes.Buffer{}
	raw := make([][]byte, 0)
	for _, d := range testData1 {
		x, _ := json.Marshal(d)
		raw = append(raw, x)
		str.Write(x)
	}

	reader := NewJSONConcatReader(bytes.NewReader(str.Bytes()))
	var x marshalTestType
	for _, d := range testData1 {
		err := reader.Unmarshal(&x)
		if err != nil {
			t.Errorf("Cannot unmarshal: %v", err)
		}
		if x != d {
			t.Errorf("Expected %+v got %+v", d, x)
		}
	}
	// Must get eof
	err := reader.Unmarshal(&x)
	if err != io.EOF {
		t.Errorf("No eof")
	}

	testRaw(NewJSONConcatReader(bytes.NewReader(str.Bytes())), raw, t)
}

func TestSeq(t *testing.T) {
	str := bytes.Buffer{}
	raw := make([][]byte, 0)
	for _, d := range testData1 {
		x, _ := json.Marshal(d)
		raw = append(raw, x)
		str.Write(x)
		str.WriteRune('\x1e')
	}

	reader := NewJSONSeqReader(bytes.NewReader(str.Bytes()))
	var x marshalTestType
	for _, d := range testData1 {
		err := reader.Unmarshal(&x)
		if err != nil {
			t.Errorf("Cannot unmarshal: %v", err)
		}
		if x != d {
			t.Errorf("Expected %+v got %+v", d, x)
		}
	}
	// Must get eof
	err := reader.Unmarshal(&x)
	if err != io.EOF {
		t.Errorf("No eof")
	}

	testRaw(NewJSONSeqReader(bytes.NewReader(str.Bytes())), raw, t)
}

func TestLenPrefixed(t *testing.T) {
	str := bytes.Buffer{}
	raw := make([][]byte, 0)
	for _, d := range testData1 {
		x, _ := json.Marshal(d)
		raw = append(raw, x)
		str.WriteString(strconv.Itoa(len(x)))
		str.Write(x)
	}

	reader := NewJSONLenPrefixedReader(bytes.NewReader(str.Bytes()))
	var x marshalTestType
	for _, d := range testData1 {
		err := reader.Unmarshal(&x)
		if err != nil {
			t.Errorf("Cannot unmarshal: %v", err)
		}
		if x != d {
			t.Errorf("Expected %+v got %+v", d, x)
		}
	}
	// Must get eof
	err := reader.Unmarshal(&x)
	if err != io.EOF {
		t.Errorf("No eof")
	}

	testRaw(NewJSONLenPrefixedReader(bytes.NewReader(str.Bytes())), raw, t)
}
