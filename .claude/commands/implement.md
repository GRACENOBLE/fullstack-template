Implement a feature following the Documentation-First workflow. This ensures every implementation is grounded in existing patterns and leaves the docs in a better state than it found them.

## Required from user before starting
- Feature description (what it does, which layer — backend / web / both)
- Any constraints or decisions already made

---

## Step 1 — Documentation check
Delegate to the `docs` agent:

> "Check existing documentation relevant to [feature]. Read _index.md for both apps,
> find the relevant topic docs, compare them against their source files, and report:
> which docs are current, which are stale, and what's missing."

Wait for the report before proceeding.

---

## Step 2 — Assess and decide
Based on the docs agent report:

| Situation | Action |
|---|---|
| Docs match code | Proceed to Step 3 |
| Docs are stale | Fix the stale docs first (delegate to `docs` agent), then Step 3 |
| No docs exist | Note it — new docs will be created in Step 4 |

Do not implement against stale documentation. Update docs first.

---

## Step 3 — Implement
Delegate to the appropriate agent (`backend`, `web`, or both).

Pass the relevant doc content as context so the agent implements to the correct patterns — not from training data.

Run `/project:check` at the end of this step. Do not proceed if it fails.

---

## Step 4 — Update documentation
Delegate to the `docs` agent:

> "Update documentation to reflect the implementation of [feature].
> Files changed: [list files]. Update last_verified dates.
> Create new topic doc if a new topic was introduced."

---

## Step 5 — Summary
Report:
- What was implemented
- Which docs were updated or created
- Output of `/project:check`
