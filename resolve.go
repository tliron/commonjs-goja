package commonjs

import (
	contextpkg "context"
	"path/filepath"

	"github.com/tliron/exturl"
)

type ResolveFunc func(context contextpkg.Context, id string, raw bool) (exturl.URL, error)

type CreateResolverFunc func(url exturl.URL, context *Context) ResolveFunc

func NewDefaultResolverCreator(urlContext *exturl.Context, path []exturl.URL, defaultExtension string) CreateResolverFunc {
	return func(url exturl.URL, jsContext *Context) ResolveFunc {
		var bases []exturl.URL

		if url != nil {
			bases = append([]exturl.URL{url.Base()}, path...)
		} else {
			bases = path
		}

		if defaultExtension == "" {
			// ResolveFunc signature
			return func(context contextpkg.Context, id string, raw bool) (exturl.URL, error) {
				return urlContext.NewValidAnyOrFileURL(context, id, bases)
			}
		} else {
			defaultExtension_ := "." + defaultExtension

			// ResolveFunc signature
			return func(context contextpkg.Context, id string, raw bool) (exturl.URL, error) {
				if !raw {
					if filepath.Ext(id) == "" {
						id += defaultExtension_
					}
				}

				return urlContext.NewValidAnyOrFileURL(context, id, bases)
			}
		}
	}
}
