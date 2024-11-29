package middleware

import (
	"deuce109/7dtd-map-server/v2/logging"
	"deuce109/7dtd-map-server/v2/models"
	"deuce109/7dtd-map-server/v2/repositories"
	"deuce109/7dtd-map-server/v2/utils/crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWTMiddleware struct {
	Secret  string
	Filters []func(interface{}) bool
}

const JWT_ERROR_MESSAGE = "Error when verifying JWT token: %s"

func (m *JWTMiddleware) decodeJwt(data string) (*models.UserInfo, int) {
	chunks := strings.Split(data, ".")

	if len(chunks) != 3 {
		return nil, http.StatusInternalServerError
	}

	msg := fmt.Sprintf("%s.%s", chunks[0], chunks[1])

	genSig := crypto.HMACSHA512(msg, m.Secret)

	sig := chunks[2]

	if crypto.HMACEquals(genSig, sig) {
		var header map[string]string

		headerString, err := base64.URLEncoding.DecodeString(chunks[0])

		if err != nil {
			logging.Error(fmt.Sprintf(JWT_ERROR_MESSAGE, err))

			return nil, 500
		}

		err = json.Unmarshal([]byte(headerString), &header)
		if err != nil {
			logging.Error(fmt.Sprintf(JWT_ERROR_MESSAGE, err))

			return nil, 500
		}

		var payload *models.UserInfo

		payloadString, err := base64.URLEncoding.DecodeString(chunks[1])

		if err != nil {
			logging.Error(fmt.Sprintf(JWT_ERROR_MESSAGE, err))

			return nil, 500
		}

		err = json.Unmarshal([]byte(payloadString), &payload)
		if err != nil {
			logging.Error(fmt.Sprintf(JWT_ERROR_MESSAGE, err))

			return nil, 500
		}
		return payload, 204
	} else {
		logging.Error("Error when verifying JWT token: invalid signature")

		return nil, 403
	}

}

func (m *JWTMiddleware) CheckJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		payload, statusCode := m.decodeJwt(strings.TrimPrefix(authHeader, "Token "))

		if statusCode != 204 {
			w.WriteHeader(statusCode)
			return
		}

		authorized := true
		for _, filter := range m.Filters {
			authorized = authorized && filter(payload)
		}

		if !authorized {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		gameId := payload.PlatformId

		r.Header["Game-Id"] = []string{gameId}

		next.ServeHTTP(w, r)

	})
}

func (m *JWTMiddleware) CreateJwt(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")

	player := repositories.GetPlayerInfo(username)

	if player == nil {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return
	}

	header := &jwtHeader{
		Alg: "HS256",
		Typ: "jwt",
	}

	jsonHeader, err := models.ToJson(header)

	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling header value %s", err))
		w.WriteHeader(500)
		return
	}

	jsonPayload, err := models.ToJson(player)

	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling payload value %s", err))
		w.WriteHeader(500)
		return
	}

	headerString := base64.URLEncoding.EncodeToString(jsonHeader)
	payloadString := base64.URLEncoding.EncodeToString(jsonPayload)

	sig := crypto.HMACSHA512(
		fmt.Sprintf(
			"%s.%s",
			headerString,
			payloadString,
		),
		m.Secret,
	)

	w.Write([]byte(fmt.Sprintf("%s.%s.%s", headerString, payloadString, sig)))

}
