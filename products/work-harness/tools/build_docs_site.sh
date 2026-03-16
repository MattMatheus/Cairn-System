#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
out_dir="${1:-$root_dir/.site}"

rm -rf "$out_dir"
mkdir -p "$out_dir"

if [[ -d "$root_dir/website/cairn-homepage" ]]; then
  cp -R "$root_dir/website/cairn-homepage/." "$out_dir/"
fi

docs_root="$out_dir/docs"
mkdir -p "$docs_root"
mkdir -p "$docs_root/products/work-harness"

cp "$root_dir/README.md" "$docs_root/README.md"
cp "$root_dir/products/work-harness/AGENTS.md" "$docs_root/products/work-harness/AGENTS.md"
cp "$root_dir/products/work-harness/HUMANS.md" "$docs_root/products/work-harness/HUMANS.md"
cp "$root_dir/products/work-harness/DEVELOPMENT_CYCLE.md" "$docs_root/products/work-harness/DEVELOPMENT_CYCLE.md"
cp -R "$root_dir/docs/decisions" "$docs_root/decisions"
if [[ -d "$root_dir/knowledge-base" ]]; then
  cp -R "$root_dir/knowledge-base" "$docs_root/knowledge-base"
fi

cat > "$docs_root/index.html" <<'HTML'
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Cairn Docs</title>
    <meta
      name="description"
      content="Cairn documentation published from repository markdown."
    />
    <link rel="canonical" href="https://cairn.teamorchestrator.com/docs/" />
    <meta property="og:type" content="website" />
    <meta property="og:site_name" content="Cairn Docs" />
    <meta property="og:title" content="Cairn Documentation" />
    <meta
      property="og:description"
      content="Canonical docs source is repository markdown. This site is generated."
    />
    <meta property="og:url" content="https://cairn.teamorchestrator.com/docs/" />
    <meta name="twitter:card" content="summary" />
    <meta name="twitter:title" content="Cairn Documentation" />
    <meta
      name="twitter:description"
      content="Canonical docs source is repository markdown. This site is generated."
    />
    <style>
      body {
        font-family: "Avenir Next", "Segoe UI", Tahoma, sans-serif;
        margin: 0;
        background: radial-gradient(circle at 75% 12%, #ffe8bf 0%, #f7f4ed 55%);
        color: #1f2a37;
      }
      main {
        max-width: 920px;
        margin: 0 auto;
        padding: 48px 20px;
      }
      h1 {
        margin: 0 0 10px;
        font-size: clamp(34px, 6vw, 56px);
      }
      p {
        color: #54606e;
      }
      .grid {
        margin-top: 16px;
        display: grid;
        gap: 12px;
        grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      }
      .card {
        background: #fffaf0;
        border: 1px solid #d9d0bf;
        border-radius: 12px;
        padding: 12px 14px;
      }
      a {
        color: #b65a2a;
        font-weight: 600;
        text-decoration: none;
      }
      a:hover {
        text-decoration: underline;
      }
      .note {
        margin-top: 20px;
        border-left: 4px solid #e1b08e;
        padding: 10px 14px;
        background: #fffaf0;
      }
    </style>
  </head>
  <body>
    <main>
      <h1>Cairn Documentation</h1>
      <p>Published docs site for Cairn. Canonical source remains markdown in git.</p>
      <div class="grid">
        <div class="card"><a href="./README.md">Repository README</a></div>
        <div class="card"><a href="./products/work-harness/HUMANS.md">work harness Operator Guide</a></div>
        <div class="card"><a href="./products/work-harness/DEVELOPMENT_CYCLE.md">work harness Development Cycle</a></div>
        <div class="card"><a href="./products/work-harness/AGENTS.md">work harness Agent Guide</a></div>
        <div class="card"><a href="./decisions/">Decision Records (ADRs)</a></div>
        <div class="card"><a href="./decisions/README.md">Decision Index</a></div>
      </div>
      <p class="note">
        Source-of-truth policy: update docs in this repository, then publish through CI.
      </p>
    </main>
  </body>
</html>
HTML

echo "Built docs site at $out_dir"
