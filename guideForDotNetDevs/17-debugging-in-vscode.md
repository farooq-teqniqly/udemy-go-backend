# 17 - Stepping through Go code in VS Code

**C#:** in Visual Studio you press F9 for a breakpoint, F5 to start, and F10/F11/Shift+F11 to step. The
debugger is built in and attaches to the CLR. Locals, Watch and Call Stack windows show state.

**Go:** the same keys work in VS Code, but the debugger is a separate tool, Delve (`dlv`), driven by the Go
extension. Examples use `templateExample` in `section10_47`.

## 1. Setup

- Install the Go extension. It may prompt you to install Delve, or run **Go: Install/Update Tools** from the
  Command Palette. `go install github.com/go-delve/delve/cmd/dlv@latest` also works.
- The extension looks for tools in `${GOPATH}/bin`, `${GOBIN}` and `PATH`. **Go: Locate Configured Go Tools**
  shows where it looks.

## 2. Breakpoint, start, step

1. Click the margin left of a line number, or press **F9**. For example, break on
   `err = template.Execute(&output, email)`.
2. Open the `package main` file and press **F5** (Run and Debug).
3. Step with the keys below.

| Action | Windows key |
|---|---|
| Continue | F5 |
| Step Over | F10 |
| Step Into | F11 |
| Step Out | Shift+F11 |
| Restart | Ctrl+Shift+F5 |
| Stop | Shift+F5 |

Visual Studio uses F9, F10, F11 and Shift+F11 the same way. Its F5, Shift+F5 and Ctrl+Shift+F5 were not on
the page checked, so those three are **Unverified**.

## 3. `launch.json`

If F5 cannot find the program, add a configuration to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [{
    "name": "Launch Go",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}/section10_47"
  }]
}
```

`mode: "auto"` picks debug or test from the active editor. `program` is the package or file to run, as an
absolute or workspace-relative path.

This repo has one module per top-level folder, no root `go.mod` and no `go.work`. That `program` must
therefore point at the module folder is my inference. The wiki does not discuss multi-module layouts, so it
is **Unverified**. A "Debug" CodeLens above `main` was also not confirmed; the wiki mentions only a
"debug test" CodeLens.

## 4. Panels

- **Variables** shows locals and lets you edit values. The Go extension hides globals unless you set
  `showGlobalVariables`. Expand `email` to see `Attachments`.
- **Watch** re-evaluates your expressions every time the debugger pauses, such as `email.Subject`.
- **Call Stack** lists goroutines, not threads. The current goroutine is marked `*`. This is the visible
  difference from the CLR debugger.
- **Debug Console** is a REPL. In Go it accepts Delve expressions and `call <func_expr>`.

## 5. Conditional breakpoints and logpoints

- **Conditional breakpoint:** right-click the margin, or run **Debug: Add Conditional Breakpoint...**. It
  takes an expression or a hit count. The Go wiki lists the hit-count operators `>`, `>=`, `<`, `<=`, `==`,
  `!=` and `% n`.
- **Logpoint:** run **Debug: Add Logpoint...**. It writes to the Debug Console without pausing, with
  expressions in braces, for example `recipient is {email.Recipient}`. Visual Studio calls these tracepoints.
  Unlike `fmt.Println`, a logpoint needs no code edit and no recompile.

## 6. Gaps against Visual Studio

- I found no Go equivalent of the Just My Code toggle in the pages checked.
- The wiki lists no Windows-specific debugging caveats. `asRoot` applies to Linux and Mac only.
- The Delve README says nothing about VS Code or Windows; it links to its installation docs.

## Sources

- Debugging Go code using VS Code (vscode-go wiki): <https://github.com/golang/vscode-go/wiki/debugging>
- VS Code debugging (breakpoints, logpoints, step keys): <https://code.visualstudio.com/docs/editor/debugging>
- Delve: <https://github.com/go-delve/delve>
- Navigate through code with the Visual Studio debugger: <https://learn.microsoft.com/en-us/visualstudio/debugger/navigating-through-code-with-the-debugger>

Fetched 2026-11-15 via a fetch tool that summarises pages. **Verified:** Delve install routes and tool lookup
path, F9 breakpoints, the step keys in the table, conditional breakpoints and hit-count operators, logpoint
syntax, `mode: "auto"`, the Variables, Watch, Call Stack and Debug Console behavior, and the Visual Studio
F9, F10, F11, Shift+F11 and Ctrl+F10 keys. **Unverified:** the multi-module `program` path advice, the
"Debug" CodeLens, Visual Studio's F5, Shift+F5 and Ctrl+Shift+F5, and the absence of Windows caveats and of a
Just My Code equivalent (absence of a mention only). The `launch.json` was not run in this repo.
