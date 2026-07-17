# scripts/deploy-prod.ps1 — rebase production onto main and force-push (with lease)
# Usage: .\scripts\deploy-prod.ps1  (or via: make deploy-prod)
#
# Fast-forwards `production` to include everything on `main` by rebasing it
# onto origin/main, then force-pushes (with lease, never a bare --force) so
# the push fails instead of clobbering anyone else's work if the remote
# `production` moved since we last fetched.

$ErrorActionPreference = "Stop"

$Remote = "origin"
$BaseBranch = "main"
$DeployBranch = "production"

$status = git status --porcelain
if ($status) {
    Write-Host "Working tree is not clean. Commit, stash, or discard changes before deploying." -ForegroundColor Red
    exit 1
}

$originalBranch = git rev-parse --abbrev-ref HEAD
$cleanExit = $false

try {
    Write-Host "Fetching from $Remote..." -ForegroundColor Yellow
    git fetch $Remote
    if ($LASTEXITCODE -ne 0) { throw "git fetch failed" }

    git show-ref --verify --quiet "refs/remotes/$Remote/$DeployBranch"
    $deployBranchExists = ($LASTEXITCODE -eq 0)
    if ($deployBranchExists) {
        git show-ref --verify --quiet "refs/heads/$DeployBranch"
        if ($LASTEXITCODE -eq 0) {
            git checkout --quiet $DeployBranch
        } else {
            git checkout --quiet -b $DeployBranch "$Remote/$DeployBranch"
        }
        if ($LASTEXITCODE -ne 0) { throw "git checkout failed" }

        # Make sure the local branch starts from exactly what's on the remote
        # before rebasing, so we never rebase stale local commits onto main.
        git reset --hard "$Remote/$DeployBranch"
        if ($LASTEXITCODE -ne 0) { throw "git reset failed" }
    } else {
        # Remote branch doesn't exist yet (e.g. first deployment).
        git checkout --quiet -b $DeployBranch "$Remote/$BaseBranch"
        if ($LASTEXITCODE -ne 0) { throw "git checkout failed" }
    }

    Write-Host "Rebasing $DeployBranch onto $Remote/$BaseBranch..." -ForegroundColor Yellow
    git rebase "$Remote/$BaseBranch"
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Rebase hit conflicts. Resolve them, then run:" -ForegroundColor Red
        Write-Host "  git rebase --continue"
        Write-Host "  git push --force-with-lease $Remote $DeployBranch"
        Write-Host "Or abort with: git rebase --abort"
        # Leave the branch mid-rebase for the user to resolve — do not clean up.
        exit 1
    }

    Write-Host ""
    Write-Host "About to force-push (with lease) $DeployBranch to $Remote. This will trigger a production deploy." -ForegroundColor Yellow
    if ($deployBranchExists) {
        $commits = git log "$Remote/$DeployBranch..$DeployBranch" --oneline
        if (-not $commits) {
            Write-Host "No new commits -- $DeployBranch already matches $BaseBranch."
        } else {
            $commits | ForEach-Object { Write-Host $_ }
        }
    } else {
        Write-Host "Creating $DeployBranch on $Remote from $BaseBranch (first deployment)."
    }

    if ($env:CONFIRM -ne "yes") {
        $reply = Read-Host "Continue? [y/N]"
        if ($reply -notmatch '^(y|yes)$') {
            Write-Host "Aborted. No push was made."
            exit 1
        }
    }

    git push --force-with-lease $Remote $DeployBranch
    if ($LASTEXITCODE -ne 0) { throw "git push failed" }

    Write-Host "Deployed: $DeployBranch pushed to $Remote." -ForegroundColor Green
    $cleanExit = $true
}
finally {
    if ($cleanExit) {
        git checkout --quiet $originalBranch 2>$null
    }
}
