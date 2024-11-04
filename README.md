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


## Integration with MatrixOne Database

```
mysql> create table chunk(t varchar, e vecf32(3072));

mysql> insert into chunk select json_unquote(json_extract(result, "$.chunk")), json_unquote(json_extract(result, "$.embedding")) from 
wasm_run_table('https://github.com/cpegeric/ollama-wasm/raw/main/ollama/ollama.wasm', 'embed', '{"model":"llama3.2"}', 'where is great wall?') as f;

mysql> select * from wasm_run_table('https://github.com/cpegeric/ollama-wasm/raw/main/ollama/ollama.wasm', 'generate', '{"model":"llama3.2"}', 'where is great
 wall?
| "The Great Wall of China is located in China, specifically along the northern border of the country. It stretches across several provinces and municipalities, including:\\n\\n1. Beijing Municipality (where the most famous and well-preserved sections are)\\n2. Tianjin Municipality\\n3. Hebei Province\\n4. Shanxi Province\\n5. Inner Mongolia Autonomous Region\\n\\nThe wall follows the mountain ranges of the northern China, winding its way through valleys and plains, and covers a total length of approximately 13,171 miles (21,196 km). It was built over several centuries to protect the Chinese Empire from invasions by nomadic tribes.\\n\\nSome popular locations to visit the Great Wall include:\\n\\n* Badaling Great Wall (near Beijing)\\n* Mutianyu Great Wall (also near Beijing)\\n* Jinshanling Great Wall (in Hebei Province)\\n* Simatai Great Wall (also in Hebei Province)\\n\\nIt's worth noting that while the wall is a UNESCO World Heritage Site and one of China's most famous landmarks, many sections have been destroyed or damaged over time due to natural erosion, human activities, or wars." |
1 row in set (3.85 sec)


```
