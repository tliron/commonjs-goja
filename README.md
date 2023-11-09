CommonJS for Goja
=================

Use the [Goja](https://github.com/dop251/goja) JavaScript engine in a
[CommonJS-style](https://wiki.commonjs.org/wiki/CommonJS) modular environment.
Supports `require`, `exports`, and more, following the
[Node.js](https://nodejs.org/api/modules.html) implementation as much as possible.

Features:

* Customize the environment with your own special APIs. Useful optional APIs are provided.
* Automatically converts Go field names to dromedary case for a more idiomatic JavaScript experience,
  e.g. `.DoTheThing` becomes `.doTheThing`.
* By default `require` supports full URLs and can resolve paths relative to the current module's
  location. But this can be customized to support your own special resolution and loading method
  (e.g. loading modules from a database).
* Optional support for watching changes to all resolved JavaScript files if they are in the local
  filesystem, allowing your application to restart or otherwise respond to live code updates.
* Optional support for `bind`, which is similar to `require` but exports the JavaScript objects,
  including functions, into a new `goja.Runtime`. This is useful for multi-threaded Go environments
  because a single `goja.Runtime` cannot be used simulatenously by more than one thread. Two variations
  exist: early binding, which creates the `Runtime` when `bind` is called (lower concurrency, higher
  performance), and late binding, which creates the `Runtime` every time the bound object is unbound
  (higher concurrency, lower performance).

Example
-------

`start.go`:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/tliron/commonjs-goja"
    "github.com/tliron/commonjs-goja/api"
    "github.com/tliron/exturl"
    "github.com/dop251/goja"
)

func main() {
    urlContext := exturl.NewContext()
    defer urlContext.Release()

    wd, _ := urlContext.NewWorkingDirFileURL()

    environment := commonjs.NewEnvironment(urlContext, wd)
    defer environment.Release()

    // Support a "console" API (console.log, console.trace, etc.)
	environment.Extensions = append(environment.Extensions, commonjs.Extension{
		Name:   "console",
		Create: api.CreateConsoleExtension(environment.Log),
	})

    // Start!
    environment.Require("./start")
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
