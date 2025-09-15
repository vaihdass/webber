package l10n

import (
	"context"
	"net/http"
	"strings"
)

const HeaderKey = "X-L10n-Header"

func ExtractLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		lang := r.Header.Get(HeaderKey)
		if lang == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx = context.WithValue(ctx, langKey{}, strings.ToLower(lang))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
