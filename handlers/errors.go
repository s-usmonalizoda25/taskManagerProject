package handlers

import (
	"errors"
	"net/http"

	"github.com/s-usmonalizoda25/taskManagerProject/pkg/errs"
	"github.com/s-usmonalizoda25/taskManagerProject/pkg/logger"
	"go.uber.org/zap"
)

func HandleError(w http.ResponseWriter, log *logger.Logger, err error) {
	switch {

	case errors.Is(err, errs.ErrEmptyTitle),
		errors.Is(err, errs.ErrInvalidStatus),
		errors.Is(err, errs.ErrInvalidID),
		errors.Is(err, errs.ErrEmptyUsername),
		errors.Is(err, errs.ErrEmptyEmail),
		errors.Is(err, errs.ErrUsernameTaken),
		errors.Is(err, errs.ErrEmailTaken):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, errs.ErrTaskNotFound),
		errors.Is(err, errs.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	default:
		log.Error("Internal server error occurred", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
