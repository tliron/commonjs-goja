package api

import (
	"io"
	"os"

	"github.com/beevik/etree"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-transcribe"
)

func CreateTranscribeExtension(stdout io.Writer, stderr io.Writer) commonjs.CreateExtensionFunc {
	return func(jsContext *commonjs.Context) any {
		return NewTranscribe(stdout, stderr)
	}
}

//
// Transcribe
//

type Transcribe struct {
	Stdout io.Writer
	Stderr io.Writer
}

func NewTranscribe(stdout io.Writer, stderr io.Writer) *Transcribe {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	return &Transcribe{
		Stdout: stdout,
		Stderr: stderr,
	}
}

// Encodes and writes the value. Supported formats are "yaml", "json",
// "xjson", "xml", "cbor", "messagepack", and "go". The "cbor" and "messsagepack"
// formats will ignore the indent argument.
func (self Transcribe) Write(writer io.Writer, value any, format string, indent string) error {
	transcriber := transcribe.Transcriber{
		Writer: writer,
		Format: format,
		Indent: indent,
	}

	return transcriber.Write(value)
}

// Encodes and writes the value. Supported formats are "yaml", "json",
// "xjson", "xml", "cbor", "messagepack", and "go". The "cbor" and "messsagepack"
// formats will be encoded in base64 and will ignore the indent argument.
func (self Transcribe) WriteText(writer io.Writer, value any, format string, indent string) error {
	transcriber := transcribe.Transcriber{
		Writer: writer,
		Format: format,
		Indent: indent,
		Base64: true,
	}

	return transcriber.Write(value)
}

// Encodes and prints the value to stdout. Supported formats are "yaml", "json",
// "xjson", "xml", "cbor", "messagepack", and "go". The "cbor" and "messsagepack"
// formats will be encoded in base64 and will ignore the indent argument.
func (self Transcribe) Print(value any, format string, indent string) error {
	transcriber := transcribe.Transcriber{
		Writer:      self.Stdout,
		Format:      format,
		ForTerminal: true,
		Indent:      indent,
		Base64:      true,
	}

	return transcriber.Write(value)
}

func (self Transcribe) Eprint(value any, format string, indent string) error {
	transcriber := transcribe.Transcriber{
		Writer:      self.Stderr,
		Format:      format,
		ForTerminal: true,
		Indent:      indent,
		Base64:      true,
	}

	return transcriber.Write(value)
}

// Encodes the value into a string.Supported formats are "yaml", "json",
// "xjson", "xml", "cbor", "messagepack", and "go". The "cbor" and "messsagepack"
// formats will be encoded in base64 and will ignore the indent argument.
func (self Transcribe) Stringify(value any, format string, indent string) (string, error) {
	transcriber := transcribe.Transcriber{
		Format: format,
		Indent: indent,
		Base64: true,
	}

	return transcriber.Stringify(value)
}

func (self Transcribe) NewXmlDocument() *etree.Document {
	return etree.NewDocument()
}
