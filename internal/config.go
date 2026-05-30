package internal

// inference config

const (
	InferenceBaseURL   = "https://api.groq.com/openai/v1"
	InferenceAPIKeyEnv = "GROQ_API_KEY"
	Model              = "openai/gpt-oss-120b"
)

// command line defaults
const (
	DefaultHTTPAddr = "localhost:11000"
	DefaultRaftAddr = "localhost:12000"
	DefaultDataDir  = "/tmp/delibera"
)
