package trace

import (
	"bytes"
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
