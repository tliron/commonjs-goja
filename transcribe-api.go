package commonjs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-transcribe"
	"github.com/tliron/yamlkeys"
)

type TranscribeAPI struct{}

func (self TranscribeAPI) ValidateFormat(code []byte, format string) error {
	return transcribe.Validate(string(code), format)
}

func (self TranscribeAPI) Decode(code []byte, format string, all bool) (ard.Value, error) {
	switch format {
	case "yaml":
		if all {
			if value, err := yamlkeys.DecodeAll(bytes.NewReader(code)); err == nil {
				value_, _ := ard.ConvertMapsToStringMaps(value)
				return value_, nil
			} else {
				return nil, err
			}
		} else {
			if value, _, err := ard.DecodeYAML(code, false); err == nil {
				value, _ = ard.ConvertMapsToStringMaps(value)
				return value, nil
			} else {
				return nil, err
			}
		}

	case "json":
		return ard.DecodeJSON(code, true)

	case "xjson":
		return ard.DecodeXJSON(code, true)

	case "xml":
		if value, err := ard.DecodeXML(code); err == nil {
			value, _ = ard.ConvertMapsToStringMaps(value)
			return value, nil
		} else {
			return nil, err
		}

	case "cbor":
		if value, err := ard.DecodeCBOR(code); err == nil {
			value, _ = ard.ConvertMapsToStringMaps(value)
			return value, nil
		} else {
			return nil, err
		}

	case "messagepack":
		return ard.DecodeMessagePack(code, true)

	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
}

func (self TranscribeAPI) Encode(value any, format string, indent string, writer io.Writer) (string, error) {
	if writer == nil {
		return transcribe.Stringify(value, format, indent, false, nil)
	} else {
		err := transcribe.Write(value, format, indent, false, writer, false, nil)
		return "", err
	}
}

func (self TranscribeAPI) NewXMLDocument() *etree.Document {
	return etree.NewDocument()
}
