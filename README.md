CommonJS for Goja
=================

Use [Goja](https://github.com/dop251/goja) JavaScript engine in a
[CommonJS-style modular](https://wiki.commonjs.org/wiki/CommonJS) environment.
Supports `require`, `exports`, and more, following the
[Node.js](https://nodejs.org/api/modules.html) implementation as much as possible.

Features:

* Customize the environment with your own special extensions.
* By default `require` supports full URLs and can resolve paths relative to the current module's
  location. But this can be customized to support your own special resolution and loading method
  (e.g. loading modules from a database).
* Optional support for watching changes to all resolved JavaScript files (if they are in the local
  filesystem), allowing your application to restart or otherwise respond accordingly to live code
  updates.
* Optional support for `bind`, which is similar to `require` but exports the JavaScript objects,
  including functions, into a new `goja.Runtime`. This is useful for multi-threaded Go environments
  because a single `goja.Runtime` cannot be used simulatenously by more than one thread. Two variations
  exist: early binding, which creates the `Runtime` when `bind` is called (low concurrency, high
  performance), and late binding, which creates the `Runtime` every time the bound object is unbound
  (high concurrency, lower performance).

Example
-------

`start.go`:

```go
package main

import (
    "log"
    "os"

    "github.com/tliron/commonjs-goja"
    "github.com/tliron/exturl"
    "github.com/dop251/goja"
)

func main() {
    urlContext := exturl.NewContext()
    defer urlContext.Release()

    wd, _ := urlContext.NewWorkingDirFileURL()

    environment := commonjs.NewEnvironment(urlContext, []exturl.URL{wd})
    defer environment.Release()

    // Implementation of "console.log"
    environment.Extensions = append(environment.Extensions, commonjs.Extension{
        Name: "console",
        Create: func(context *commonjs.Context) goja.Value {
            return context.Environment.Runtime.ToValue(ConsoleAPI{})
        },
    })

    // Support for "bind" (late binding)
    environment.Extensions = append(environment.Extensions, commonjs.Extension{
        Name:   "bind",
        Create: commonjs.CreateLateBindExtension,
    })

    // Start!
    environment.RequireID("./start")
}

var logger = log.New(os.Stdout, "console: ", log.LstdFlags)

type ConsoleAPI struct{}

func (self ConsoleAPI) Log(message string) {
    logger.Println(message)
}
```

`start.js`:

```js
const hello = require('./lib/hello');

hello.sayit();
```

`lib/hello.js`:

```js
exports.sayit = function() {
    console.log('hi!');
};
```
