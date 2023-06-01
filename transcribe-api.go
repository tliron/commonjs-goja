package commonjs

import (
	"fmt"
	"io"
	"strings"

	"github.com/beevik/etree"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/yamlkeys"
)

type TranscribeAPI struct{}

func (self TranscribeAPI) ValidateFormat(code string, format string) error {
	return transcribe.Validate(code, format)
}

func (self TranscribeAPI) Decode(code string, format string, all bool) (ard.Value, error) {
	switch format {
	case "yaml", "":
		if all {
			if value, err := yamlkeys.DecodeAll(strings.NewReader(code)); err == nil {
				value_, _ := ard.NormalizeStringMaps(value)
				return value_, err
			} else {
				return nil, err
			}
		} else {
			value, _, err := ard.DecodeYAML(code, false)
			value, _ = ard.NormalizeStringMaps(value)
			return value, err
		}

	case "json":
		value, _, err := ard.DecodeJSON(code, false)
		value, _ = ard.NormalizeStringMaps(value)
		return value, err

	case "cjson":
		value, _, err := ard.DecodeCompatibleJSON(code, false)
		value, _ = ard.NormalizeStringMaps(value)
		return value, err

	case "xml":
		value, _, err := ard.DecodeCompatibleXML(code, false)
		value, _ = ard.NormalizeStringMaps(value)
		return value, err

	case "cbor":
		value, _, err := ard.DecodeCBOR(code)
		value, _ = ard.NormalizeStringMaps(value)
		return value, err

	case "messagepack":
		value, _, err := ard.DecodeMessagePack(code)
		value, _ = ard.NormalizeStringMaps(value)
		return value, err

	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
}

func (self TranscribeAPI) Encode(value any, format string, indent string, writer io.Writer) (string, error) {
	if writer == nil {
		return transcribe.Encode(value, format, indent, false)
	} else {
		err := transcribe.Write(value, format, indent, false, writer)
		return "", err
	}
}

func (self TranscribeAPI) NewXMLDocument() *etree.Document {
	return etree.NewDocument()
}
