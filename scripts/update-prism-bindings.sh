#!/usr/bin/env bash
# Refresh the checked-in Go bindings using the published Prism release.
set -euo pipefail

if [[ $# -gt 1 || ( $# -eq 1 && $1 != --check ) ]]; then
    printf 'Usage: %s [--check]\n' "${0##*/}" >&2
    exit 2
fi

version=2.4
sha256=e6b6b685f73d65c8ed0781458a909c002143c943382a4f8ca60e9bf234fa9b0a
release_url="https://github.com/polypheny/Polypheny-Prism-API/releases/download/${version}/go-api-files-${version}.tar.gz"
repo_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)
target="$repo_dir/prism"
work_dir=$(mktemp -d "${TMPDIR:-/tmp}/polypheny-prism-go.XXXXXX")
trap 'rm -rf "$work_dir"' EXIT

archive="$work_dir/prism.tar.gz"
curl --fail --silent --show-error --location --retry 3 --output "$archive" "$release_url"

if command -v sha256sum >/dev/null 2>&1; then
    printf '%s  %s\n' "$sha256" "$archive" | sha256sum -c -
else
    printf '%s  %s\n' "$sha256" "$archive" | shasum -a 256 -c -
fi

tar -xzf "$archive" -C "$work_dir"
generated="$work_dir/prism"
shopt -s nullglob
bindings=("$generated"/*.pb.go)
if [[ ${#bindings[@]} -eq 0 || ! -f "$generated/version.go" ]]; then
    printf 'The verified Prism archive does not contain the expected Go bindings.\n' >&2
    exit 1
fi
bindings+=("$generated/version.go")

if [[ ${1:-} == --check ]]; then
    status=0
    for file in "${bindings[@]}"; do
        if ! diff -u "$target/${file##*/}" "$file"; then
            status=1
        fi
    done
    for file in "$target"/*.pb.go; do
        if [[ ! -f "$generated/${file##*/}" ]]; then
            printf 'Unexpected generated binding: %s\n' "$file" >&2
            status=1
        fi
    done
    exit "$status"
fi

mkdir -p "$target"
cp "${bindings[@]}" "$target/"
for file in "$target"/*.pb.go; do
    if [[ ! -f "$generated/${file##*/}" ]]; then
        rm "$file"
    fi
done
