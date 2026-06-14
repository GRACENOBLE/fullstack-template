/**
 * PreToolUse hook — block dangerous Bash commands before execution.
 * Exit 2 to block + feed the message back to Claude so it can adjust.
 * Exit 0 to allow.
 */
const chunks = [];
for await (const chunk of process.stdin) chunks.push(chunk);
const input = JSON.parse(Buffer.concat(chunks).toString());

const command = (input.tool_input?.command ?? '').trim();

const BLOCKED = [
  [
    /git\s+push\s+(.*\s+)?--force|-f\b.*origin|origin.*-f\b/,
    'Force push is blocked. Open a PR against main instead of force-pushing.',
  ],
  [
    /git\s+push\s+\S+\s+main\b/,
    'Direct push to main is blocked. Create a feature branch and open a PR.',
  ],
  [
    /git\s+reset\s+--hard/,
    'git reset --hard is blocked. Stash your changes or create a backup branch first.',
  ],
  [
    /\brm\s+-[a-z]*r[a-z]*f[a-z]*\s+[\/~]/,
    'Recursive force-delete from root or home is blocked.',
  ],
  [
    /\bDROP\s+(TABLE|DATABASE|SCHEMA)\b/i,
    'DROP TABLE/DATABASE/SCHEMA is blocked. Write a migration file instead.',
  ],
  [
    />\s*\.env(\.|$)/,
    'Overwriting .env is blocked. Edit it directly with your editor.',
  ],
  [
    /git\s+commit\s+.*--no-verify/,
    '--no-verify is blocked. Fix the pre-commit hook failure instead of bypassing it.',
  ],
];

for (const [pattern, message] of BLOCKED) {
  if (pattern.test(command)) {
    process.stderr.write(`[guard] Blocked: ${message}\n`);
    process.exit(2);
  }
}

process.exit(0);
