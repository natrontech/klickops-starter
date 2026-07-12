# Git Conventions

## When to commit

Commit automatically when a logical unit of work is complete - don't ask,
and don't batch unrelated changes. A commit is one coherent change that
could be reverted on its own:

- a new feature (migration + handler + tests + UI together)
- a bug fix
- a dependency or tooling change
- a docs-only change

Before every commit, the checks must pass:

```bash
make check && make test
```

If either fails, fix the cause first. Never commit around a failure and
never use `--no-verify`.

## Conventional Commits

All commit messages follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
```

### Types

- `feat` - new feature or capability
- `fix` - bug fix
- `refactor` - code change that neither fixes a bug nor adds a feature
- `docs` - documentation only
- `test` - adding or updating tests
- `chore` - tooling, dependencies, build config
- `ci` - GitHub Actions workflow changes
- `perf` - performance improvement

### Scopes

- `api` - backend handlers and routes
- `db` - database layer and migrations
- `storage` - object storage layer
- `ui` - frontend
- `deps` - dependency updates

Scope is optional when the change spans the whole repo.

### Rules

- Lowercase everything, imperative mood ("add" not "added"), no period at
  the end, description under 72 characters.
- Body is optional - use it to explain *why* when the diff doesn't say it.

### Examples

```
feat(api): add customer export endpoint
fix(ui): show upload error inline instead of failing silently
docs: explain the S3 binding variables
chore(deps): pin typescript to 6
```

## Never

- Force-push to main.
- Rewrite published history (`rebase`, `commit --amend` after push).
- Commit secrets, `.env` files, or generated build output.
