package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

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

type EmbeddingResult struct {
	Chunk     string    `json:"chunk"`
	Embedding []float32 `json:"embedding"`
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

type ChunkIterator struct {
	chunk_size    int
	chunk_overlap int
	data          []byte
}

func (c ChunkIterator) Chunks(yield func([]byte) bool) {
	length := len(c.data)

	next_offset := 0
	for offset := 0; offset < length; offset = next_offset {
		avail := length - offset
		csize := 0
		if c.chunk_size >= avail {
			csize = avail
			next_offset = offset + csize
		} else {
			csize = c.chunk_size
			next_offset = offset + csize - c.chunk_overlap
		}

		if !yield(c.data[offset : offset+csize]) {
			return
		}
	}
}

//export chunk
func chunk() int32 {

	var err error
	chunk_size := 1024
	chunk_size_str, ok := pdk.GetConfig("chunk_size")
	if ok {
		chunk_size, err = strconv.Atoi(chunk_size_str)
		if err != nil {
			pdk.SetError(err)
			return 1
		}
	}

	chunk_overlap := 20
	chunk_overlap_str, ok := pdk.GetConfig("chunk_overlap")
	if ok {
		chunk_overlap, err = strconv.Atoi(chunk_overlap_str)
		if err != nil {
			pdk.SetError(err)
			return 2
		}
	}

	if chunk_overlap >= chunk_size {
		pdk.SetError(errors.New("chunk_size must be larger than chunk_overlap"))
		return 2
	}

	data := pdk.Input()
	iter := ChunkIterator{chunk_size, chunk_overlap, data}

	var result []string
	for c := range iter.Chunks {
		result = append(result, string(c))
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		pdk.SetError(err)
		return 3
	}
	pdk.Output(bytes)
	return 0
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

	chunk_size := 1024
	chunk_size_str, ok := pdk.GetConfig("chunk_size")
	if ok {
		chunk_size, err = strconv.Atoi(chunk_size_str)
		if err != nil {
			pdk.SetError(err)
			return 2
		}
	}

	chunk_overlap := 20
	chunk_overlap_str, ok := pdk.GetConfig("chunk_overlap")
	if ok {
		chunk_overlap, err = strconv.Atoi(chunk_overlap_str)
		if err != nil {
			pdk.SetError(err)
			return 3
		}
	}
	if chunk_overlap >= chunk_size {
		pdk.SetError(errors.New("chunk_size must be larger than chunk_overlap"))
		return 2
	}

	model, ok := pdk.GetConfig("model")
	if !ok {
		pdk.SetError(errors.New("model not found in config"))
		return 4
	}

	data := pdk.Input()

	iter := ChunkIterator{chunk_size, chunk_overlap, data}
	var input []string
	for t := range iter.Chunks {
		input = append(input, string(t))
	}

	embedreq := EmbeddingRequest{Model: model, Input: input}
	payload, err := json.Marshal(embedreq)
	if err != nil {
		pdk.SetError(err)
		return 5
	}

	// create an HTTP Request (without relying on WASI), set headers as needed
	req := pdk.NewHTTPRequest(pdk.MethodPost, u.String())
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(payload)
	// send the request, get response back (can check status on response via res.Status())
	res := req.Send()
	if res.Status() != 200 {
		pdk.SetError(errors.New(fmt.Sprintf("HTTP NOT OK (%d)", res.Status())))
		return 6
	}

	var e EmbeddingResponse
	err = json.Unmarshal(res.Body(), &e)
	if err != nil {
		pdk.SetError(err)
		return 7
	}

	if len(input) != len(e.Embeddings) {
		pdk.SetError(errors.New("embed: generated embedding size not match (#input != #embedding)"))
		return 8
	}

	result := make([]EmbeddingResult, len(input))
	for i, chunk_text := range input {
		result[i] = EmbeddingResult{chunk_text, e.Embeddings[i]}
	}

	bytes, err := json.Marshal(result)
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

	model, ok := pdk.GetConfig("model")
	if !ok {
		pdk.SetError(errors.New("model not found in config"))
		return 2
	}

	prompt := pdk.InputString()
	genreq := GenerateRequest{Model: model, Prompt: prompt, Stream: false}

	payload, err := json.Marshal(genreq)
	if err != nil {
		pdk.SetError(err)
		return 3
	}

	// create an HTTP Request (without relying on WASI), set headers as needed
	req := pdk.NewHTTPRequest(pdk.MethodPost, u.String())
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(payload)
	// send the request, get response back (can check status on response via res.Status())
	res := req.Send()
	if res.Status() != 200 {
		pdk.SetError(errors.New(fmt.Sprintf("HTTP NOT OK (%d)", res.Status())))
		return int32(res.Status())
	}

	var g GenerateResponse
	err = json.Unmarshal(res.Body(), &g)
	if err != nil {
		pdk.SetError(err)
		return 4
	}

	output := []string{g.Response}
	bytes, err := json.Marshal(output)
	if err != nil {
		pdk.SetError(err)
		return 5
	}

	pdk.Output(bytes)

	return 0

}

func main() {}
