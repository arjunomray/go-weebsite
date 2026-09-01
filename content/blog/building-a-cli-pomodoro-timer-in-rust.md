---
title: Building a CLI Pomodoro timer in Rust
date: 2026-08-14
tags: rust, cli, learning
summary: Building and shipping a minimal CLI Pomodoro timer with Rust and publishing it to crates.io.
stage: published
---
I wanted a small, self-contained project to actually learn Rust with — not a tutorial, something I'd use. A Pomodoro timer felt right: small scope, clear behavior, easy to know when it's "done."

The result is **rustime** — a CLI Pomodoro timer with configurable session and break durations, named sessions, and a minimal output mode. Published it to crates.io so I could install it the same way I'd install any other CLI tool.

### What it does

```bash
rustime -s 25 -b 5 -t "deep work"
```

Configurable session time, break time, session titles, a `--list-sessions` flag to see history, and a `--minimal-version` for quieter output.

### What I actually learned

* **Rust's ownership model** is unforgiving in exactly the way people warn you about — I fought the borrow checker more than I expected on something this small.
* **Clap** (the CLI arg parser) made the interface trivial compared to hand-rolling flag parsing in Node.
* **Publishing to crates.io** for the first time — different ceremony than npm, but not harder.

### Next up

Moving the terminal UI to something more structured (probably [ratatui](https://github.com/ratatui/ratatui)) instead of raw stdout, and maybe Discord rich presence for session status.

[→ source on GitHub](https://github.com/arjunomray/rustime)
