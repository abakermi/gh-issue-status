# gh-issue-status

A GitHub CLI extension for managing issue status in project boards.

## Installation

```bash
gh extension install abdelhak/gh-issue-status
```


## Usage

Change issue status interactively:
```bash
gh issue-status change
```

Change specific issue status:
```bash
gh issue-status change 123 --status in-progress --project "My Project"
```

## Commands

- `gh issue-status change [issue-number]` - Change issue status in project board

## Flags

- `--status, -s` - New status (todo, in-progress, done)
- `--project, -p` - Project name or number
- `--repo, -R` - Repository (OWNER/REPO format)
- `--json, -j` - Output in JSON format

## Examples

```bash
# Interactive mode
gh issue-status change

# Direct mode
gh issue-status change 42 --status done --project "Sprint 1"

# With specific repo
gh issue-status change 42 --repo owner/repo --status in-progress
```
