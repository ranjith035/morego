# Engineering Guidelines

This document establishes the technical rules, conventions, and processes for developers contributing code and APIs to the Mobile Automation Platform.

---

## 1. General Philosophy

*   **API Stability:** Treat every public API as a 10-year compatibility contract. Changing a public signature requires deprecation cycles and major version releases.
*   **Composition Over Inheritance:** Prefer simple interface composition over complex type hierarchy tree structures.
*   **No Global State:** Global variables and global singleton states make concurrent execution (e.g. running multiple parallel mobile sessions) impossible or highly prone to race conditions. Instantiation and state must be encapsulated in local context structs.
*   **Zero-Copy & Zero-Allocation:** Write hot execution paths (event loops, buffer readers, locator matchers) to avoid heap allocations. Re-use byte buffers and slice capacities where possible.

---

## 2. Go Coding Standards

We follow idiomatic Go practices as outlined in *Effective Go* and the Go Code Review Comments.

### 2.1. Project & Package Layout
*   Avoid large, monolithic files. Split logical concerns into distinct files within a package.
*   Keep packages small and focused. If a package contains multiple unrelated structures, split it.
*   Internal utility functions should stay in `internal/` packages to prevent exposing implementation details to SDK clients.

### 2.2. Error Handling
*   **Always return errors:** Never discard errors. Handled errors must be logged, wrapped, or returned up the call stack.
*   **Wrap errors with context:** Use `%w` to wrap underlying errors.
    ```go
    if err != nil {
        return fmt.Errorf("failed to locate element %s: %w", locator, err)
    }
    ```
*   **No Panics:** Libraries must not call `panic()`. Recover from potential panics in critical event loops to ensure session crashes do not bring down the entire execution engine.

### 2.3. Concurrency
*   Always associate goroutines with a lifetime boundary using `context.Context`.
*   Ensure all goroutines exit cleanly when the context is cancelled to prevent goroutine/memory leaks.
*   Run the Go race detector (`go test -race`) during local verification and CI.

---

## 3. gRPC API Guidelines

Since the Go Core communicates with all SDKs via gRPC, protocol definitions are the ultimate source of truth.

### 3.1. Versioning
*   Protocol buffer files are located in `/proto/`.
*   All proto APIs must be versioned. For example, use package name `automation.session.v1`.
*   Backward-compatible changes (adding fields, adding RPC methods) are allowed.
*   Breaking changes (renaming fields, changing field types, removing fields/methods) are strictly forbidden within the same version. Create a `v2` package if a breaking change is required.

### 3.2. Error Delivery
*   Use standard gRPC error codes (e.g., `InvalidArgument`, `NotFound`, `DeadlineExceeded`) for protocol-level issues.
*   Embed detailed structural error information (e.g., failed assertion details, element screenshot bytes upon failure) inside the gRPC `status.Status` detail payload instead of parsing raw error strings.

---

## 4. Release Process & Versioning

We adhere strictly to Semantic Versioning (SemVer) 2.0.0.

### 4.1. Version Format
`MAJOR.MINOR.PATCH`
*   **MAJOR:** Incremented when backwards-incompatible API changes are introduced.
*   **MINOR:** Incremented when backward-compatible features are added (e.g., a new locator strategy `GetByRole` or a new SDK language binding).
*   **PATCH:** Incremented for backwards-compatible bug fixes.

### 4.2. Release Cycle
1.  **Release Branches:** Releases are tagged from the `main` branch. Development happens on `develop`.
2.  **Release Candidate:** When preparing a release, a release candidate branch is branched from `develop` (e.g., `release/v1.2.0-rc1`). Only bug fixes are permitted on this branch.
3.  **Final Tag:** Once validated, the release branch is merged into `main`, tagged (e.g., `v1.2.0`), and merged back into `develop`.
4.  **Changelog:** A changelog is automatically compiled from commit messages and PR templates.

---

## 5. Architectural Decision Records (ADRs)

All significant structural design decisions must be documented using ADRs.

*   ADRs are stored in `/docs/adr/`.
*   Filename format: `NNNN-short-descriptive-title.md` (e.g., `0002-use-grpc-over-http.md`).
*   Each ADR follows a standard template containing:
    *   **Title & Index**
    *   **Status** (Proposed, Accepted, Rejected, Deprecated)
    *   **Context** (The problem and background context)
    *   **Decision** (What we decided to do and why)
    *   **Consequences** (Tradeoffs, new requirements, or risks created by the decision)
