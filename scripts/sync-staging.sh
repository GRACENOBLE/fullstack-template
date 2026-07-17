#!/usr/bin/env bash
# scripts/sync-staging.sh — rebase staging onto main and force-push (with lease)
#
# Run automatically by .github/workflows/sync-staging.yml after every merge to
# main. Non-interactive by design — there is no confirmation prompt, unlike
# scripts/deploy-prod.sh which gates production deploys behind a manual step.
#
# Can also be run locally (git identity must already be configured).

set -euo pipefail

REMOTE="${REMOTE:-origin}"
BASE_BRANCH="main"
SYNC_BRANCH="staging"

echo "Fetching from ${REMOTE}..."
git fetch "$REMOTE"

if git show-ref --verify --quiet "refs/remotes/$REMOTE/$SYNC_BRANCH"; then
  REMOTE_SYNC_EXISTS=true
  if git show-ref --verify --quiet "refs/heads/$SYNC_BRANCH"; then
    git checkout --quiet "$SYNC_BRANCH"
  else
    git checkout --quiet -b "$SYNC_BRANCH" "$REMOTE/$SYNC_BRANCH"
  fi
  # Make sure the local branch starts from exactly what's on the remote before
  # rebasing, so we never rebase stale local commits onto main by accident.
  git reset --hard "$REMOTE/$SYNC_BRANCH"
else
  # Remote branch doesn't exist yet (e.g. first run of this template).
  REMOTE_SYNC_EXISTS=false
  git checkout --quiet -b "$SYNC_BRANCH" "$REMOTE/$BASE_BRANCH"
fi

echo "Rebasing ${SYNC_BRANCH} onto ${REMOTE}/${BASE_BRANCH}..."
if ! git rebase "$REMOTE/$BASE_BRANCH"; then
  echo "Rebase of ${SYNC_BRANCH} onto ${REMOTE}/${BASE_BRANCH} hit conflicts." >&2
  echo "This needs manual resolution — sync did not run:" >&2
  echo "  git fetch $REMOTE" >&2
  echo "  git checkout $SYNC_BRANCH && git reset --hard $REMOTE/$SYNC_BRANCH" >&2
  echo "  git rebase $REMOTE/$BASE_BRANCH   # resolve conflicts, then --continue" >&2
  echo "  git push --force-with-lease $REMOTE $SYNC_BRANCH" >&2
  git rebase --abort
  exit 1
fi

if [ "$REMOTE_SYNC_EXISTS" = true ]; then
  COMMITS="$(git log "$REMOTE/$SYNC_BRANCH..$SYNC_BRANCH" --oneline)"
  if [ -z "$COMMITS" ]; then
    echo "No new commits — ${SYNC_BRANCH} already matches ${BASE_BRANCH}."
    exit 0
  fi
  echo "Syncing ${SYNC_BRANCH}:"
  echo "$COMMITS"
else
  echo "Creating ${SYNC_BRANCH} on ${REMOTE} from ${BASE_BRANCH} (first run)."
fi

git push --force-with-lease "$REMOTE" "$SYNC_BRANCH"
echo "Synced: ${SYNC_BRANCH} pushed to ${REMOTE}."
