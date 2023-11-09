package api

import (
	"bytes"
	"fmt"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-ard"
	"github.com/tliron/yamlkeys"
)

// ([commonjs.CreateExtensionFunc] signature)
func CreateARDExtension(jsContext *commonjs.Context) any {
	return ARD{}
}

//
// ARD
//

type ARD struct{}

func (self ARD) Decode(code []byte, format string, all bool) (ard.Value, error) {
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
		if value, err := ard.DecodeCBOR(code, false); err == nil {
			value, _ = ard.ConvertMapsToStringMaps(value)
			return value, nil
		} else {
			return nil, err
		}

	case "messagepack":
		return ard.DecodeMessagePack(code, false, true)

	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
}

func (self ARD) ValidateFormat(code []byte, format string) error {
	return ard.Validate(code, format)
}
