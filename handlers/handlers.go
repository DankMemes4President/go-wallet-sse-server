package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"go-wallet-sse-server/config"
	"go-wallet-sse-server/errors"
	"go-wallet-sse-server/internal/response"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleSSE(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set response headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Add the channel to the PubSub pubSub
		// We get the userID from the request, we then subscribe him to his specific Redis channel

		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				token := headerParts[1]
				claims, err := jwt.HMACCheck([]byte(token), []byte(app.Config.Jwt.SecretKey))
				if err != nil {
					errors.InvalidAuthenticationToken(w, r, app)
					return
				}

				if !claims.Valid(time.Now()) {
					errors.InvalidAuthenticationToken(w, r, app)
					return
				}

				if claims.Issuer != app.Config.BaseURL {
					errors.InvalidAuthenticationToken(w, r, app)
					return
				}

				if !claims.AcceptAudience(app.Config.BaseURL) {
					errors.InvalidAuthenticationToken(w, r, app)
					return
				}

				userID, err := strconv.Atoi(claims.Subject)
				if err != nil {
					errors.ServerError(w, r, err, app)
					return
				}
				userChannelName := fmt.Sprintf("user#%d", userID)
				err = app.PubSub.Subscribe(app.Rdb.Context(), userChannelName)
				if err != nil {
					app.Logger.Debugf("Error %s", userChannelName)
					app.Logger.Debugf(err.Error())
					//errors.ServerError(w, r, err, app) // This means even if the channel of that user does not exist, it'll throw a 500. This is not optimal
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer app.PubSub.Close()
				statusCheckStruct := response.StatusCheckMessage{
					Type: "status_check",
					Payload: response.StatusCheckPayload{
						"ok",
					},
				}
				statusCheckMsg, _ := json.Marshal(&statusCheckStruct)

				fmt.Fprintf(w, "data: %s\n\n", strings.TrimSuffix(string(statusCheckMsg), "\n\n"))
				w.(http.Flusher).Flush()

				channel := app.PubSub.Channel()

				for msg := range channel {
					app.Logger.Infof(msg.Channel, msg.Payload)
					fmt.Fprintf(w, "data: %s\n\n", strings.TrimSuffix(msg.Payload, "\n\n"))
					w.(http.Flusher).Flush()
				}
			}
		}
	}
}
