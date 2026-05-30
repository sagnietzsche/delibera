## delibera

consensus driven multi-agent reasoning through committed deliberation


## Using `hashicorp/raft`
- When using library, need to implement 3 main interfaces:
    - FSM (Finite State Machine) : Custom application logic. FSM applies the committed log entries to our actual system state.
    - LogStore : underlying storage mechanism to durably persist Raft log entries 
    - StableStore: used to store a stable Raft configuration state such as the current term and the voted candidates.



## references:

1. Improving Factuality and Reasoning in Language Models through Multiagent Debate (ICML 2024)
 - What if agents vote on reasoning steps rather than debate until convergence?
2. ReConcile: Round-Table Conference Improves Reasoning via Consensus among Diverse LLMs
 - Round table of diverse models that can discuss, share confidence scores and reach a confidence-weighted consensus. 
 - They report improvements over both single-agent and prior multi-agent methods.

3. CONSENSAGENT (ACL Findings 2025)
http://papers.cool/venue/2025.findings-acl.1141%40ACL

4. Can LLM Agents really debate ? 
https://huggingface.co/papers/2511.07784?

## Benchmarking Dataset
- MMLU
- TruthfulQA
- GSM8K

