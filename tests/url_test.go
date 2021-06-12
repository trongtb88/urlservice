package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/trongtb88/urlservice/api/constants"
	"github.com/trongtb88/urlservice/api/entity"
	"github.com/trongtb88/urlservice/api/utils"
	"gopkg.in/go-playground/assert.v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidShortCode(t *testing.T) {

	shortCode := "a@"
	params := entity.HttpRequestUrl{
		Url:       "https://github.com/trongtb88",
		Shortcode: shortCode,
	}
	body, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", "/shortcode", bytes.NewReader(body))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Shorten)
	handler.ServeHTTP(rr, req)

	var response entity.HttpResponseUrl

	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Printf("Cannot convert to json: %v", err)
	}
	assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
}

func TestCreateSuccessShortCode(t *testing.T) {

	shortCode := utils.GenerateRandom(constants.SHORT_CODE_REGEX)
	params := entity.HttpRequestUrl{
		Url:       "https://github.com/trongtb88",
		Shortcode: shortCode,
	}
	body, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", "/shortcode", bytes.NewReader(body))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Shorten)
	handler.ServeHTTP(rr, req)

	var response entity.HttpResponseUrl

	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Printf("Cannot convert to json: %v", err)
	}
	assert.Equal(t, rr.Code, http.StatusCreated)
	assert.Equal(t, response.Shortcode, shortCode)
}

func TestDuplicateShortCode(t *testing.T) {

	shortCode := utils.GenerateRandom(constants.SHORT_CODE_REGEX)
	params := entity.HttpRequestUrl{
		Url:       "https://github.com/trongtb88",
		Shortcode: shortCode,
	}
	body, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", "/shortcode", bytes.NewReader(body))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Shorten)
	handler.ServeHTTP(rr, req)

	var response entity.HttpResponseUrl

	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Printf("Cannot convert to json: %v", err)
	}

	req, err = http.NewRequest("POST", "/shortcode", bytes.NewReader(body))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(server.Shorten)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		fmt.Printf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusConflict)
}

func TestRedirectShortCode(t *testing.T) {

	shortCode := utils.GenerateRandom(constants.SHORT_CODE_REGEX)
	params := entity.HttpRequestUrl{
		Url:       "https://github.com/trongtb88",
		Shortcode: shortCode,
	}

	body, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", "/shortcode", bytes.NewReader(body))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Shorten)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusCreated)

	url := "/{shortcode}"

	replacer := strings.NewReplacer("{shortcode}", shortCode)
	url = replacer.Replace(url)

	req2, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	req2 = mux.SetURLVars(req2, map[string]string{"shortcode": shortCode})
	rr2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(server.RedirectURL)
	handler2.ServeHTTP(rr2, req2)

	assert.Equal(t, rr2.Code, http.StatusFound)
}

func TestStatsShortCode(t *testing.T) {

	shortCode := utils.GenerateRandom(constants.SHORT_CODE_REGEX)
	params := entity.HttpRequestUrl{
		Url:       "https://github.com/trongtb88",
		Shortcode: shortCode,
	}

	body, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", "/shortcode", bytes.NewReader(body))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Shorten)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusCreated)

	url := "/{shortcode}"
	replacer := strings.NewReplacer("{shortcode}", shortCode)
	url = replacer.Replace(url)

	// Redirect 10 times
	for i := 0; i < 10; i++ {
		req2, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req2 = mux.SetURLVars(req2, map[string]string{"shortcode": shortCode})
		rr2 := httptest.NewRecorder()
		handler2 := http.HandlerFunc(server.RedirectURL)
		handler2.ServeHTTP(rr2, req2)
	}
	// end redirect


	// get stats
	url = "/{shortcode}"
	replacer = strings.NewReplacer("{shortcode}", shortCode)
	url = replacer.Replace(url)

	req3, err := http.NewRequest("GET", url +"/stats", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	req3 = mux.SetURLVars(req3, map[string]string{"shortcode": shortCode})

	rr3 := httptest.NewRecorder()

	handler3 := http.HandlerFunc(server.Stats)
	handler3.ServeHTTP(rr3, req3)

	response := entity.HttpResponseStatsUrl{}

	err = json.Unmarshal([]byte(rr3.Body.String()), &response)
	if err != nil {
		fmt.Printf("Cannot convert to json: %v", err)
	}

	assert.Equal(t, rr3.Code, http.StatusOK)
	assert.Equal(t, response.RedirectCount, int64(10))
}
