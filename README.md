# shell-zsh — A Unix Shell Built From Scratch in Go

A POSIX-inspired shell implemented from first principles — no readline library,
no shell-parsing package. Just raw syscalls, a hand-rolled parser, and signal
handling done the hard way.

[![asciicast](https://asciinema.org/a/HkXTQXoc9ky3e7jn.svg)](https://asciinema.org/a/HkXTQXoc9ky3e7jn)


## Why I built this

I wanted to actually understand what a shell does between you pressing Enter
and a process running — process creation, signal delivery, file descriptor
redirection, terminal control — instead of treating it as a black box. This
is the result of building that understanding one syscall at a time.

## Features

- **Command execution** — forks and execs external programs, resolves `$PATH`
- **Input/output redirection** — `<`, `>`, `>>` implemented via file descriptor manipulation
- **Tab autocompletion** — completes commands and file paths as you type
- **Signal handling** — Ctrl-C / Ctrl-Z routed correctly to foreground processes rather than killing the shell itself
- **Environment variables** — get, set, and export vars; expansion inside commands (`$VAR`)
- **Command history** — recall and re-run previous commands
- **Animated welcome message** — because a shell can have some personality

## How it works

A quick look under the hood:

1. **Read** — the shell reads a line from the terminal in raw mode, driving the tab-completion and history logic as you type
2. **Parse** — the input is tokenized and split into a command plus any redirection/argument metadata
3. **Set up redirection** — file descriptors are duped/rewired *before* exec, so the child process just sees a normal stdin/stdout
4. **Execute** — the shell forks and execs the target binary, tracking its PID as the foreground process
5. **Signal routing** — the shell installs its own signal handlers and forwards relevant signals to the foreground child, so Ctrl-C interrupts your program instead of the shell.

## What's next

- Pipes (`|`)
- Job control (`bg`, `fg`, `jobs`)

## Tech

Written in Go, using only the standard library — no external shell, parsing,
or readline packages.
