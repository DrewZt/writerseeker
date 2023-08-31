package writerseeker

import (
	"bytes"
	"errors"
	"io"
)

// WriterSeeker is an in-memory io.WriteSeeker implementation
type WriterSeeker struct {
	Buf bytes.Buffer
	pos int
}

// Write writes to the Buffer of this WriterSeeker instance
func (ws *WriterSeeker) Write(p []byte) (n int, err error) {
	// If the offset is past the end of the Buffer, grow the Buffer with null bytes.
	if extra := ws.pos - ws.Buf.Len(); extra > 0 {
		if _, err := ws.Buf.Write(make([]byte, extra)); err != nil {
			return n, err
		}
	}

	// If the offset isn't at the end of the Buffer, write as much as we can.
	if ws.pos < ws.Buf.Len() {
		n = copy(ws.Buf.Bytes()[ws.pos:], p)
		p = p[n:]
	}

	// If there are remaining bytes, append them to the Buffer.
	if len(p) > 0 {
		var bn int
		bn, err = ws.Buf.Write(p)
		n += bn
	}

	ws.pos += n
	return n, err
}

// Seek seeks in the Buffer of this WriterSeeker instance
func (ws *WriterSeeker) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = ws.pos + offs
	case io.SeekEnd:
		newPos = ws.Buf.Len() + offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	ws.pos = newPos
	return int64(newPos), nil
}

// Reader returns an io.Reader. Use it, for example, with io.Copy, to copy the content of the WriterSeeker Buffer to an io.Writer
func (ws *WriterSeeker) Reader() io.Reader {
	return bytes.NewReader(ws.Buf.Bytes())
}

// Close :
func (ws *WriterSeeker) Close() error {
	return nil
}

// BytesReader returns a *bytes.Reader. Use it when you need a reader that implements the io.ReadSeeker interface
func (ws *WriterSeeker) BytesReader() *bytes.Reader {
	return bytes.NewReader(ws.Buf.Bytes())
}
