#!/usr/bin/env node
const { execFileSync } = require("child_process");
const path = require("path");

const ext = process.platform === "win32" ? ".exe" : "";
const bin = path.join(__dirname, "..", "bin", "xfchat_cli" + ext);

try {
  execFileSync(bin, process.argv.slice(2), {
    stdio: "inherit",
    env: {
      ...process.env,
      XFCHAT_CLI_SKILLS_DIR: path.join(__dirname, "..", "skills"),
    },
  });
} catch (e) {
  process.exit(e.status || 1);
}
