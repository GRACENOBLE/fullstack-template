#!/usr/bin/env bash
# scripts/deploy-prod.sh — rebase production onto main and force-push (with lease)
# Usage: ./scripts/deploy-prod.sh  (or via: make deploy-prod)
#
# Fast-forwards `production` to include everything on `main` by rebasing it
# onto origin/main, then force-pushes (with lease, never a bare --force) so
# the push fails instead of clobbering anyone else's work if the remote
# `production` moved since we last fetched.

set -euo pipefail

REMOTE="origin"
BASE_BRANCH="main"
DEPLOY_BRANCH="production"

RED='\033[1;31m'
YELLOW='\033[1;33m'
GREEN='\033[1;32m'
RESET='\033[0m'

if [ -n "$(git status --porcelain)" ]; then
  echo -e "${RED}Working tree is not clean. Commit, stash, or discard changes before deploying.${RESET}" >&2
  exit 1
fi

ORIGINAL_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
cleanup() {
  git checkout --quiet "$ORIGINAL_BRANCH" 2>/dev/null || true
}
trap cleanup EXIT

echo -e "${YELLOW}Fetching ${REMOTE}/${BASE_BRANCH} and ${REMOTE}/${DEPLOY_BRANCH}...${RESET}"
git fetch "$REMOTE" "$BASE_BRANCH" "$DEPLOY_BRANCH"

if git show-ref --verify --quiet "refs/heads/$DEPLOY_BRANCH"; then
  git checkout --quiet "$DEPLOY_BRANCH"
else
  git checkout --quiet -b "$DEPLOY_BRANCH" "$REMOTE/$DEPLOY_BRANCH"
fi
# Make sure the local branch starts from exactly what's on the remote before
# rebasing, so we never rebase stale local commits onto main by accident.
git reset --hard "$REMOTE/$DEPLOY_BRANCH"

echo -e "${YELLOW}Rebasing ${DEPLOY_BRANCH} onto ${REMOTE}/${BASE_BRANCH}...${RESET}"
if ! git rebase "$REMOTE/$BASE_BRANCH"; then
  echo -e "${RED}Rebase hit conflicts. Resolve them, then run:${RESET}" >&2
  echo "  git rebase --continue" >&2
  echo "  git push --force-with-lease $REMOTE $DEPLOY_BRANCH" >&2
  echo "Or abort with: git rebase --abort" >&2
  # Leave the branch mid-rebase for the user to resolve — do not clean up.
  trap - EXIT
  exit 1
fi

echo ""
echo -e "${YELLOW}About to force-push (with lease) ${DEPLOY_BRANCH} to ${REMOTE}. This will trigger a production deploy.${RESET}"
COMMITS="$(git log "$REMOTE/$DEPLOY_BRANCH..$DEPLOY_BRANCH" --oneline)"
if [ -z "$COMMITS" ]; then
  echo "No new commits — ${DEPLOY_BRANCH} already matches ${BASE_BRANCH}."
else
  echo "$COMMITS"
fi

if [ "${CONFIRM:-}" != "yes" ]; then
  read -r -p "Continue? [y/N] " reply
  case "$reply" in
    [yY][eE][sS]|[yY]) ;;
    *) echo "Aborted. No push was made."; exit 1 ;;
  esac
fi

git push --force-with-lease "$REMOTE" "$DEPLOY_BRANCH"
echo -e "${GREEN}Deployed: ${DEPLOY_BRANCH} pushed to ${REMOTE}.${RESET}"
