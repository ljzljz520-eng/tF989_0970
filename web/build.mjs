import { accessSync } from "node:fs";
import { spawnSync } from "node:child_process";

const sources = ["assets/admin.js", "assets/labels.js"];
const required = [...sources, "assets/app.css", "views/admin.html", "views/labels.html"];

for (const file of required) {
  accessSync(new URL(file, import.meta.url));
}

for (const file of sources) {
  const result = spawnSync(process.execPath, ["--check", file], {
    cwd: new URL(".", import.meta.url),
    stdio: "inherit"
  });
  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

console.log("frontend sources verified");
