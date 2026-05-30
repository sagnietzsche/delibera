## delibera

consensus driven multi-agent reasoning through committed deliberation

## Architecture 

For the high level architecture and why is this useuful : [architecture](./docs/ARCHITECTURE.md)

## Taskfile Reference

| Task | Description |
|---|---|
| `task build` | Full build — runs `fmt` + `vet` first, embeds git version via ldflags, outputs to `bin/delibera` |
| `task build:fast` | Build without fmt/vet for quick iteration |
| `task install` | `go install` to `$GOPATH/bin` |
| `task deps` | `go mod download` + `go mod tidy` |
| `task fmt` | `go fmt ./...` |
| `task vet` | `go vet ./...` |
| `task test` | Full test suite with `-race` and 60s timeout |
| `task test:short` | Short/fast tests only |
| `task clean` | Remove `bin/`, `dist/`, and `/tmp/delibera` cluster data |
| `task dev:node` | Build + start a single `node1` on `localhost:11000` |
| `task dev:cluster` | Build + start a 3-node local cluster |
| `task dev:ask -- "..."` | Build + run `ask` with your question |
| `task release` | Cross-compile for `linux/amd64`, `linux/arm64`, `darwin/arm64`, `darwin/amd64` into `dist/` |


## Package Structure 

- The Cobra app entrypoint is `main.go`, which calls the `cmd` package. Node startup logic lives in `internal/node`, and Cobra commands live directly under `cmd`.

- When using library, need to implement 3 main interfaces:
    - FSM (Finite State Machine) : Custom application logic. FSM applies the committed log entries to our actual system state.
    - LogStore : underlying storage mechanism to durably persist Raft log entries 
    - StableStore: used to store a stable Raft configuration state such as the current term and the voted candidates.

## References:

1. Improving Factuality and Reasoning in Language Models through Multiagent Debate (ICML 2024)
 - Why ? 
    - What if agents vote on reasoning steps rather than debate until convergence?
 - Ref:  https://arxiv.org/abs/2305.14325 
2. ReConcile: Round-Table Conference Improves Reasoning via Consensus among Diverse LLMs
 - Why ?
    - Round table of diverse models that can discuss, share confidence scores and reach a confidence-weighted consensus. 
    - They report improvements over both single-agent and prior multi-agent methods.
 - Ref: https://arxiv.org/abs/2309.13007 

3. CONSENSAGENT (ACL Findings 2025)
http://papers.cool/venue/2025.findings-acl.1141%40ACL

4. Can LLM Agents really debate ? 
https://huggingface.co/papers/2511.07784?

## Benchmarking Dataset
- MMLU
- TruthfulQA
- GSM8K

## LICENSE

MIT 
