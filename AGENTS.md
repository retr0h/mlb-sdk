# AGENTS.md

Test: `just test` | Before committing: `just ready`

Read [docs/development.md](docs/development.md) first. It covers setup, the
layout, regenerating the client, the eight-step recipe for adding an endpoint,
the OpenAPI and public surface authoring rules, the table-driven testing rules,
the 100% coverage gate and commit messages.
[docs/contributing.md](docs/contributing.md) covers branches and pull requests.
All of it applies to agents exactly as it applies to people. This file carries
only what is specific to agents.

## Running tools

Invoke tools through `mise`, not from your path:

```bash
mise exec -- just test
```

`mise` is active in a person's shell and supplies the versions `.mise.toml`
declares. An agent's shell has no activation, so a bare `just` resolves to
whatever is installed globally, usually an older version.

The symptom is a check that fails here and passes in continuous integration, on
a file nobody edited. When that happens, establish which version ran before
treating the failure as real.

## Porting endpoints

To port endpoints from
[toddrob99/MLB-StatsAPI](https://github.com/toddrob99/MLB-StatsAPI), use the
[`ingest-from-toddrob99`](.claude/skills/ingest-from-toddrob99/SKILL.md) skill
rather than improvising the translation. It carries the Python to OpenAPI
mapping and a manifest of what has already been ported. Claude Code loads it on
its own. Other agents should read `SKILL.md` and follow it.

## Commit trailer

When committing via Claude Code, end the message with:

```
🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```
