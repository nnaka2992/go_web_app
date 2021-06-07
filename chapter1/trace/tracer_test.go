package trace

import (
	"bytes"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("New returns nil")
	} else {
		tracer.Trace("Hello Tracer package")
		if buf.String() != "Hello Tracer package\n" {
			t.Errorf("'%s' is outputted", buf.String())
		}
	}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
