package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/extism/go-pdk"
)

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingResponse struct {
	Model           string      `json:"model"`
	Embeddings      [][]float32 `json:"embeddings"`
	TotalDuration   int64       `json:"total_duration"`
	LoadDuration    int64       `json:"load_duration"`
	PromptEvalCount int         `json:"prompt_eval_count"`
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func getApiUrl(apiURI string) (*url.URL, error) {
	address, ok := pdk.GetConfig("address")
	var u string
	var err error
	if ok {
		u, err = url.JoinPath(address, apiURI)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = url.JoinPath("http://localhost:11434", apiURI)
		if err != nil {
			return nil, err
		}
	}

	return url.Parse(u)
}

//export embed
func embed() int32 {

	apiURI := "/api/embed"
	u, err := getApiUrl(apiURI)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	payload := pdk.Input()

	// create an HTTP Request (without relying on WASI), set headers as needed
	req := pdk.NewHTTPRequest(pdk.MethodPost, u.String())
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(payload)
	// send the request, get response back (can check status on response via res.Status())
	res := req.Send()
	if res.Status() != 200 {
		pdk.SetError(errors.New(fmt.Sprintf("HTTP NOT OK (%d)", res.Status())))
		return 1
	}

	var e EmbeddingResponse
	err = json.Unmarshal(res.Body(), &e)
	if err != nil {
		pdk.SetError(err)
		return 2
	}

	bytes, err := json.Marshal(e.Embeddings)
	if err != nil {
		pdk.SetError(err)
		return 3
	}

	pdk.Output(bytes)

	return 0
}

//export generate
func generate() int32 {
	apiURI := "/api/generate"
	u, err := getApiUrl(apiURI)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	payload := pdk.Input()

	// create an HTTP Request (without relying on WASI), set headers as needed
	req := pdk.NewHTTPRequest(pdk.MethodPost, u.String())
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(payload)
	// send the request, get response back (can check status on response via res.Status())
	res := req.Send()
	if res.Status() != 200 {
		pdk.SetError(errors.New(fmt.Sprintf("HTTP NOT OK (%d)", res.Status())))
		return 1
	}

	var g GenerateResponse
	err = json.Unmarshal(res.Body(), &g)
	if err != nil {
		pdk.SetError(err)
		return 2
	}

	output := []string{g.Response}
	bytes, err := json.Marshal(output)
	if err != nil {
		pdk.SetError(err)
		return 3
	}

	pdk.Output(bytes)

	return 0

}

func main() {}
