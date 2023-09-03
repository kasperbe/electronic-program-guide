package router

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kasperbe/electronic-program-guide/epg"
)

func unmarshal(req *http.Request) (epg.TranslateInput, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return epg.TranslateInput{}, err
	}
	defer req.Body.Close()

	var body epg.TranslateInput
	err = json.Unmarshal(data, &body)
	if err != nil {
		return epg.TranslateInput{}, err
	}

	return body, nil
}

func TranslateEPG(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		return
	}

	input, err := unmarshal(req)
	if err != nil {
		BadRequest(w, err)
		return
	}

	program := epg.New(input)
	PlainText(w, program.ToString())
}

// Helpers for responses.
// Usually server libraries include these, but we just roll these ourselves for simplicity this time.
func PlainText(w http.ResponseWriter, data string) {
	w.Header().Add("content-type", "text/plain")
	w.WriteHeader(200)
	io.WriteString(w, data)
}

func BadRequest(w http.ResponseWriter, err error) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(400)
	io.WriteString(w, fmt.Sprintf(`{ "error": "%s" }`, err))
}

func InternalServerError(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(500)
	io.WriteString(w, `{ "error": "internal server error" }`)
}
