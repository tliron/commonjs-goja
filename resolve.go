package commonjs

import (
	contextpkg "context"
	"path/filepath"

	"github.com/tliron/exturl"
)

type ResolveFunc func(context contextpkg.Context, id string, raw bool) (exturl.URL, error)

type CreateResolverFunc func(url exturl.URL, context *Context) ResolveFunc

func NewDefaultResolverCreator(urlContext *exturl.Context, basePaths []exturl.URL, defaultExtension string, allowFilePaths bool) CreateResolverFunc {
	// Copy to protect against later changes
	basePaths = append(basePaths[:0:0], basePaths...)

	// CreateResolverFunc signature
	return func(url exturl.URL, jsContext *Context) ResolveFunc {
		if url != nil {
			basePaths = append([]exturl.URL{url.Base()}, basePaths...)
		}

		if defaultExtension == "" {
			if allowFilePaths {
				return func(context contextpkg.Context, id string, raw bool) (exturl.URL, error) {
					return urlContext.NewValidAnyOrFileURL(context, id, basePaths)
				}
			} else {
				return func(context contextpkg.Context, id string, raw bool) (exturl.URL, error) {
					return urlContext.NewValidURL(context, id, basePaths)
				}
			}
		} else {
			defaultExtension_ := "." + defaultExtension // new var for capture

			if allowFilePaths {
				return func(context contextpkg.Context, id string, raw bool) (exturl.URL, error) {
					if !raw {
						if filepath.Ext(id) == "" {
							id += defaultExtension_
						}
					}

					return urlContext.NewValidAnyOrFileURL(context, id, basePaths)
				}
			} else {
				return func(context contextpkg.Context, id string, raw bool) (exturl.URL, error) {
					if !raw {
						if filepath.Ext(id) == "" {
							id += defaultExtension_
						}
					}

					return urlContext.NewValidURL(context, id, basePaths)
				}
			}
		}
	}
}
