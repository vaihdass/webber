package l10n

import (
	"context"
	"net/http"

	"github.com/vaihdass/webber/errors/errh"
)

type Localizer func(errorType, language string) (string, bool)

type LocalizedErrorHandler struct {
	handler     *errh.ErrorHandler
	defaultLang string
	localizerFn Localizer
}

func NewLocalizedErrorHandler(
	handler *errh.ErrorHandler,
	defaultLanguage string,
	localizerFn Localizer,
) *LocalizedErrorHandler {
	return &LocalizedErrorHandler{
		handler:     handler,
		defaultLang: defaultLanguage,
		localizerFn: localizerFn,
	}
}

func (h *LocalizedErrorHandler) Handle(ctx context.Context, operation string, err error, options ...errh.Option) error {
	if err == nil {
		return nil
	}

	err = h.rewrapMessage(ctx, err)

	return h.handler.Handle(ctx, operation, err, options...)
}

func (h *LocalizedErrorHandler) HandleHTTP(
	ctx context.Context, w http.ResponseWriter, r *http.Request,
	operation string, err error, options ...errh.Option,
) {
	if err == nil {
		return
	}

	err = h.rewrapMessage(ctx, err)

	h.handler.HandleHTTP(ctx, w, r, operation, err, options...)
}
