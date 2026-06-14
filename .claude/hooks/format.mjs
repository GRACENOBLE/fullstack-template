/**
 * PostToolUse hook — auto-format after Edit, Write, or MultiEdit.
 * Runs gofmt on Go files, Prettier on TS/JS/CSS files.
 * Exit 0 always — format failures are non-blocking.
 */
import { execSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { extname, resolve, join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const projectDir = resolve(dirname(fileURLToPath(import.meta.url)), '..', '..');

const chunks = [];
for await (const chunk of process.stdin) chunks.push(chunk);
const input = JSON.parse(Buffer.concat(chunks).toString());

const filePath = input.tool_input?.file_path;
if (!filePath) process.exit(0);

const abs = resolve(filePath);
const ext = extname(abs);

try {
  if (ext === '.go') {
    execSync(`gofmt -w "${abs}"`, { stdio: 'pipe' });
  } else if (['.ts', '.tsx', '.js', '.jsx', '.css'].includes(ext)) {
    const prettier = join(projectDir, 'frontend', 'node_modules', '.bin', 'prettier');
    if (existsSync(prettier)) {
      execSync(`"${prettier}" --write "${abs}"`, { stdio: 'pipe' });
    }
  }
} catch {
  // Non-blocking — formatter not installed or file has syntax errors Claude must fix
}

process.exit(0);
