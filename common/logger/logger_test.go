package logger

import (
	"errors"
	"testing"
)

type Writer struct {
	nBytes int
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.nBytes = len(p)
	return len(p),nil
}

func (w *Writer) Reset() {
	w.nBytes = 0
}

func (w *Writer) GetWriteNBytes() int {
	return w.nBytes
}

func TestLogger(t *testing.T) {
	writer := &Writer{}
	logger := New(writer,ERROR)
	logger.Info("hello world")
	if writer.GetWriteNBytes() <= 0 {
		t.Error("log writer failed")
		return
	}
	writer.Reset()
	logger.Panic(errors.New("hello panic"))
	if writer.GetWriteNBytes() > 0 {
		t.Error("log writer failed")
		return
	}
}
