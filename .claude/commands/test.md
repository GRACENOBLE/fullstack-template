Run the full test suite and report results.

## Backend unit + integration tests
```bash
cd backend && make test
```

## Backend integration tests only (requires Docker running)
```bash
cd backend && make itest
```

Show full output for any failing test, including the test name, file, line number, and error message.

End with a summary:
- Total tests run
- Pass / Fail count
- Any tests skipped and why
