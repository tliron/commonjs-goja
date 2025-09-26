package api

import (
	"fmt"
	"html"
	urlpkg "net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonlog"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-kutil/util"
)

// ([commonjs.CreateExtensionFunc] signature)
func CreateUtilExtension(jsContext *commonjs.Context) any {
	return NewUtil(jsContext.Environment.Log)
}

//
// Util
//

type Util struct {
	log commonlog.Logger
}

func NewUtil(log commonlog.Logger) *Util {
	return &Util{
		log: log,
	}
}

func (self *Util) StringToBytes(string_ string) []byte {
	return util.StringToBytes(string_)
}

// Another way to achieve this in JavaScript: String.fromCharCode.apply(null, bytes)
func (self *Util) BytesToString(bytes []byte) string {
	return util.BytesToString(bytes)
}

// Encode bytes as base64
func (self *Util) Btoa(bytes []byte) string {
	return util.ToBase64(bytes)
}

// Decode base64 to bytes
func (self Transcribe) Atob(b64 string) ([]byte, error) {
	return util.FromBase64(b64)
}

func (self *Util) DeepCopy(value ard.Value) ard.Value {
	return ard.Copy(value)
}

func (self *Util) DeepEquals(a ard.Value, b ard.Value) bool {
	return ard.Equals(a, b)
}

func (self *Util) IsType(value ard.Value, typeName string) (bool, error) {
	// Special case whereby an integer stored as a float type has been optimized to an integer type
	if (typeName == "!!float") && util.IsInteger(value) {
		return true, nil
	}

	if validate, ok := ard.TypeValidators[ard.TypeName(typeName)]; ok {
		return validate(value), nil
	} else {
		return false, fmt.Errorf("unsupported type name: %q", typeName)
	}
}

// Construct a URL via its parts. Supports parts:
//
//   - "scheme": as string
//   - "username": as string
//   - "password": as string
//   - "host": as string
//   - "port": as unsigned integer
//   - "path": as string
//   - "query": map of strings to either a list of multiple values or
//     or a single value as string
//   - "fragment": as string
func (self *Util) Url(config ard.StringMap) (string, error) {
	var invalidKeys []string
	for key := range config {
		switch key {
		case "scheme", "username", "password", "host", "port", "path", "query", "fragment":
		default:
			invalidKeys = append(invalidKeys, key)
		}
	}
	if len(invalidKeys) > 0 {
		return "", fmt.Errorf("invalid keys for \"url\": %s", util.JoinQuote(invalidKeys, ", "))
	}

	config_ := ard.With(config).ConvertSimilar().NilMeansZero()

	var url urlpkg.URL

	url.Scheme, _ = config_.Get("scheme").String()

	if username, ok := config_.Get("username").String(); ok {
		if password, ok := config_.Get("password").String(); ok {
			url.User = urlpkg.UserPassword(username, password)
		} else {
			url.User = urlpkg.User(username)
		}
	}

	url.Host, _ = config_.Get("host").String()

	if port, ok := config_.Get("port").UnsignedInteger(); ok {
		url.Host += ":" + strconv.FormatUint(port, 10)
	}

	url.Path, _ = config_.Get("path").String()
	if (url.Path != "") && (!strings.HasPrefix(url.Path, "/")) {
		url.Path = "/" + url.Path
	}

	if query, ok := config_.Get("query").StringMap(); ok {
		values := make(urlpkg.Values)
		for key, value := range query {
			value_ := ard.With(value).ConvertSimilar().NilMeansZero()
			if listValue, ok := value_.StringList(); ok {
				values[key] = listValue
			} else if singleValue, ok := value_.String(); ok {
				values[key] = []string{singleValue}
			}
		}
		url.RawQuery = values.Encode()
	}

	url.Fragment, _ = config_.Get("fragment").String()

	return url.String(), nil
}

func (self *Util) EscapeHtml(text string) string {
	return html.EscapeString(text)
}

func (self *Util) UnescapeHtml(text string) string {
	return html.UnescapeString(text)
}

func (self *Util) Hash(value ard.Value) (uint64, error) {
	return hashstructure.Hash(value, hashstructure.FormatV2, nil)
}

func (self *Util) Sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func (self *Util) Fail(message string) {
	util.Fail(message)
}

func (self *Util) Failf(format string, args ...any) {
	util.Failf(format, args...)
}

func (self *Util) TimeFromUnix(seconds int64, nanoseconds int64) time.Time {
	return time.Unix(seconds, nanoseconds)
}

func (self *Util) FormatTime(t time.Time, format string) string {
	if format == "" {
		format = time.RFC3339Nano
	}
	return t.Format(format)
}

func (self *Util) Now() time.Time {
	return time.Now()
}

func (self *Util) Mutex() util.RWLocker {
	return util.NewDefaultRWLocker()
}

var onces sync.Map

func (self *Util) Once(name string, value goja.Value, this goja.Value, arguments []goja.Value) error {
	if call, ok := goja.AssertFunction(value); ok {
		var err error

		call_ := func() error {
			_, err := call(this, arguments...)
			return err
		}

		call__ := func() {
			err = call_()
		}

		once, _ := onces.LoadOrStore(name, new(sync.Once))
		once.(*sync.Once).Do(call__)

		return err
	} else {
		return fmt.Errorf("not a function: %T", value)
	}
}

// Goroutine
func (self *Util) Go(value goja.Value, this goja.Value, arguments []goja.Value) error {
	if call, ok := goja.AssertFunction(value); ok {
		call_ := func() error {
			_, err := call(this, arguments...)
			return err
		}

		go commonlog.CallAndLogError(call_, "Util.Go", self.log)

		return nil
	} else {
		return fmt.Errorf("not a \"function\": %T", value)
	}
}
