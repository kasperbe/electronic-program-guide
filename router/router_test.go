package router

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateEndpoint(t *testing.T) {
	tt := map[string]struct {
		input              string
		expectedStatusCode int
	}{
		"Empty string": {input: "", expectedStatusCode: 400},
		"Broken Json":  {input: "{", expectedStatusCode: 400},
		"Json":         {input: "{}", expectedStatusCode: 200},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {

			req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(tc.input)))
			if err != nil {
				t.Logf("new request: %v", err)
				t.FailNow()
			}

			writer := httptest.NewRecorder()

			TranslateEPG(writer, req)

			assert.Equal(t, tc.expectedStatusCode, writer.Result().StatusCode)
		})
	}
}

func TestEmptyPayload(t *testing.T) {
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		t.Logf("new request: %v\n", err)
		t.FailNow()
	}

	writer := httptest.NewRecorder()

	TranslateEPG(writer, req)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Logf("read all: %v\n", err)
	}

	expected := `Monday: Nothing aired today
Tuesday: Nothing aired today
Wednesday: Nothing aired today
Thursday: Nothing aired today
Friday: Nothing aired today
Saturday: Nothing aired today
Sunday: Nothing aired today`

	assert.Equal(t, expected, string(data))
}

func TestWithPayload(t *testing.T) {
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`{ "monday": [ { 
 "title": "Nyhederne", 
 "state": "begin", 
 "time": 21600 
 }, 
 { 
 "title": "Nyhederne", 
 "state": "end", 
 "time": 36000 
 } 
 ] 
}`)))
	if err != nil {
		t.Logf("new request: %v\n", err)
		t.FailNow()
	}

	writer := httptest.NewRecorder()

	TranslateEPG(writer, req)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Logf("read all: %v\n", err)
	}

	expected := `Monday: Nyhederne 7 - 11
Tuesday: Nothing aired today
Wednesday: Nothing aired today
Thursday: Nothing aired today
Friday: Nothing aired today
Saturday: Nothing aired today
Sunday: Nothing aired today`

	assert.Equal(t, expected, string(data))
}
