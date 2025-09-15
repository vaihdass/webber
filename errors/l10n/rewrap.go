package l10n

import (
	"context"

	"github.com/vaihdass/webber/errors/errh"
	"github.com/vaihdass/webber/errors/xerr"
)

type langKey struct{}

func (h *LocalizedErrorHandler) rewrapMessage(ctx context.Context, err error) error {
	xErr, ok := xerr.From(err)
	if !ok {
		return err
	}

	newMsg, localized := h.localizerFn(xErr.Type(), h.extractLanguage(ctx))
	if !localized {
		return err
	}

	return errh.TryRewrapTypedErr(err, newMsg)
}

func (h *LocalizedErrorHandler) extractLanguage(ctx context.Context) string {
	lang, ok := ctx.Value(langKey{}).(string)
	if !ok {
		return h.defaultLang
	}

	return lang
}
