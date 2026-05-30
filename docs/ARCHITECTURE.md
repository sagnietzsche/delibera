## Core Idea:

- Agents do not directly produce the answers.
- Agents produce reasoning proposals.
- Reasoning proposals must be committedd by quorum before they become part of the shared reasoning state.

## HLD 

```
                 User Task
                      |
                      v
             +----------------+
             | Task Planner   |
             +----------------+
                      |
                      v
        +-----------------------------+
        | Shared Deliberation Log      |
        +-----------------------------+
                      |
       -----------------------------------
       |                |               |
       v                v               v
  Reasoner A      Reasoner B      Reasoner C
       |                |               |
       -----------------------------------
                      |
                      v
             Consensus Layer
                      |
          Commit / Reject Decision
                      |
                      v
        Shared Deliberation State
                      |
                      v
             Final Synthesizer
                      |
                      v
                 Final Answer
```

## Core Flow:
- Question -> Reasoning Step -> Consensus -> Commit -> Next Reasoning Step