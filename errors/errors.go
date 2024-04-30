package errors

import (
	"fmt"
	"go-wallet-sse-server/config"
	"go-wallet-sse-server/internal/response"
	"net/http"
	"strings"
)

// I have changed the signatures of lot of these functions to inject dependencies, in cases where this signature needs
// to follow a constraint we'll take help of factory functions to create our functions and inject dependencies simultaneously.

func ServerError(w http.ResponseWriter, r *http.Request, err error, app *config.Application) {
	//trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//app.Logger.Output(5, trace)
	app.Logger.Errorf(err.Error())
	message := "The server encountered a problem and could not process your request"
	ErrorMessage(w, r, http.StatusInternalServerError, message, nil, app)
}

func InvalidAuthenticationToken(w http.ResponseWriter, r *http.Request, app *config.Application) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	ErrorMessage(w, r, http.StatusUnauthorized, "Invalid authentication token", headers, app)
}

func ErrorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header, app *config.Application) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.Logger.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NotFound(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "The requested resource could not be found"
		ErrorMessage(w, r, http.StatusNotFound, message, nil, app)
	}
}
func MethodNotAllowed(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
		ErrorMessage(w, r, http.StatusMethodNotAllowed, message, nil, app)
	}
}
