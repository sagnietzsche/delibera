## Step1: Defining the Shared State:
- Creating a clear domain model before coding 
- Need
    - Task:
        - id 
        - title 
        - description 
        - status: pending | planning | ready | in_progress | blocked | done | failed 
        - parent_task_id 
        - created_at 
        - updated_at 
    - Plan:
        - id
        - task_id
        - version 
        - status: draft | accepted | rejected | superseeded 
    - PlanStep:
        - id
        - plan_id
        - order 
        - description 
        - status 
        - assigned_node_id 
        - dependencies 
    - DeliberationEvent
        - id 
        - task_id 
        - author: node_id | agent_id | user 
        - type 
        - content 
        - created_at 
        - causal_parent_id / reply_to_id

- Starting with a append-only event log and materialize current state from the events.

## Step 2: Make The Deliberation Log Append-Only

  Treat the shared deliberation log as your product-level history.

  Example event types:

  task.created
  task.updated
  plan.proposed
  plan.revised
  plan.accepted
  step.added
  step.started
  step.completed
  step.blocked
  thought.appended
  evidence.added
  decision.proposed
  decision.accepted
  decision.rejected
  inference.requested
  inference.completed
  tool.result.recorded

  Important distinction:

  Raft log = infrastructure replication log
  Deliberation log = semantic application log

  The deliberation log should be written by commands that go through Raft.

  For example:

  AppendDeliberationEvent(task_id, event_type, content, author)

  That command gets committed through Raft, then each node applies it to its local state.

## Step 3: Add Commands, Not Direct Writes

  Define a small set of commands that can be replicated:

  CreateTask
  UpdateTaskStatus
  ProposePlan
  AcceptPlan
  AppendDeliberationEvent
  StartStep
  CompleteStep
  BlockStep
  RecordInferenceResult

  Every node should apply the same committed command and reach the same state.

  Avoid commands like:

  RunLLMAndCreatePlan

  because inference is nondeterministic. Instead split it:

  RequestPlanGeneration
  RecordPlanProposal

  The LLM output becomes data, then that data is committed.

## Step 4: Build A Minimal Planning Loop

  Start with one simple flow:

  User creates task
  Leader appends task.created
  Planner agent reads task
  Planner calls OpenAI
  Planner proposes plan
  Cluster commits plan.proposed
  A decision rule accepts/rejects plan
  Cluster commits plan.accepted
  Steps become executable

  Do not build multi-agent deliberation first. First build one planner producing one plan and writing it into the shared log.

  Once that works, add more agents.

## Step 5: Decide Your Deliberation Policy

  You need a rule for turning discussion into decisions.

  Simple options:

  1. Leader-decides
     The Raft leader or coordinator agent accepts a plan.
  2. Quorum vote
     Nodes vote on plan proposals.
  3. Evaluator-decides
     One evaluator model critiques the plan and accepts/rejects.
  4. User-approval
     Human approves the final plan.

  For now, I’d use:

  planner proposes -> evaluator critiques -> coordinator accepts

  Then all accepted decisions are committed through Raft.

## Step 6: Keep Inference Outside The Deterministic State Machine

  This is critical.

  Bad:

  Apply command:
    call OpenAI
    mutate state based on result

  Good:

  Node calls OpenAI outside Raft
  Node gets result
  Node submits RecordInferenceResult command
  Raft commits result
  All nodes apply result

  That way replicas do not diverge.

## Step 7: Add Read Models

  Once you have the event log, create materialized views:

  Current task tree
  Current accepted plan
  Timeline / deliberation feed
  Pending decisions
  Blocked steps
  Node activity

  These are derived from committed events. If they break, you can rebuild them from the deliberation log.

## Step 8: Minimum Viable Milestone

  Your next concrete milestone should be:

  Create task -> generate plan -> append shared deliberation events -> accept plan -> show task state consistently on all 3 nodes

  Do not start with fancy scheduling, multiple agents, retries, delegation, or autonomous execution yet.

  The first win is proving that all three nodes agree on:

  same task
  same plan
  same deliberation history
  same step statuses

  Implementation Order

  1. Define event schema.
  2. Define command schema.
  3. Add Raft command application for those events.
  4. Build task state projection from event log.
  5. Add planner inference that proposes a plan.
  6. Store the plan proposal as committed log data.
  7. Add simple accept/reject decision.
  8. Expose API/UI endpoints to inspect task, plan, and deliberation timeline.
  9. Add tests that kill/restart nodes and verify the same deliberation log is reconstructed.