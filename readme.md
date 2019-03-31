# JSON Streams

[![GoDoc](https://godoc.org/github.com/bserdar/jsonstream?status.svg)](https://godoc.org/github.com/bserdar/jsonstream)

This library supports streaming JSON streaming conventions described
in this [Wikipedia page](https://en.wikipedia.org/wiki/JSON_streaming).

This is for the concatenated JSON streams, where each JSON document is
concatenated one after the other:

```
rd:=jsonstream.NewConcatReader(reader)
var entry MyEntry
for {
  err:=rd.Unmarshal(&entry)
  if err==io.EOF {
    break
  }
  if err!=nil {
    return err
  }
  processEntry(entry)
}
```

This is for JSON streams with a separate JSON document in each line
(NDJSON):

```
ndLinesReader:=jsonstream.NewLineReader(reader)
ndLinesWriter:=jsonstream.NewLineWriter(writer)
```

This is for JSON streams separated with record separator delimiter:

```
seqReader:=jsonstream.NewSeqReader(reader) // 0x1e record separator
seqWriter:=jsonstream.NewSeqWriter(writer) 
```
```
seqReader:=jsonstream.NewSeqReaderWithSep(reader,'\n') // Line-separated JSON
seqWriter:=jsonstream.NewSeqWriterWithSep(writer,'\n')
```

This is for JSON streams with length prefixed JSON stream, where each
JSON document is prefixed by its byte length:

```
lpReader:=jsonstream.NewLenPrefixedReader(reader)
lpWriter:=jsonstream.NewLenPrefixesWriter(writer)
```
## APIs

All four stream readers/writers use the same APIs.

### Readers

```
data, err:=reader.ReadRaw()
```

ReadRaw reads the next JSON document. Only the ConcatReader validates
that the JSON document is a valid document, the remaining readers
simply read until the next delimiter. The returned byte array is a
newly allocated copy of the underlying read buffer. Some of the
readers use buffered input, so the state of the underlying reader is
unknown.

```
var data myStruct
err:=reader.Unmarshal(&data)
```

Unmarshals the next entry from the input. For concatenated JSON,
errors invalidate the rest of the stream. For others stream processing
can continue.


### Writers

```
err:=writer.WriteRaw(data)
```

WriteRaw simply writes the []byte data to the output, with the correct
delimiter. For NDJSON, WriteRaw removes the newline characters from
data.

```
err:=writer.Marshal(data)
```

Marshal first encodes data to JSON, and then writes it to the output.
