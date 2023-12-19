package commonjs

import (
	contextpkg "context"
	"path/filepath"

	"github.com/tliron/exturl"
)

type ResolveFunc func(context contextpkg.Context, id string, bareId bool) (exturl.URL, error)

type CreateResolverFunc func(fromUrl exturl.URL, jsContext *Context) ResolveFunc

func NewDefaultResolverCreator(defaultExtension string, allowFilePaths bool, urlContext *exturl.Context, basePaths ...exturl.URL) CreateResolverFunc {
	// CreateResolverFunc signature
	return func(fromUrl exturl.URL, jsContext *Context) ResolveFunc {
		basePaths_ := basePaths // new var for capture
		if fromUrl != nil {
			basePaths_ = append([]exturl.URL{fromUrl.Base()}, basePaths_...)
		}

		if defaultExtension == "" {
			if allowFilePaths {
				// ResolveFunc signature
				return func(context contextpkg.Context, id string, bareId bool) (exturl.URL, error) {
					return urlContext.NewValidAnyOrFileURL(context, id, basePaths_)
				}
			} else {
				// ResolveFunc signature
				return func(context contextpkg.Context, id string, bareId bool) (exturl.URL, error) {
					return urlContext.NewValidURL(context, id, basePaths_)
				}
			}
		} else {
			defaultExtension_ := "." + defaultExtension // new var for capture

			if allowFilePaths {
				// ResolveFunc signature
				return func(context contextpkg.Context, id string, bareId bool) (exturl.URL, error) {
					if !bareId {
						if filepath.Ext(id) == "" {
							id += defaultExtension_
						}
					}

					return urlContext.NewValidAnyOrFileURL(context, id, basePaths_)
				}
			} else {
				// ResolveFunc signature
				return func(context contextpkg.Context, id string, bareId bool) (exturl.URL, error) {
					if !bareId {
						if filepath.Ext(id) == "" {
							id += defaultExtension_
						}
					}

					return urlContext.NewValidURL(context, id, basePaths_)
				}
			}
		}
	}
}
