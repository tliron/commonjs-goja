package commonjs

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/util"
)

type UtilAPI struct{}

func (self TranscribeAPI) StringToBytes(string_ string) []byte {
	return util.StringToBytes(string_)
}

// Another way to achieve this in JavaScript: String.fromCharCode.apply(null, bytes)
func (self TranscribeAPI) BytesToString(bytes []byte) string {
	return util.BytesToString(bytes)
}

// Encode bytes as base64
func (self TranscribeAPI) Btoa(bytes []byte) string {
	return util.ToBase64(bytes)
}

// Decode base64 to bytes
func (self TranscribeAPI) Atob(b64 string) ([]byte, error) {
	return util.FromBase64(b64)
}

func (self UtilAPI) DeepCopy(value ard.Value) ard.Value {
	return ard.Copy(value)
}

func (self UtilAPI) DeepEquals(a ard.Value, b ard.Value) bool {
	return ard.Equals(a, b)
}

func (self UtilAPI) IsType(value ard.Value, type_ string) (bool, error) {
	// Special case whereby an integer stored as a float type has been optimized to an integer type
	if (type_ == "!!float") && util.IsInteger(value) {
		return true, nil
	}

	if validate, ok := ard.TypeValidators[ard.TypeName(type_)]; ok {
		return validate(value), nil
	} else {
		return false, fmt.Errorf("unsupported type: %q", type_)
	}
}

func (self UtilAPI) Hash(value ard.Value) (string, error) {
	if hash, err := hashstructure.Hash(value, hashstructure.FormatV2, nil); err == nil {
		return strconv.FormatUint(hash, 10), nil
	} else {
		return "", err
	}
}

func (self UtilAPI) Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (self UtilAPI) Sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func (self UtilAPI) Now() time.Time {
	return time.Now()
}

func (self UtilAPI) Mutex() util.RWLocker {
	return util.NewDefaultRWLocker()
}

var onces sync.Map

func (self UtilAPI) Once(name string, value goja.Value) error {
	if call, ok := goja.AssertFunction(value); ok {
		once, _ := onces.LoadOrStore(name, new(sync.Once))
		once.(*sync.Once).Do(func() {
			if _, err := call(nil); err != nil {
				log.Errorf("%s", err.Error())
			}
		})
		return nil
	} else {
		return fmt.Errorf("not a \"function\": %T", value)
	}
}

// Goroutine
func (self UtilAPI) Go(value goja.Value) error {
	if call, ok := goja.AssertFunction(value); ok {
		go func() {
			if _, err := call(nil); err != nil {
				log.Errorf("%s", err.Error())
			}
		}()
		return nil
	} else {
		return fmt.Errorf("not a \"function\": %T", value)
	}
}
