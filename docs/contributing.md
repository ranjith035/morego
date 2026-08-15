# Contributing Guidelines

We welcome contributions to the Mobile Automation Platform! This guide explains how to set up your environment, follow our development process, and submit pull requests.

---

## 1. Development Environment Setup

To contribute to the Go Core Engine and the documentation, you will need:

*   **Go:** Version 1.22 or higher.
*   **Make:** GNU Make to run build targets.
*   **Docker:** Recommended for verifying container builds.
*   **Protocol Buffers Compiler (`protoc`):** Required if you are modifying `.proto` files in Phase 2+.

### Clone the Repository
```bash
git clone https://github.com/mobile-automation/mobile-framework.git
cd mobile-framework
```

### Initialize Workspace
Run `make init` to synchronize workspace packages:
```bash
make init
```

---

## 2. Branching & PR Workflow

We follow a structured pull request lifecycle:

1.  **Fork and Branch:** Create a branch from `develop` with a descriptive name:
    *   `feature/feature-name` for new features.
    *   `bugfix/bug-description` for bug fixes.
    *   `docs/doc-updates` for documentation.
2.  **Commit Standards:** Follow semantic commit messages:
    *   `feat: add auto-wait mechanism`
    *   `fix: resolve memory leak in session manager`
    *   `docs: update architectural diagrams`
3.  **Run Checks Locally:** Before pushing, ensure your changes compile, pass formatting rules, and pass tests:
    *   Format: `make fmt`
    *   Vet: `make lint`
    *   Test: `make test`
4.  **Create a Pull Request:** Submit your PR targeting the `develop` branch. Fill out the pull request template completely.
5.  **Review Process:** At least one core maintainer must approve the PR. The GitHub CI workflow must pass cleanly.

---

## 3. Coding Guidelines & Standards

*   **Format:** All Go code must be formatted using standard `gofmt` (or `goimports`).
*   **Error Handling:** Never ignore errors. Always wrap errors with context or return them to the caller.
*   **Panics:** Never call `panic()` or exit abruptly in library modules. Only CLI entrypoints under `cmd` are allowed to perform fatal exits under extreme startup failure conditions.
*   **Concurrency:** Use standard channel patterns or the `sync` package. Avoid spawning untracked goroutines; always tie goroutine lifecycles to a `context.Context`.
*   **Tests:** Every new feature requires unit tests. Use table-driven tests where appropriate.
