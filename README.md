# Ollama Client in Web Assembly

We have implemented three export functions such as 'chunk', 'embed' and 'generate'.

1. 'chunk' function to cut the text into multiple chunks
```
Configuration parameters:
1. chunk_size = size of the chunk, default 1024
2. chunk_overlap, number of byte overlapping, default 20

Input:
string of text input

Output:
JSON Array of chunks
```

2. 'embed' function to convert text into chunks of (embedding, chunk_text) pair.
```
Configuration parameters:
1. model, model name such as llama3.2
2. chunk_size = size of the chunk, default 1024
3. chunk_overlap, number of byte overlapping, default 20
4. address, address of the ollama host, default is http://localhost:11434.

Input:
input text in string format

Output:
JSON Array of Chunk with text and embedding. i.e.  [{"chunk":"chunk text 1", "embedding":[1,2....]},...]
```


3. 'generate' function to generate result from prompt.

```
Configuration parameters:
1. model, model name such as llama3.2
2. address, address of the ollama host, default is http://localhost:11434.

Input:
prompt in string format

Output:
JSON Array of generated text

