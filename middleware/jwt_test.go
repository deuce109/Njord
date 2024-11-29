package middleware

import (
	"bytes"
	"crypto/rand"
	"deuce109/7dtd-map-server/v2/readers"
	"deuce109/7dtd-map-server/v2/utils/crypto"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

var mw *JWTMiddleware

var jwt string

// const TEST_AUTH_HEADER = "Token eyJhbGciOiJIUzI1NiIsInR5cCI6Imp3dCJ9.eyJ1aWQiOiIxMjM0NTYiLCJwbHRpZCI6IjMyMTU0OTg0MTA2IiwidXNyIjoidGVzdFVzZXJOYW1lIiwibGV2ZWwiOjAsInBsdCI6IkFueSJ9.fNHKZNtkztaO9ZxVowxB3dwt093URSXndPqE_Aij72s="

func TestMain(m *testing.M) {

	rng := &crypto.Random{
		Reader: rand.Reader,
	}

	secret, _ := rng.GetRandomSecret(32)

	alwaysTrue := func(interface{}) bool { return true }

	mw = &JWTMiddleware{
		Secret:  secret,
		Filters: [](func(interface{}) bool){alwaysTrue},
	}

	readers.GetUserInfo("../test_data/good", "../test_data/good")

	m.Run()

}

func TestCreateJwt(t *testing.T) {

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	writer.WriteField("username", "testUserName")
	writer.Close()

	w := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/login", body)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	mw.CreateJwt(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("CreateJwt returned a non-200 status code: %d", w.Code)
		return
	}

	resBody, _ := io.ReadAll(res.Body)

	jwt = string(resBody)

	match, _ := regexp.Match(`[A-Za-z0-9+\/]+={0,3}\.[A-Za-z0-9+\/]+={0,3}\.[A-Za-z0-9+\/]+={0,3}`, resBody)

	if !match {
		t.Errorf("CreateJwt returned and invalid JWT")
	}
}

func TestCreateJwtReturns404(t *testing.T) {

	body := new(bytes.Buffer)

	w := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/login", body)

	mw.CreateJwt(w, req)

	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("CreateJwt returned a non-404 status code: %d", w.Code)
	}
}

func TestCheckJwt(t *testing.T) {

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("Game-Id")
		if val == "" {
			t.Error("Game-Id not present")
		}
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	writer.Close()

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", body)

	req.Header.Set("Authorization", "Token "+jwt)

	mw.CheckJwt(nextHandler).ServeHTTP(w, req)
}

func TestCheckJwtReturns500WithPoorlyFormatedAuthHeader(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	writer.Close()

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", body)

	req.Header.Set("Authorization", "toekn")

	mw.CheckJwt(nextHandler).ServeHTTP(w, req)

	if w.Result().StatusCode != 500 {
		t.Error("Handler should return status code 500")
	}
}

func TestCheckJwtReturns401WithFailedFilters(t *testing.T) {
	alwaysFalse := func(interface{}) bool { return false }

	falseMw := &JWTMiddleware{
		Secret:  mw.Secret,
		Filters: [](func(interface{}) bool){alwaysFalse},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	writer.Close()

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", body)

	req.Header.Set("Authorization", "Token "+jwt)

	falseMw.CheckJwt(nextHandler).ServeHTTP(w, req)

	code := w.Result().StatusCode

	if code != 401 {
		t.Errorf("Handler should return status code 401 but returned %d", code)
	}
}

func TestCheckJwtReturns401WithEmptyAuthHeaders(t *testing.T) {

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	writer.Close()

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", body)

	req.Header.Set("Authorization", "")

	mw.CheckJwt(nextHandler).ServeHTTP(w, req)

	code := w.Result().StatusCode

	if code != 401 {
		t.Errorf("Handler should return status code 401 but received %d", code)
	}
}

func TestDeocdeJwtReturns403WithImproperSignature(t *testing.T) {

	headerString := "tt"
	payloadString := "tt"
	sig := "tt"

	testJwt := fmt.Sprintf("%s.%s.%s", headerString, payloadString, sig)

	_, code := mw.decodeJwt(testJwt)

	if code != 403 {
		t.Errorf("Handler should return status code 403, got %d", code)
	}
}

func TestDeocdeJwtReturns500WithPoorlyFormatedJwtHeaderUrlValue(t *testing.T) {

	headerString := "tt"
	payloadString := "tt"

	sig := crypto.HMACSHA512(
		fmt.Sprintf(
			"%s.%s",
			headerString,
			payloadString,
		),
		mw.Secret,
	)
	testJwt := fmt.Sprintf("%s.%s.%s", headerString, payloadString, sig)

	_, code := mw.decodeJwt(testJwt)

	if code != 500 {
		t.Errorf("Handler should return status code 500, got %d", code)
	}
}

func TestDeocdeJwtReturns500WithPoorlyFormatedJwtPayloadUrlValue(t *testing.T) {

	headerString := strings.Split(jwt, ".")[0]
	payloadString := "tt"

	sig := crypto.HMACSHA512(
		fmt.Sprintf(
			"%s.%s",
			headerString,
			payloadString,
		),
		mw.Secret,
	)
	testJwt := fmt.Sprintf("%s.%s.%s", headerString, payloadString, sig)

	_, code := mw.decodeJwt(testJwt)

	if code != 500 {
		t.Errorf("Handler should return status code 500, got %d", code)
	}
}

func TestDeocdeJwtReturns500WithPoorlyFormatedJwtHeaderJsonValue(t *testing.T) {

	headerString := "eyJ0ZXN0dGVzdCJ9"
	payloadString := "tt"

	sig := crypto.HMACSHA512(
		fmt.Sprintf(
			"%s.%s",
			headerString,
			payloadString,
		),
		mw.Secret,
	)
	testJwt := fmt.Sprintf("%s.%s.%s", headerString, payloadString, sig)

	_, code := mw.decodeJwt(testJwt)

	if code != 500 {
		t.Errorf("Handler should return status code 500, got %d", code)
	}
}

func TestDeocdeJwtReturns500WithPoorlyFormatedJwtPayloadJsonValue(t *testing.T) {

	headerString := strings.Split(jwt, ".")[0]
	payloadString := "eyJ0ZXN0dGVzdCJ9"

	sig := crypto.HMACSHA512(
		fmt.Sprintf(
			"%s.%s",
			headerString,
			payloadString,
		),
		mw.Secret,
	)
	testJwt := fmt.Sprintf("%s.%s.%s", headerString, payloadString, sig)

	_, code := mw.decodeJwt(testJwt)

	if code != 500 {
		t.Errorf("Handler should return status code 500, got %d", code)
	}
}
