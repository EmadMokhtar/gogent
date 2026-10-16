# M0 · Contracts and Scaffolding Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create Gogent's multi-module repository, CI, milestone issues, and every frozen contract package, with working test models, so the eight wave-1 milestones can start in parallel.

**Architecture:** One core module (`github.com/EmadMokhtar/gogent`) holds the contract packages: `usage`, `message`, `model`, `run`, `tool`, `output`, `capability`, `runtime`, `event`, the root `gogent` skeleton, and `model/modeltest`. Nine empty nested modules reserve the places for providers, MCP, OpenTelemetry, durable execution, evals, and examples. The import graph runs one way only: `usage` ← `message` ← `model` ← `run` ← `tool`/`output` ← `capability` ← `runtime` ← root.

**Tech Stack:** Go 1.27; `github.com/google/jsonschema-go` v0.4.3 (tests only); golangci-lint v2.15.0; govulncheck v1.9.0; actionlint v1.7.12; GitHub Actions `actions/checkout@v7`, `actions/setup-go@v7`, `golangci/golangci-lint-action@v9`, `amannn/action-semantic-pull-request@v6`; MkDocs Material (skeleton only).

**Spec:** `docs/superpowers/specs/2026-10-16-gogent-design.md`. Read sections 4, 5, 7, and 8.3 before starting.

## Global Constraints

- Every `go.mod` says `go 1.27` and has no `toolchain` line.
- Root module path: `github.com/EmadMokhtar/gogent`.
- No pre-release versions (`-alpha`, `-beta`, `-rc`) in any `go.mod` or `go.sum`.
- At the end of M0 the core module depends only on the standard library, plus `github.com/google/jsonschema-go v0.4.3` for tests.
- Every exported name has a doc comment. The linter enforces this (revive rules `exported` and `package-comments`).
- Library code never logs and never prints.
- Messages and parts are stored and switched on as **values**, never pointers.
- Optional numeric settings are pointers; `nil` means "not set".
- All tests pass with `go test -race -count=1`.
- Commit messages and the PR title follow Conventional Commits with a lowercase description.
- The PR is created with `--assignee EmadMokhtar` (bare login, no `@`). Check the assignment with `gh pr view --json assignees`.
- The PR body contains `Closes #<M0 issue number>`. After creating it, check the link with `gh pr view --json closingIssuesReferences`.
- **API-first gate:** after Task 3, stop and wait for the owner's approval of the API commit before starting Task 4.

## Decisions This Plan Adds to the Spec

The spec leaves package-level placement to the builder. These choices change where some spec names live; the owner reviews them with this plan.

1. **New leaf package `usage`** holds `Usage` and `Limits`. The spec placed `Usage` in `model` and limits in `run`. Reason: `message.Response` carries usage, and `run` and `model` both need it, so a shared leaf package avoids an import cycle.
2. **`model.ToolDef`, `model.OutputMode`, and `model.OutputObject`** live in `model`, not in `tool` and `output`. Reason: model requests carry them, while `tool` and `output` import `run`, which imports `model`.
3. **Instructions are a field** (`message.Request.Instructions`), not a part. This follows PydanticAI: instructions from older requests are never re-sent.
4. **M0 hook stages:** run, model request, tool validate, tool execute, and output, plus the instructions, toolset, settings, and tool-preparation providers. The node, event-stream, and deferred-call stages, and the `Hooks` before/after helper, are added later as new optional interfaces (A1, A4, A5). Adding an optional interface does not break existing code.
5. **`usage.Usage.TotalTokens` is a method**, not a field, so it can never disagree with the input and output counts.
6. **Tool results carry usage** (`tool.Result.Usage`). A delegated agent's usage then flows back to the parent as data, which also works across durable-execution boundaries.

## Review Focus

1. **A caller stops a stream early** (`break` inside `for ... range m.RequestStream(...)`). Expected: no further `yield` call (Go panics if `yield` runs after it returned false). Test: Task 9, `TestFunctionModelStreamBreakEarly`.
2. **An explicit zero setting overrides a non-zero one** (agent temperature 0.7, run temperature 0). Expected: merged temperature 0. Test: Task 6, `TestMergeSettingsExplicitZeroOverrides`.
3. **History written by a newer Gogent contains an unknown `kind`.** Expected: `Unmarshal` returns an error that wraps `message.ErrUnknownKind` and names the kind; no panic. Test: Task 5, `TestUnmarshalUnknownPartKind`.
4. **Model IDs with extra colons or a missing half** (`"ollama:llama3:8b"`, `"anthropic"`, `"anthropic:"`, `":x"`). Expected: split on the first colon only; a missing half wraps `model.ErrInvalidID`. Test: Task 6, `TestParse`.
5. **A recursive JSON Schema** (`$ref` that points to itself through a required property). Expected: `modeltest.GenerateJSON` returns `modeltest.ErrSchemaTooDeep` instead of overflowing the stack. Test: Task 8, `TestGenerateRecursiveRefFails`.

---

## Before You Start

- Work in a new git worktree of `~/Projects/gogent`, on branch `feat/m0-contracts`, created from `main`.
- The design document and this plan must be on `main` first. If they are not, ask the owner to merge `docs/gogent-design` first, or branch from `docs/gogent-design` with the owner's agreement.
- Pushing over SSH may fail when 1Password cannot prompt. If it does, push over HTTPS: `git -c credential.helper='!gh auth git-credential' push https://github.com/EmadMokhtar/gogent.git feat/m0-contracts`.

---

### Task 1: Repository Scaffolding, CI, and Docs Skeleton

**Files:**
- Create: `go.mod`, `doc.go`, `Makefile`, `.golangci.yml`, `.github/workflows/ci.yml`, `.github/workflows/pr-title.yml`, `mkdocs.yml`, `docs/index.md`, `docs/contributing.md`, `CONTRIBUTING.md`
- Create: `go.mod` and `doc.go` in each of the nine nested module directories (table below)
- Modify: `README.md` (replace), `.gitignore` (append `site/`)

**Interfaces:**
- Produces: the `make` targets `test`, `lint`, `vuln`, `tidy`, `tidy-check`, `cover`, `record`, and `release-prep`, which every later milestone uses. It also produces the nine nested modules, so no session ever creates a module itself.

- [ ] **Step 1: Create the core module**

`go.mod` declares `module github.com/EmadMokhtar/gogent` and `go 1.27`. `doc.go` holds `package gogent` with a package comment: Gogent is a framework for building LLM agents in Go, modeled on PydanticAI; link to the design document.

- [ ] **Step 2: Create the nine nested modules**

| Directory | Module path suffix | Package | `replace` target | Built by |
|-----------|--------------------|---------|------------------|----------|
| `providers/anthropic` | `/providers/anthropic` | `anthropic` | `../..` | P0 |
| `providers/openai` | `/providers/openai` | `openai` | `../..` | P1 |
| `providers/google` | `/providers/google` | `google` | `../..` | P2 |
| `mcp` | `/mcp` | `mcp` | `..` | X1 |
| `otel` | `/otel` | `otel` | `..` | O1 |
| `durable/temporal` | `/durable/temporal` | `temporal` | `../..` | D1 |
| `durable/dbos` | `/durable/dbos` | `dbos` | `../..` | D2 |
| `evals` | `/evals` | `evals` | `..` | E1 |
| `examples` | `/examples` | `examples` | `..` | I1 |

Each `go.mod` has three lines of content: `module github.com/EmadMokhtar/gogent<suffix>`, `go 1.27`, and `replace github.com/EmadMokhtar/gogent => <replace target>`. It has **no** `require` line: `go mod tidy` adds one when the first import appears. Each `doc.go` has a package comment of the form "Package anthropic will hold the Anthropic model provider. Milestone P0 builds it."

- [ ] **Step 3: Write the `Makefile`**

The module list must skip `.git` and `.claude` (worktrees live under `.claude/worktrees/` and contain their own `go.mod` files).

```make
MODULES := $(shell find . -name go.mod -not -path './.git/*' -not -path './.claude/*' -exec dirname {} \; | sort)
PROVIDER_MODULES := $(filter ./providers/%,$(MODULES))
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.9.0
CORE := github.com/EmadMokhtar/gogent
```

| Target | Command run in every module (stop at the first failure) |
|--------|----------------------------------------------------------|
| `test` | `go test -race -count=1 ./...` |
| `lint` | `golangci-lint run ./...` |
| `vuln` | `go run $(GOVULNCHECK) ./...` |
| `tidy` | `go mod tidy` |
| `tidy-check` | `go mod tidy -diff` |
| `cover` | `go test -race -count=1 -coverprofile=coverage.out ./...`, then `go tool cover -func=coverage.out \| tail -1` |
| `record` | Only in `$(PROVIDER_MODULES)`: `GOGENT_RECORD=1 go test -tags live -count=1 ./...` |
| `release-prep` | Fails with a usage message unless `VERSION` is set. In every module except `.` whose `go.mod` mentions `$(CORE)`: `go mod edit -require=$(CORE)@$(VERSION) && go mod tidy` |

Each loop prints `== <module>` before it runs the command.

- [ ] **Step 4: Write `.golangci.yml`**

```yaml
version: "2"
linters:
  default: standard
  enable: [revive, errorlint, gocritic, misspell, unparam, bodyclose]
  settings:
    revive:
      rules:
        - name: exported
        - name: package-comments
formatters:
  enable: [gofmt, goimports]
```

- [ ] **Step 5: Run the local checks**

Run: `make tidy-check test lint vuln`
Expected: every module prints `== <dir>`; `test` reports `[no test files]`; lint reports `0 issues.`; govulncheck reports `No vulnerabilities found.`; exit code 0.

- [ ] **Step 6: Write the CI workflows**

`.github/workflows/ci.yml` runs on pushes to `main` and on pull requests, with `permissions: contents: read`. It has three jobs on `ubuntu-latest`. Each job starts with `actions/checkout@v7` and `actions/setup-go@v7` (`cache-dependency-path: '**/go.sum'`).

| Job | Go version | Steps after setup |
|-----|------------|-------------------|
| `test` | matrix `['1.27.x', 'stable']` | `make tidy-check`, then `make test` |
| `lint` | `1.27.x` | `golangci/golangci-lint-action@v9` with `version: v2.15.0` and `install-only: true`, then `make lint` |
| `vuln` | `1.27.x` | `make vuln` |

`stable` makes CI test each new Go release automatically when it ships (spec section 4).

`.github/workflows/pr-title.yml` runs on `pull_request_target` (types `opened`, `edited`, `synchronize`, `reopened`), with `permissions: pull-requests: read`. It runs one job with `amannn/action-semantic-pull-request@v6` and `GITHUB_TOKEN`. It never checks out code.

- [ ] **Step 7: Validate the workflows**

Run: `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12`
Expected: no output, exit code 0.

- [ ] **Step 8: Write the docs skeleton and README**

- `mkdocs.yml`: `site_name: Gogent`, `repo_url: https://github.com/EmadMokhtar/gogent`, `theme: {name: material}`, and `exclude_docs: superpowers/`. Leave out a `nav` key, so navigation follows the folder structure (spec decision 11).
- `docs/index.md`: two short paragraphs on what Gogent is, plus links to the design document and the contributing page.
- `docs/contributing.md`: rules 1 to 11 from spec section 7, copied word for word. Add a "Run the checks" section listing the `make` targets from Step 3.
- `CONTRIBUTING.md`: one line that links to `docs/contributing.md`.
- `README.md`: the title, one paragraph on what Gogent is, the status "Pre-alpha: milestone M0", and links to the design document and `CONTRIBUTING.md`.
- `.gitignore`: append `site/`.

- [ ] **Step 9: Commit**

```bash
git add -A
git commit -m "build: add module layout, makefile, ci and docs skeleton"
```

---

### Task 2: Milestone Issues and Labels

**Files:**
- Create: `scripts/milestone-issues.sh`

**Interfaces:**
- Produces: 21 GitHub issues titled `<ID> · <name>` (for example `M0 · Contracts and scaffolding`), labelled `milestone` and `wave-<N>`, and assigned to `EmadMokhtar`. The M0 issue number is used by Task 3.

- [ ] **Step 1: Write the script**

Behavior:
- Read the spec path from `$1` (default: `docs/superpowers/specs/2026-10-16-gogent-design.md`).
- Read ID-to-wave pairs from the table in section 8.2 (rows of the form `| A1 | Agent loop | 1 | ...`).
- For each heading `#### <ID> · <rest>` in sections 8.3 to 8.7:
  - The title is `<ID> · <name>`, where `<name>` is the text up to the next ` · ` (so `I1 · Vertical slice + orchestrator demo · Release v0.1.0` gives `I1 · Vertical slice + orchestrator demo`).
  - The body is every line up to the next line that starts with `#`, followed by a blank line and `Design: docs/superpowers/specs/2026-10-16-gogent-design.md`.
- Create labels with `gh label create <name> --force`: `milestone`, and `wave-0` through `wave-4`.
- Skip a milestone if `gh issue list --state all --search "\"<ID> ·\" in:title" --json number --jq length` prints a number other than `0`.
- Otherwise run `gh issue create --title ... --body ... --label milestone --label wave-<N> --assignee EmadMokhtar`.
- With `DRY_RUN=1`, print `<title> (wave <N>)` for each milestone and create nothing.
- Start the script with `set -euo pipefail`.

- [ ] **Step 2: Dry run**

Run: `DRY_RUN=1 bash scripts/milestone-issues.sh`
Expected: 21 lines, from `M0 · Contracts and scaffolding (wave 0)` to `R1 · Docs polish + release v0.3.0 (wave 4)`.

- [ ] **Step 3: Create the issues**

Run: `bash scripts/milestone-issues.sh`
Then: `gh issue list --label milestone --state all --json number --jq length`
Expected: `21`. Running the script a second time creates nothing.

- [ ] **Step 4: Commit**

```bash
git add scripts/milestone-issues.sh
git commit -m "chore: add script that creates milestone issues"
```

---

### Task 3: API-First Skeleton (Owner Gate)

This task writes every exported declaration in M0 with its doc comment. Function bodies are placeholders:

- In the lower packages, bodies `panic("gogent: not implemented")`. Tasks 4 to 9 replace every one of them.
- In the root package, bodies return `ErrNotImplemented`, and option constructors return options that do nothing. These stay after M0 is merged; A1 replaces them.
- `Example` functions are included, but without `// Output:` lines yet. Tasks 4 to 9 add those lines.

**Files:** one file per concern, as listed below each package.

**Interfaces:**
- Produces: everything below. Every later task and milestone uses these exact names.

#### Package `usage` (`usage/usage.go`, `usage/limits.go`)

```go
type Usage struct {
	Requests         int            `json:"requests,omitzero"`
	ToolCalls        int            `json:"tool_calls,omitzero"`
	InputTokens      int            `json:"input_tokens,omitzero"`
	OutputTokens     int            `json:"output_tokens,omitzero"`
	CacheReadTokens  int            `json:"cache_read_tokens,omitzero"`
	CacheWriteTokens int            `json:"cache_write_tokens,omitzero"`
	Details          map[string]int `json:"details,omitzero"` // provider-specific counts, such as reasoning tokens
}
func (u Usage) Add(o Usage) Usage // field-wise sum; Details merged by key into a NEW map
func (u Usage) TotalTokens() int  // InputTokens + OutputTokens

type Limits struct { // zero means "no limit" for every field
	Requests, ToolCalls, InputTokens, OutputTokens, TotalTokens int
}
func DefaultLimits() Limits // Limits{Requests: 50}, the PydanticAI default
func (l Limits) CheckBeforeRequest(u Usage) error          // fails if u.Requests+1 > l.Requests
func (l Limits) CheckBeforeToolCalls(u Usage, n int) error // fails if u.ToolCalls+n > l.ToolCalls
func (l Limits) CheckTokens(u Usage) error                 // checks input, then output, then total

type LimitExceededError struct {
	Limit string // "requests", "tool_calls", "input_tokens", "output_tokens" or "total_tokens"
	Max   int
	Value int
}
func (e *LimitExceededError) Error() string // usage: limit "requests" exceeded (51 > 50)
```

#### Package `message` (`message/message.go`, `message/parts.go`, `message/json.go`)

```go
type Message interface{ isMessage() } // implemented by Request and Response (value receivers)

type Request struct {
	Parts        []RequestPart
	Instructions string // instructions for this request only; never re-sent from history
	Timestamp    time.Time
}
type Response struct {
	Parts              []ResponsePart
	Usage              usage.Usage // Requests is 1 for every response
	ModelName          string
	ProviderName       string
	ProviderResponseID string
	FinishReason       FinishReason
	Timestamp          time.Time
}

type RequestPart interface{ isRequestPart() }
type ResponsePart interface{ isResponsePart() }
type UserContent interface{ isUserContent() } // A6 adds image, audio, video and document content

type SystemPromptPart struct{ Content string }
type UserPromptPart struct{ Content []UserContent }
type TextContent struct{ Text string }
type ToolReturnPart struct {
	ToolName, ToolCallID string
	Content              json.RawMessage // the tool's result as JSON
	IsError              bool            // true when the tool failed with tool.Fail
}
type RetryPromptPart struct{ ToolName, ToolCallID, Content string } // ToolName is empty when not about a tool
type TextPart struct{ Content string }
type ThinkingPart struct{ Content, Signature, Provider string } // Provider names who made Signature
type ToolCallPart struct {
	ToolName, ToolCallID string
	Args                 json.RawMessage
}

type FinishReason string
const (
	FinishStop          FinishReason = "stop"
	FinishLength        FinishReason = "length"
	FinishToolCall      FinishReason = "tool_call"
	FinishContentFilter FinishReason = "content_filter"
	FinishError         FinishReason = "error"
)

func UserPrompt(text string) UserPromptPart // one TextContent

func Marshal(msgs []Message) ([]byte, error) // JSON array; nil gives []
func Unmarshal(data []byte) ([]Message, error)
var ErrUnknownKind = errors.New("message: unknown kind")
```

JSON encoding rules (Task 5 implements them):
- Every message, part, and user content object has a `kind` field.
  - Message kinds: `request`, `response`.
  - Request part kinds: `system-prompt`, `user-prompt`, `tool-return`, `retry-prompt`.
  - Response part kinds: `text`, `thinking`, `tool-call`.
  - User content kinds: `text`.
- Field names are snake_case: `parts`, `instructions`, `timestamp`, `usage`, `model_name`, `provider_name`, `provider_response_id`, `finish_reason`, `content`, `text`, `tool_name`, `tool_call_id`, `args`, `is_error`, `signature`, `provider`.
- Empty optional fields are left out (`omitzero`). `parts` is always present.

#### Package `model` (`model/model.go`, `model/params.go`, `model/settings.go`, `model/stream.go`, `model/errors.go`, `model/registry.go`)

```go
type Model interface {
	Name() string     // for example "claude-sonnet-5-5"
	Provider() string // for example "anthropic"; the registry key
	Profile() Profile
	Request(ctx context.Context, msgs []message.Message, s Settings, p RequestParams) (message.Response, error)
	// RequestStream yields PartStart and PartDelta events, then exactly one Done
	// as its last event. It yields (nil, err) once on failure and then stops.
	RequestStream(ctx context.Context, msgs []message.Message, s Settings, p RequestParams) iter.Seq2[StreamEvent, error]
}

type Profile struct {
	SupportsTools, SupportsNativeOutput, SupportsThinking bool
	DefaultOutputMode                                     OutputMode
}

type OutputMode string
const (OutputText OutputMode = "text"; OutputTool OutputMode = "tool"; OutputNative OutputMode = "native"; OutputPrompted OutputMode = "prompted")

type ToolKind string
const (ToolFunction ToolKind = "function"; ToolOutput ToolKind = "output"; ToolExternal ToolKind = "external"; ToolUnapproved ToolKind = "unapproved")

type ToolDef struct {
	Name, Description string
	Parameters        json.RawMessage // a JSON Schema object
	Strict            *bool           // nil means the provider's default
	Kind              ToolKind
	Sequential        bool // must not run at the same time as other tools
}
type OutputObject struct {
	Name, Description string
	Schema            json.RawMessage
	Strict            *bool
}
type RequestParams struct {
	FunctionTools, OutputTools []ToolDef
	OutputMode                 OutputMode
	OutputObject               *OutputObject // for native and prompted modes
	AllowTextOutput            bool
}

type Settings struct {
	MaxTokens         *int
	Temperature, TopP *float64
	Seed              *int
	StopSequences     []string
	ParallelToolCalls *bool
	Timeout           time.Duration
	ProviderSettings  []any // provider-specific option structs; a provider uses the last value of its own type
}
func MergeSettings(layers ...Settings) Settings // order: model, agent, capability, run
func Ptr[T any](v T) *T

type StreamEvent interface{ isStreamEvent() }
type PartStart struct{ Index int; Part message.ResponsePart }
type PartDelta struct{ Index int; Delta Delta }
type Done struct{ Response message.Response }
type Delta interface{ isDelta() }
type TextDelta struct{ Content string }
type ThinkingDelta struct{ Content, Signature string }
type ToolCallDelta struct{ ToolName, ToolCallID, ArgsJSON string } // ArgsJSON is a fragment
func Collect(stream iter.Seq2[StreamEvent, error]) (message.Response, error)
var ErrIncompleteStream = errors.New("model: stream ended without a Done event")

type APIError struct {
	Provider   string
	StatusCode int // 0 when there was no HTTP response
	Message    string
	Retryable  bool
	RetryAfter time.Duration
	Err        error
}
func (e *APIError) Error() string
func (e *APIError) Unwrap() error

type Factory func(name string) (Model, error)
func Register(provider string, f Factory) // panics on an empty provider, a nil factory or a duplicate
func Parse(id string) (Model, error)      // "provider:name"; splits on the first colon
func Providers() []string                 // sorted
var (
	ErrUnknownProvider = errors.New("model: unknown provider")
	ErrInvalidID       = errors.New("model: invalid model id")
)
```

#### Package `run` (`run/run.go`)

```go
type Info struct {
	RunID      string
	Model      model.Model
	Usage      usage.Usage
	Messages   []message.Message // read-only for tools and hooks
	Step       int
	ToolName   string // set while a tool runs
	ToolCallID string
	Retry      int // retries used so far by the current tool or output
	MaxRetries int
}
type Untyped struct { Info; Deps any }
type Context[D any] struct { Info; Deps D }
func Typed[D any](u *Untyped) (*Context[D], error) // nil Deps gives the zero D
func (c *Context[D]) Untyped() *Untyped
var ErrDepsType = errors.New("run: dependencies have the wrong type")
```

#### Package `tool` (`tool/tool.go`, `tool/errors.go`, `tool/deferred.go`)

```go
type Tool interface {
	Def() model.ToolDef
	Call(ctx context.Context, rc *run.Untyped, args json.RawMessage) (Result, error)
}
type Result struct {
	Content json.RawMessage // sent back to the model
	Usage   usage.Usage     // usage spent inside the tool, such as a delegated agent's; added to the run
}
type Toolset interface {
	Tools(ctx context.Context, rc *run.Untyped) ([]Tool, error)
}
type DepsTyped interface{ DepsType() reflect.Type } // optional; gogent.New checks it

type RetryError struct{ Message string }
type FailError struct{ Message string }
type ApprovalRequiredError struct{ Metadata map[string]any }
type CallDeferredError struct{ Metadata map[string]any }
// Error() texts: "tool: retry requested: <Message>", "tool: failed: <Message>",
// "tool: approval required", "tool: call deferred".
func Retry(msg string) error                     // returns *RetryError
func Fail(msg string) error                      // returns *FailError
func ApprovalRequired(meta map[string]any) error // returns *ApprovalRequiredError
func CallDeferred(meta map[string]any) error     // returns *CallDeferredError

type DeferredRequests struct {
	Calls     []message.ToolCallPart    // external tools to run outside the agent
	Approvals []message.ToolCallPart    // tools waiting for approval
	Metadata  map[string]map[string]any // keyed by tool call ID
}
```

#### Package `output` (`output/output.go`)

```go
type Spec interface {
	OutputType() reflect.Type
	Params(p model.Profile) (Params, error)
	// Process turns a candidate final answer into the output value.
	// A *tool.RetryError asks the model to try again.
	Process(ctx context.Context, rc *run.Untyped, c Candidate) (any, error)
}
type Params struct {
	Mode      model.OutputMode
	Tools     []model.ToolDef
	Object    *model.OutputObject
	AllowText bool
}
type Candidate struct {
	ToolCall *message.ToolCallPart // set in tool mode
	Text     string                // set in text, native and prompted modes
	Partial  bool                  // true while streaming (A4)
}
```

#### Package `capability` (`capability/capability.go`, `capability/chain.go`)

```go
type Capability interface{ ID() string }

type InstructionsProvider interface {
	Instructions(ctx context.Context, rc *run.Untyped) (string, error)
}
type ToolsetProvider interface{ Toolsets() []tool.Toolset }
type SettingsProvider interface {
	Settings(ctx context.Context, rc *run.Untyped) (model.Settings, error)
}
type ToolPreparer interface {
	PrepareTools(ctx context.Context, rc *run.Untyped, defs []model.ToolDef) ([]model.ToolDef, error)
}

type RunHandler func(ctx context.Context, rc *run.Untyped) (any, error)
type RunWrapper interface{ WrapRun(next RunHandler) RunHandler }

type ModelRequest struct {
	Model    model.Model
	Messages []message.Message
	Settings model.Settings
	Params   model.RequestParams
}
type ModelRequestHandler func(ctx context.Context, rc *run.Untyped, req ModelRequest) (message.Response, error)
type ModelRequestWrapper interface{ WrapModelRequest(next ModelRequestHandler) ModelRequestHandler }

type ToolValidateHandler func(ctx context.Context, rc *run.Untyped, call message.ToolCallPart, def model.ToolDef) (json.RawMessage, error)
type ToolValidateWrapper interface{ WrapToolValidate(next ToolValidateHandler) ToolValidateHandler }

type ToolExecuteHandler func(ctx context.Context, rc *run.Untyped, call message.ToolCallPart, def model.ToolDef, args json.RawMessage) (tool.Result, error)
type ToolExecuteWrapper interface{ WrapToolExecute(next ToolExecuteHandler) ToolExecuteHandler }

type OutputHandler func(ctx context.Context, rc *run.Untyped, c output.Candidate) (any, error)
type OutputWrapper interface{ WrapOutput(next OutputHandler) OutputHandler }

// Each Chain function wraps final with every capability that implements the
// matching wrapper interface. caps[0] is the outermost wrapper. Nil entries
// and capabilities without the interface are skipped.
func ChainRun(caps []Capability, final RunHandler) RunHandler
func ChainModelRequest(caps []Capability, final ModelRequestHandler) ModelRequestHandler
func ChainToolValidate(caps []Capability, final ToolValidateHandler) ToolValidateHandler
func ChainToolExecute(caps []Capability, final ToolExecuteHandler) ToolExecuteHandler
func ChainOutput(caps []Capability, final OutputHandler) OutputHandler
```

#### Package `runtime` (`runtime/runtime.go`)

Types only. A1 builds the default runtime. The name matches the spec. Importers that also need the standard library's `runtime` use an alias.

```go
type Effect interface{ isEffect() }
type RequestModel struct{ Messages []message.Message; Step int }
type CallTool struct{ Call message.ToolCallPart; Step int }
type ProcessOutput struct{ Candidate output.Candidate; Step int }

type Outcome interface{ isOutcome() }
type ModelOutcome struct{ Request message.Request; Response message.Response } // Request is what was actually sent
type ToolOutcome struct { // exactly one of Return, RetryPrompt or Deferred is set
	Return      *message.ToolReturnPart
	RetryPrompt *message.RetryPromptPart
	Deferred    *Deferred
	Usage       usage.Usage
}
type DeferredKind string
const (DeferredApproval DeferredKind = "approval"; DeferredExternal DeferredKind = "external")
type Deferred struct {
	Call     message.ToolCallPart
	Kind     DeferredKind
	Metadata map[string]any
}
type OutputOutcome struct { // Output must be JSON-serializable for durable runs
	Output      any
	RetryPrompt *message.RetryPromptPart
}

type Handler interface { // performs one effect; built by the agent
	RequestModel(ctx context.Context, e RequestModel) (ModelOutcome, error)
	CallTool(ctx context.Context, e CallTool) (ToolOutcome, error)
	ProcessOutput(ctx context.Context, e ProcessOutput) (OutputOutcome, error)
}
type Runtime interface {
	// Execute performs a batch of effects, possibly at the same time, and
	// returns one outcome per effect in the same order. A runtime is bound to
	// one run and its context when it is created.
	Execute(batch []Effect) ([]Outcome, error)
}
```

#### Package `event` (`event/event.go`)

```go
type Event interface{ isEvent() } // A4 adds the event types
```

#### Root package `gogent` (`agent.go`, `options.go`, `result.go`, `node.go`, `errors.go`, `stream.go`)

`stream.go` holds only `RunStream`, because A4 owns `stream*.go`.

```go
var (
	ErrNotImplemented   = errors.New("gogent: not implemented")
	ErrConcurrencyLimit = errors.New("gogent: concurrency queue is full")
)
type UnexpectedModelBehaviorError struct{ Message string } // "gogent: unexpected model behavior: <Message>"
type PanicError struct{ Value any; Stack []byte }          // "gogent: panic: <Value>"

type Agent[D, O any] struct{ /* unexported fields */ }
func New[D, O any](m model.Model, opts ...Option) (*Agent[D, O], error)
func (a *Agent[D, O]) Name() string
func (a *Agent[D, O]) Run(ctx context.Context, deps D, prompt string, opts ...RunOption) (*Result[O], error)
func (a *Agent[D, O]) RunAs[T any](ctx context.Context, deps D, prompt string, opts ...RunOption) (*Result[T], error)
func (a *Agent[D, O]) Iter(ctx context.Context, deps D, prompt string, opts ...RunOption) iter.Seq2[Node, error]
func (a *Agent[D, O]) RunStream(ctx context.Context, deps D, prompt string, opts ...RunOption) iter.Seq2[event.Event, error]

type Result[O any] struct {
	Output   O
	Deferred *tool.DeferredRequests // non-nil when the run paused (A5)
	/* unexported: all messages, index of the first new message, usage */
}
func (r *Result[O]) Usage() usage.Usage
func (r *Result[O]) AllMessages() []message.Message
func (r *Result[O]) NewMessages() []message.Message
func (r *Result[O]) Paused() bool

type Node interface{ isNode() }
type UserPromptNode struct{ Prompt string }
type ModelRequestNode struct{ Request message.Request }
type CallToolsNode struct{ Response message.Response }
type EndNode[O any] struct{ Result *Result[O] }

type EndStrategy string
const (EndGraceful EndStrategy = "graceful"; EndEarly EndStrategy = "early"; EndExhaustive EndStrategy = "exhaustive")

type Option func(*agentConfig) // agentConfig is unexported
func WithName(name string) Option
func WithInstructions(text string) Option
func WithInstructionsFunc[D any](f func(ctx context.Context, rc *run.Context[D]) (string, error)) Option
func WithSystemPrompt(text string) Option
func WithTools(tools ...tool.Tool) Option
func WithToolsets(sets ...tool.Toolset) Option
func WithOutput(spec output.Spec) Option
func WithOutputValidator[D, O any](f func(ctx context.Context, rc *run.Context[D], out O) (O, error)) Option
func WithCapabilities(caps ...capability.Capability) Option
func WithSettings(s model.Settings) Option
func WithLimits(l usage.Limits) Option
func WithToolRetries(n int) Option   // default 1
func WithOutputRetries(n int) Option // default 1
func WithEndStrategy(s EndStrategy) Option // default EndGraceful
func WithConcurrencyLimit(running, queued int) Option

type RunOption func(*runConfig) // runConfig is unexported
func WithHistory(msgs []message.Message) RunOption
func WithRunInstructions(text string) RunOption
func WithRunSettings(s model.Settings) RunOption
func WithRunLimits(l usage.Limits) RunOption
func WithRunModel(m model.Model) RunOption
func WithRunOutput(spec output.Spec) RunOption
```

#### Package `model/modeltest` (`model/modeltest/function.go`, `model/modeltest/testmodel.go`, `model/modeltest/generate.go`, `model/modeltest/options.go`)

```go
type Info struct {
	Settings model.Settings
	Params   model.RequestParams
}
type Func func(ctx context.Context, msgs []message.Message, info Info) (message.Response, error)

type FunctionModel struct{ /* unexported */ } // Name() "function", Provider() "function"
func NewFunctionModel(f Func, opts ...Option) *FunctionModel

type TestModel struct{ /* unexported */ } // Name() "test", Provider() "test"
func NewTestModel(opts ...Option) *TestModel
const DefaultText = "success (no tool calls)"

type Option func(*config) // shared; options a model does not use are ignored
func WithName(name string) Option
func WithProfile(p model.Profile) Option // default: tools and native output supported, thinking not, DefaultOutputMode OutputTool
func WithText(text string) Option        // TestModel's final text; default DefaultText
func WithToolCalls(names ...string) Option // TestModel calls only these function tools; default all

func GenerateJSON(schema json.RawMessage) (json.RawMessage, error)
var (
	ErrSchemaTooDeep = errors.New("modeltest: schema is too deep or recursive")
	ErrUnsatisfiable = errors.New("modeltest: schema accepts no value")
)
```

Add `var _ model.Model = (*FunctionModel)(nil)` and `var _ model.Model = (*TestModel)(nil)` compile-time checks.

- [ ] **Step 1: Write all declarations above with doc comments, placeholder bodies, and these Example functions:** `message.ExampleMarshal`, `model.ExampleParse`, `model.ExampleMergeSettings`, `modeltest.ExampleNewFunctionModel`.

- [ ] **Step 2: Check it builds and lints**

Run: `go build ./... && go vet ./... && make lint`
Expected: exit code 0, `0 issues.` in every module.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "feat: add m0 public api skeleton"
```

- [ ] **Step 4: Push and open a draft PR**

Run:
```bash
gh pr create --draft --base main --title "feat: add m0 contracts and scaffolding" --assignee EmadMokhtar --body "Closes #<M0 issue number>"
```
Then: `gh pr view --json assignees,closingIssuesReferences`
Expected: assignee `EmadMokhtar`, and one closing issue (the M0 issue).

- [ ] **Step 5: Owner gate**

Ask the owner to review the API commit. Stop until the owner approves. Apply any requested changes as new commits on the same branch, and record each change in the PR description. Then continue with Task 4.

---

### Task 4: Usage Totals and Limits

**Files:**
- Modify: `usage/usage.go`, `usage/limits.go`
- Test: `usage/usage_test.go`, `usage/limits_test.go`

**Interfaces:**
- Consumes: Task 3 `usage` declarations.
- Produces: working `Add`, `TotalTokens`, `DefaultLimits`, the `Check*` methods, and `LimitExceededError`.

- [ ] **Step 1: Write the failing tests**

- `TestAddSumsAllFields`: `Usage{Requests:1, InputTokens:10, Details: {"reasoning":2}}.Add(Usage{Requests:1, OutputTokens:5, Details: {"reasoning":3, "audio":1}})` equals `Usage{Requests:2, InputTokens:10, OutputTokens:5, Details: {"reasoning":5, "audio":1}}`.
- `TestAddDoesNotAliasDetails`: after `c := a.Add(b)`, setting `c.Details["x"] = 9` leaves `a.Details` and `b.Details` unchanged.
- `TestTotalTokens`: `Usage{InputTokens: 7, OutputTokens: 3}.TotalTokens() == 10`.
- `TestDefaultLimits`: equals `Limits{Requests: 50}`.
- `TestZeroLimitMeansNoLimit`: `Limits{}` returns nil from all three checks for `Usage{Requests: 1_000_000, ToolCalls: 1_000_000, InputTokens: 1_000_000}` and `n = 1_000_000`.
- `TestCheckBeforeRequest`: with `Limits{Requests: 50}`, `Usage{Requests: 49}` gives nil. `Usage{Requests: 50}` gives a `*LimitExceededError` (found with `errors.As`) equal to `{Limit: "requests", Max: 50, Value: 51}`.
- `TestCheckBeforeToolCalls`: `Limits{ToolCalls: 3}`, `Usage{ToolCalls: 1}`, `n = 3` gives `{Limit: "tool_calls", Max: 3, Value: 4}`. `n = 2` gives nil.
- `TestCheckTokensOrder`: `Limits{InputTokens: 10, TotalTokens: 10}` with `Usage{InputTokens: 11, OutputTokens: 5}` reports `Limit == "input_tokens"`.
- `TestLimitExceededErrorText`: `(&LimitExceededError{Limit: "requests", Max: 50, Value: 51}).Error()` equals `usage: limit "requests" exceeded (51 > 50)`.

- [ ] **Step 2: Run them to see them fail**

Run: `go test ./usage/`
Expected: FAIL (panic: gogent: not implemented).

- [ ] **Step 3: Implement the bodies.** `Add` allocates a new `Details` map only when either side has entries.

- [ ] **Step 4: Run them to see them pass**

Run: `go test -race ./usage/`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add usage/
git commit -m "feat(usage): add usage totals and limits"
```

---

### Task 5: Message JSON Encoding

**Files:**
- Modify: `message/message.go`, `message/parts.go`, `message/json.go`
- Test: `message/json_test.go`, `message/testdata/history.golden.json`, `message/example_test.go`

**Interfaces:**
- Consumes: Task 3 `message` declarations, `usage.Usage`.
- Produces: `Marshal`, `Unmarshal`, `UserPrompt`, and the JSON format that D1, D2, and A5 depend on.

- [ ] **Step 1: Write the failing tests**

Use one fixture history, in this order:
1. A `Request` with a `SystemPromptPart`, `UserPrompt("Weather in Amsterdam?")`, `Instructions`, and a timestamp from `time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)`.
2. A `Response` with a `ThinkingPart`, a `TextPart`, and `ToolCallPart{Args: json.RawMessage(`{"city":"Amsterdam"}`)}`, plus full `Usage`, model and provider names, and `FinishToolCall`.
3. A `Request` with a `ToolReturnPart{IsError: true}` and a `RetryPromptPart`.

Write all JSON in fixtures **compact**, because `json.RawMessage` is compacted when encoded.

- `TestRoundTrip`: `Unmarshal(Marshal(fixture))` is `reflect.DeepEqual` to the fixture.
- `TestGoldenHistory`: `Marshal(fixture)`, indented with two spaces, equals `testdata/history.golden.json` byte for byte. An `-update` flag rewrites the file.
- `TestKindValues`: each part type marshals with the kind string listed in Task 3.
- `TestUnmarshalUnknownPartKind`: `[{"kind":"request","parts":[{"kind":"video-call"}]}]` gives an error where `errors.Is(err, ErrUnknownKind)` is true and `err.Error()` contains `video-call`.
- `TestUnmarshalUnknownMessageKind`: `[{"kind":"memo"}]` gives the same kind of error, containing `memo`.
- `TestUnmarshalMissingKind`: `[{"parts":[]}]` gives an `ErrUnknownKind` error.
- `TestEmptyHistory`: `Marshal(nil)` gives `[]`, and `Unmarshal([]byte("[]"))` gives a slice of length 0 with no error.
- `TestUserPrompt`: `UserPrompt("hi")` equals `UserPromptPart{Content: []UserContent{TextContent{Text: "hi"}}}`.
- `ExampleMarshal`: marshals `[]Message{Request{Parts: []RequestPart{UserPrompt("hi")}}}` and prints `[{"kind":"request","parts":[{"kind":"user-prompt","content":[{"kind":"text","text":"hi"}]}]}]`.

- [ ] **Step 2: Run them to see them fail**

Run: `go test ./message/`
Expected: FAIL

- [ ] **Step 3: Implement the encoding.**
- Each concrete type gets a `MarshalJSON` that adds its `kind`.
- `Request`, `Response`, and `UserPromptPart` get an `UnmarshalJSON` that reads `parts` or `content` as `[]json.RawMessage`. It decodes `kind` first, then dispatches through one unexported table that maps kind to a decode function.
- Errors wrap `ErrUnknownKind` as `fmt.Errorf("%w: %q", ErrUnknownKind, kind)`.

- [ ] **Step 4: Run them to see them pass**

Run: `go test -race ./message/`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add message/
git commit -m "feat(message): add stable json encoding for message history"
```

---

### Task 6: Model Settings, Registry, Stream Collector, and API Errors

**Files:**
- Modify: `model/settings.go`, `model/registry.go`, `model/stream.go`, `model/errors.go`
- Test: `model/settings_test.go`, `model/registry_test.go`, `model/stream_test.go`, `model/errors_test.go`, `model/example_test.go`

**Interfaces:**
- Consumes: Task 3 `model` declarations, `message.Response`.
- Produces: `MergeSettings`, `Ptr`, `Register`, `Parse`, `Providers`, `Collect`, and `APIError`.

- [ ] **Step 1: Write the failing tests**

The registry is global, so each test registers provider names that no other test uses (prefix `fake-<testname>`). Fake models come from a small test-only type in `registry_test.go`.

- `TestMergeSettingsLaterWins`: `MergeSettings(Settings{Temperature: Ptr(0.2), MaxTokens: Ptr(100)}, Settings{Temperature: Ptr(0.7)})` gives temperature 0.7 and max tokens 100.
- `TestMergeSettingsExplicitZeroOverrides`: `MergeSettings(Settings{Temperature: Ptr(0.7)}, Settings{Temperature: Ptr(0.0)})` gives temperature 0.0.
- `TestMergeSettingsNilKeeps`: a later layer with a nil `Temperature` keeps the earlier value.
- `TestMergeSettingsSlicesAndTimeout`: a later non-nil `StopSequences` replaces the earlier one, and a nil one keeps it. A later non-zero `Timeout` wins. `ProviderSettings` from all layers are concatenated in order.
- `TestMergeSettingsDoesNotAlias`: changing `result.StopSequences[0]` leaves the input layer unchanged.
- `TestParse` (table):

  | Input | Result |
  |-------|--------|
  | `fake-parse:m1` | provider `fake-parse`, name `m1` |
  | `fake-parse:llama3:8b` | name `llama3:8b` |
  | `fake-parse` | error wrapping `ErrInvalidID` |
  | `fake-parse:` | error wrapping `ErrInvalidID` |
  | `:m1` | error wrapping `ErrInvalidID` |
  | empty string | error wrapping `ErrInvalidID` |
  | `nope:m1` | error wrapping `ErrUnknownProvider`, with text containing `"nope"` and `forgot to import the provider package?` |

- `TestRegisterPanics`: the panic messages are exactly `model: Register provider name is empty`, `model: Register factory is nil`, and `model: Register called twice for provider "fake-dup"`.
- `TestProvidersSorted`: after registering `fake-z` and `fake-a`, `Providers()` lists `fake-a` before `fake-z`.
- `TestRegistryConcurrent`: 50 goroutines call `Parse` while another goroutine registers 50 new providers. The test passes under `-race`.
- `TestCollect`:
  - A stream of `PartStart`, `PartDelta`, and `Done{Response: r}` gives `r`.
  - A stream that yields `(nil, errBoom)` gives `errBoom`.
  - A stream that ends without `Done` gives `ErrIncompleteStream`.
  - A stream that emits events after `Done` is not read past `Done`, and the test's iterator records that `yield` returned false.
- `TestAPIErrorText`: `&APIError{Provider: "anthropic", StatusCode: 429, Message: "rate limited"}` prints `model: anthropic API error 429: rate limited`. With `StatusCode: 0` it prints `model: anthropic API error: rate limited`. `errors.Unwrap` returns `Err`.
- `ExampleParse`: registers `example` with a factory for a fake model, calls `Parse("example:small")`, and prints `example small`.
- `ExampleMergeSettings`: prints `0.7 100`.

- [ ] **Step 2: Run them to see them fail**

Run: `go test ./model/`
Expected: FAIL

- [ ] **Step 3: Implement the bodies.**
- The registry uses a `sync.RWMutex` around a `map[string]Factory`.
- `Parse` wraps errors as `fmt.Errorf("%w: %q", ErrInvalidID, id)` and `fmt.Errorf("%w: %q (forgot to import the provider package?)", ErrUnknownProvider, provider)`.
- `MergeSettings` copies slices so the result never shares memory with a layer.

- [ ] **Step 4: Run them to see them pass**

Run: `go test -race ./model/`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add model/
git commit -m "feat(model): add settings merge, model registry and stream collector"
```

---

### Task 7: Tool Errors, Typed Run Context, and Capability Chains

**Files:**
- Modify: `tool/errors.go`, `run/run.go`, `capability/chain.go`
- Test: `tool/errors_test.go`, `run/run_test.go`, `capability/chain_test.go`

**Interfaces:**
- Consumes: Task 3 `tool`, `run`, and `capability` declarations.
- Produces: working error constructors, `run.Typed`, `(*run.Context[D]).Untyped`, and the five `Chain*` functions, with the ordering A1 and D1 rely on.

- [ ] **Step 1: Write the failing tests**

- `TestToolErrors` (table): `Retry("bad city")`, `Fail("no network")`, `ApprovalRequired(map[string]any{"why": "cost"})`, and `CallDeferred(nil)` each match their pointer type through `errors.As`, also after wrapping with `fmt.Errorf("x: %w", err)`. Their `Error()` texts are exactly `tool: retry requested: bad city`, `tool: failed: no network`, `tool: approval required`, and `tool: call deferred`.
- `TestTypedOK`: `Typed[myDeps](&Untyped{Deps: myDeps{N: 1}, Info: Info{Step: 2}})` gives `Deps.N == 1` and `Step == 2`.
- `TestTypedWrongType`: `Typed[myDeps](&Untyped{Deps: "text"})` gives an error where `errors.Is(err, ErrDepsType)` is true, and the text contains `run.myDeps` and `string`.
- `TestTypedNilDeps`: `Typed[myDeps](&Untyped{})` gives the zero `myDeps` and no error. `Typed[*myDeps](&Untyped{})` gives a nil pointer and no error.
- `TestUntypedRoundTrip`: `c.Untyped().Deps` equals `c.Deps`, and the `Info` values are equal.
- `TestChainOrder`: one table-driven case per Chain function. Each case uses two capabilities `A` and `B` whose wrappers append `A>` before calling `next` and `<A` after it, with a final handler that appends `final`. The log must be `[A> B> final <B <A]`.
- `TestChainSkips`: a `nil` entry and a capability that implements no wrapper are skipped, and the log is `[A> final <A]`.
- `TestChainShortCircuit`: if `A` returns without calling `next`, the final handler never runs and `A`'s return value comes back.
- `TestChainEmpty`: with no capabilities, the returned handler calls the final handler exactly once.

- [ ] **Step 2: Run them to see them fail**

Run: `go test ./tool/ ./run/ ./capability/`
Expected: FAIL

- [ ] **Step 3: Implement the bodies.** `Typed` reports the wanted type with `reflect.TypeFor[D]()` and the actual type with `%T`. The five Chain functions share one unexported generic helper.

- [ ] **Step 4: Run them to see them pass**

Run: `go test -race ./tool/ ./run/ ./capability/`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add tool/ run/ capability/
git commit -m "feat: add tool errors, typed run context and capability chains"
```

---

### Task 8: Example JSON From a Schema

**Files:**
- Modify: `model/modeltest/generate.go`
- Test: `model/modeltest/generate_test.go`
- Modify: `go.mod`, `go.sum` (adds `github.com/google/jsonschema-go v0.4.3`, used only in tests)

**Interfaces:**
- Consumes: Task 3 `modeltest` declarations.
- Produces: `GenerateJSON`, which `TestModel` (Task 9) uses for tool arguments and output.

`GenerateJSON` decodes the schema into `map[string]any` and walks it. It uses no schema library. The rules, applied in this order at each node:

| Schema | Generated value |
|--------|-----------------|
| boolean schema `true` | `null` |
| boolean schema `false` | error wrapping `ErrUnsatisfiable` |
| has `$ref` (`#/$defs/X` or `#/definitions/X`) | the value for the referenced schema, resolved against the root |
| has `const` | that value |
| has `enum` | its first value |
| has `anyOf` or `oneOf` | the value for the first branch |
| `type` is a list | the first type that is not `"null"`; `null` if every entry is `"null"` |
| `object` | an object with one generated value for each name in `required`; other properties are left out |
| `array` | `minItems` generated values (default 0) from `items` |
| `string` | `"a"`, repeated `minLength` times when `minLength` > 1 |
| `integer` | `minimum` if set; `exclusiveMinimum + 1` if set; otherwise `0` |
| `number` | `minimum` if set; `exclusiveMinimum + 1` if set; otherwise `0` |
| `boolean` | `false` |
| `null` | `null` |
| no `type` and none of the above | `null` |

The recursion depth is limited to 16. Deeper schemas return `ErrSchemaTooDeep`. Output is compact JSON, with object keys in sorted order.

- [ ] **Step 1: Write the failing tests**

- `TestGenerateBasicTypes` (table): `{"type":"string"}` gives `"a"`, `{"type":"integer"}` gives `0`, `{"type":"number"}` gives `0`, `{"type":"boolean"}` gives `false`, `{"type":"null"}` gives `null`.
- `TestGenerateObjectRequiredOnly`: `{"type":"object","properties":{"a":{"type":"string"},"b":{"type":"integer"}},"required":["a"]}` gives `{"a":"a"}`.
- `TestGenerateEnumAndConst`: `{"enum":["x","y"]}` gives `"x"`; `{"const":7}` gives `7`.
- `TestGenerateMinimums`: `{"type":"string","minLength":3}` gives `"aaa"`; `{"type":"integer","minimum":5}` gives `5`; `{"type":"integer","exclusiveMinimum":5}` gives `6`; `{"type":"array","items":{"type":"integer"},"minItems":2}` gives `[0,0]`.
- `TestGenerateRef`: both `#/$defs/X` and `#/definitions/X` resolve.
- `TestGenerateRecursiveRefFails`: `{"$ref":"#/$defs/N","$defs":{"N":{"type":"object","properties":{"next":{"$ref":"#/$defs/N"}},"required":["next"]}}}` gives an error wrapping `ErrSchemaTooDeep`.
- `TestGenerateNullableTypeList`: `{"type":["string","null"]}` gives `"a"`; `{"type":["null"]}` gives `null`.
- `TestGenerateAnyOf`: `{"anyOf":[{"type":"integer"},{"type":"string"}]}` gives `0`.
- `TestGenerateBooleanSchemas`: `true` gives `null`; `false` gives an error wrapping `ErrUnsatisfiable`.
- `TestGenerateValidatesAgainstInferredSchemas`:
  - Infer schemas with `jsonschema.For[T](nil)` for three test structs. Together they must cover a nested struct, a slice, a `map[string]int`, a pointer field, an `omitempty` field, and a `time.Time` field.
  - Generate a value for each schema and decode it into `any`.
  - Check that `schema.Resolve(nil)` followed by `resolved.Validate(value)` returns nil.

- [ ] **Step 2: Run them to see them fail**

Run: `go test ./model/modeltest/ -run Generate`
Expected: FAIL

- [ ] **Step 3: Implement `GenerateJSON`, then add the test dependency**

Run: `go get github.com/google/jsonschema-go@v0.4.3 && go mod tidy`

- [ ] **Step 4: Run them to see them pass**

Run: `go test -race ./model/modeltest/ -run Generate`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add model/modeltest/ go.mod go.sum
git commit -m "feat(modeltest): generate example json from a json schema"
```

---

### Task 9: FunctionModel, TestModel, and the Fake Loop

**Files:**
- Modify: `model/modeltest/function.go`, `model/modeltest/testmodel.go`, `model/modeltest/options.go`
- Test: `model/modeltest/function_test.go`, `model/modeltest/testmodel_test.go`, `model/modeltest/loop_test.go`, `model/modeltest/example_test.go`

**Interfaces:**
- Consumes: `GenerateJSON` (Task 8), `model.Collect` (Task 6), `message.Marshal` and `message.Unmarshal` (Task 5).
- Produces: the two test models that every wave-1 milestone uses.

**FunctionModel behavior:**
- `Request` returns `ctx.Err()` without calling the function if the context is already done.
- It then calls the function with `Info{Settings, Params}`.
- If the function left them empty, it fills `ModelName` with `Name()`, `ProviderName` with `Provider()`, and `Usage.Requests` with 1.
- `RequestStream` calls `Request`, then yields one `PartStart{Index: i, Part: part}` per part in order, then `Done{Response}`. On error it yields `(nil, err)` once. It returns as soon as `yield` returns false.

**TestModel behavior on each `Request`:**
1. If the last message is not a `message.Request`, return an error.
2. If no response in `msgs` contains a `ToolCallPart`, and the selected function tools are not empty, return one `ToolCallPart` per selected tool, in the order of `Params.FunctionTools`. Arguments come from `GenerateJSON(def.Parameters)`, and the finish reason is `FinishToolCall`.
3. Otherwise return the final answer:
   - In tool mode with output tools, return a `ToolCallPart` for `OutputTools[0]` with generated arguments (finish reason `FinishToolCall`).
   - In native or prompted mode with an `OutputObject`, return a `TextPart` holding `GenerateJSON(OutputObject.Schema)` (finish reason `FinishStop`).
   - Otherwise return a `TextPart` with the configured text (finish reason `FinishStop`).
4. Tool call IDs are `call_<n>`. Here `n` counts every `ToolCallPart` already in `msgs`, plus this part's position starting at 1. IDs are therefore unique within a run and the same on every run.
5. `Usage` is `{Requests: 1}`.

- [ ] **Step 1: Write the failing tests**

- `TestFunctionModelPassesInputs`: the function receives the exact messages, settings, and params.
- `TestFunctionModelFillsDefaults`: an empty response comes back with `ModelName` `function`, `ProviderName` `function`, and `Usage.Requests` 1. Values the function already set are kept.
- `TestFunctionModelCancelledContext`: with a cancelled context, the error is `context.Canceled` and the function is never called.
- `TestFunctionModelStreamEvents`: for a response with two parts, the events are `PartStart` with index 0, `PartStart` with index 1, then `Done`. `model.Collect` on the stream equals `Request`'s result.
- `TestFunctionModelStreamBreakEarly`: a consumer that breaks after the first event causes no panic, and the function is called once.
- `TestFunctionModelStreamError`: when the function returns `errBoom`, the stream yields exactly one `(nil, errBoom)`.
- `TestTestModelCallsToolsThenText`: given function tools `a` and `b` and a text output mode, the first request returns calls to `a` and `b` (IDs `call_1` and `call_2`). After a request that returns both results, it returns `TextPart{Content: DefaultText}`.
- `TestTestModelSelectedTools`: with `WithToolCalls("b")`, only `b` is called.
- `TestTestModelOutputTool`: in tool mode with output tool `final_result` (schema `{"type":"object","properties":{"answer":{"type":"string"}},"required":["answer"]}`), the final response is a call to `final_result` with arguments `{"answer":"a"}`.
- `TestTestModelNativeOutput`: in native mode with the same schema as an `OutputObject`, the final response is `TextPart{Content: `{"answer":"a"}`}`.
- `TestTestModelRejectsNonRequestLast`: when the history ends with a `Response`, it returns an error.
- `TestFakeLoopWithFunctionModel` (the M0 "done when" test):
  - The function returns `TextPart{Content: "Sunny in Amsterdam"}` when the last request holds a `ToolReturnPart`.
  - Otherwise it returns `ToolCallPart{ToolName: "get_weather", ToolCallID: "c1", Args: `{"city":"Amsterdam"}`}`.
  - The loop starts from `[Request{Parts: [UserPrompt("Weather?")]}]`. On each turn it appends the response. If the response has tool calls, it appends a `Request` with `ToolReturnPart{ToolName: "get_weather", ToolCallID: "c1", Content: `"sunny"`}`; otherwise it stops. It allows at most 5 turns.
  - Expected: 4 messages, with the final text as shown. `message.Unmarshal(message.Marshal(history))` is `reflect.DeepEqual` to the history.
- `TestFakeLoopWithTestModel`: the same loop with a `TestModel`, function tool `get_weather` (schema `{"type":"object","properties":{"city":{"type":"string"}},"required":["city"],"additionalProperties":false}`), and text output. The first response calls `get_weather` with `{"city":"a"}`, and the second is `DefaultText`.
- `ExampleNewFunctionModel`: a function model that always answers `hello`. Prints `hello`.

- [ ] **Step 2: Run them to see them fail**

Run: `go test ./model/modeltest/`
Expected: FAIL

- [ ] **Step 3: Implement both models and the options.**

- [ ] **Step 4: Run them to see them pass**

Run: `go test -race ./model/modeltest/`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add model/modeltest/
git commit -m "feat(modeltest): add function and test models"
```

---

### Task 10: Final Verification and Ready for Review

**Files:** none new. Updates the PR description only.

- [ ] **Step 1: Run every check**

Run: `make tidy-check test lint vuln`
Expected: exit code 0 in every module.

- [ ] **Step 2: Check that no placeholder panic remains**

Run: `grep -rn 'panic("gogent: not implemented")' --include=*.go . || echo none`
Expected: `none`. The only placeholders left are root-package bodies that return `ErrNotImplemented`.

- [ ] **Step 3: Check dependencies**

Run: `go list -m all` in the repository root.
Expected: `github.com/EmadMokhtar/gogent` and `github.com/google/jsonschema-go v0.4.3`, plus only what jsonschema-go itself requires.

Run: `grep -rEn -- '-(alpha|beta|rc)' --include=go.mod --include=go.sum . || echo none`
Expected: `none`.

- [ ] **Step 4: Check the doc comments**

Run: `go doc -all ./message | head -40` and `go doc -all ./model | head -40`
Expected: every exported name shows a comment. Lint already enforces this; this step is a visual check.

- [ ] **Step 5: Push and mark the PR ready**

Push the branch. When CI finishes, check `gh pr checks` and confirm that every job passed. Update the PR description with the "Done when" list from spec section 8.3, each item checked. Then run `gh pr ready`.
