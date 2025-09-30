#!/usr/bin/env bash

set -euo pipefail

OWNER="$(gh api user | jq -r '.login')"
REPO="topn"
TAP_REPO="homebrew-tap"  # set empty to skip tap cleanup
BREW_TAP_REPO="$(brew --repository "$OWNER"/homebrew-tap)"

echo "== Uninstall local brew artifacts =="
brew uninstall --force topn || true
brew cleanup -s topn || true
rm -rf ~/Library/Caches/Homebrew/topn* 2>/dev/null || true

echo "== Delete GitHub releases (and tags) =="
gh auth status >/dev/null
for tag in $(gh release list -R "$OWNER/$REPO" --limit 1000 | awk '{print $1}'); do
    echo "Deleting release $tag"
    gh release delete "$tag" -R "$OWNER/$REPO" --yes --cleanup-tag
done

echo "== Delete local tags =="
git tag -d "$(git --no-pager tag -l)"

echo "== Delete leftover remote tags =="
git ls-remote --tags "https://github.com/$OWNER/$REPO.git" | awk -F/ '{print $3}' | while read -r t; do
    [[ -z "$t" ]] && continue
    echo "Deleting remote tag: $t"
    git push "https://github.com/$OWNER/$REPO.git" :refs/tags/"$t" || true
done

echo "== Delete local tags =="
git -C "$BREW_TAP_REPO" tag -l | xargs -r -n1 git -C "$BREW_TAP_REPO" tag -d || true

if [[ -n "$TAP_REPO" ]]; then
    echo "== Remove formula from tap =="
    pushd "$BREW_TAP_REPO" >/dev/null
    git pull
    git rm -f Formula/topn.rb || true
    git commit -m "Remove topn formula (resetting)" || true
    git push || true
    popd >/dev/null
    brew untap "$OWNER/$TAP_REPO" || true
fi

echo "Done. Recreate a fresh tag (e.g. v0.1.0) and release when ready."
