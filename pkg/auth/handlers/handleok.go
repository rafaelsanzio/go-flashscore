package handlers

import (
	"net/http"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type AuthOkPayload struct {
	Health int    `json:"health,omitempty"`
	Test   string `json:"test,omitempty"`
}

func HandleAuthOK(w http.ResponseWriter, r *http.Request) {
	dataReturn := AuthOkPayload{
		Health: 1,
		Test:   "Everthing is OK",
	}

	data, err_ := jsonMarshal(dataReturn)
	if err_ != nil {
		_ = errs.ErrMarshalingJson.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	_, err_ = write(w, data)

	if err_ != nil {
		_ = errs.ErrResponseWriter.Throwf(applog.Log, errs.ErrFmt, err_)
		errs.HttpInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
