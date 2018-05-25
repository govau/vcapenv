package vcapenv

import (
	"fmt"
	"io"
	"os"
)

type VariableExporter interface {
	Export(variable, value string)
}

type VariableExporterFunc func(variable, value string)

func (f VariableExporterFunc) Export(variable, value string) {
	f(variable, value)
}

type BashExporter struct {
	io.Writer
}

func (bash BashExporter) Export(variable, value string) {
	fmt.Fprintf(bash, "export %s=%s\n", variable, value)
}

func Setenv(variable, value string) {
	os.Setenv(variable, value)
}

var SetenvExporter VariableExporter = VariableExporterFunc(Setenv)
