package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	extism "github.com/extism/go-sdk"
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

func main() {

	cfg := make(map[string]string)
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: "ollama.wasm",
			},
		},
		AllowedHosts: []string{"localhost"},
		Config:       cfg,
	}

	ctx := context.Background()
	config := extism.PluginConfig{
		EnableWasi: true,
	}
	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})
	if err != nil {
		fmt.Printf("Failed to initialize plugin: %v\n", err)
		os.Exit(1)
	}

	payload := EmbeddingRequest{
		Model: "llama3.2",
		Input: []string{"how are you?"},
	}
	reqbody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	exit, out, err := plugin.Call("embed", reqbody)
	if err != nil {
		fmt.Printf("plugin call %v\n", err)
		os.Exit(int(exit))
	}

	fmt.Println(string(out))

	greq := GenerateRequest{
		Model: "llama3.2",
		Prompt: `Question: Who is the queen of england. Please summarize the answer below with the question.
                Here is the answer
                1. Charles is the king now
                2. Elizaberth is the queen last year
                3. Bloody mary is the greatest queen
                4. Edward is the strongest king ever`,
		Stream: false,
	}

	reqbody, err = json.Marshal(greq)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exit, out, err = plugin.Call("generate", reqbody)
	if err != nil {
		fmt.Printf("plugin call %v\n", err)
		os.Exit(int(exit))
	}
	fmt.Println(string(out))
}
