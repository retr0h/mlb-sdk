mod? go '.just/remote/go.mod.just'
mod? just '.just/remote/just.mod.just'

# --- Fetch ---

# Fetch shared justfiles from osapi-justfiles
fetch:
    mkdir -p .just/remote
    curl -sSfL https://raw.githubusercontent.com/osapi-io/osapi-justfiles/refs/heads/main/go.mod.just -o .just/remote/go.mod.just
    curl -sSfL https://raw.githubusercontent.com/osapi-io/osapi-justfiles/refs/heads/main/go.just -o .just/remote/go.just
    curl -sSfL https://raw.githubusercontent.com/osapi-io/osapi-justfiles/refs/heads/main/just.mod.just -o .just/remote/just.mod.just
    curl -sSfL https://raw.githubusercontent.com/osapi-io/osapi-justfiles/refs/heads/main/just.just -o .just/remote/just.just

# --- Top-level orchestration ---

# Install all dependencies
deps:
    just go::deps
    just go::mod

# Run all tests
test:
    just just::fmt-check
    just go::test

# Format, lint before committing. `go::vet` runs golangci-lint per the

# osapi-justfiles convention — it's the lint step despite the name.
ready:
    just just::fmt
    just go::fmt
    just go::vet

# Run code generation (oapi-codegen)
generate:
    just go::generate
