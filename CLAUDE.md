# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose of this repo

Follow-along code for the Udemy course "Go (Golang) Backend Masterclass: REST APIs & Microservices" (https://www.udemy.com/course/learn-golang-like-google-engineers-do). The owner is an expert .NET/C# engineer learning Go, so this is a learning repo, not a production codebase.

## How to answer questions

The user will ask many questions: C# -> Go concept mappings, idiomatic Go, when/why to use pointers, returning a pointer vs. a value, etc.

- Search official Go sources and style guides (go.dev/doc, Effective Go, Go Code Review Comments, Go spec, Go blog, pkg.go.dev, Google Go Style Guide) and other relevant sites rather than answering from memory.
- Always cite the source of the answer (link or document/section name).
- Where useful, contrast with the C# equivalent.

## Layout

There is no `go.work` (it is gitignored) and no root `go.mod`. Every top-level directory is its own independent Go module with a single `main.go` and no third-party dependencies:

- `sectionNN_MM/` — scratch exercises named after the course section and lecture number (e.g. `section07_31`).
- `customLogger`, `safeMath`, `salesOrderProcessor`, `contactManagement` — larger exercise/mini-app projects from the end of a section.

Module paths follow `teqniqly.com/<snake_case_name>` (e.g. `teqniqly.com/section07_31`, `teqniqly.com/contact_management`), and the `go` directive is the installed toolchain version (currently 1.27.x). Directories are not related to each other, so changes rarely span modules.

Exercises are often written as a sequence of experiments in one `main()`, with earlier experiments left commented out. Preserve commented-out code unless asked to remove it.

## Commands

Because there is no root module, run Go commands from inside the module directory (or use `-C`):

```
go -C section07_31 run .              # run a module from the repo root

go -C section07_31 build ./...
go -C section07_31 vet ./...
gofmt -l .                            # list unformatted files (run from repo root; walks all modules)
gofmt -w section07_31/main.go         # format in place
```

There are currently no `_test.go` files and no linter config (no `.golangci.yml`). If tests are added, run a single one with `go -C <module> test -run '^TestName$' ./...`.

New exercise: create a directory, then `go mod init teqniqly.com/<name>` inside it.

## Environment notes

- Windows 11; PowerShell is the primary shell, Git Bash is also available.
- `.vscode/settings.json` is gitignored.
- Workflow so far: each section's work lands on a branch and is squash-merged to `main` via a PR (see `git log`).
