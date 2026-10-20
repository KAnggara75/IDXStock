#!/usr/bin/env bash
set -euo pipefail

# Script to compress database/master_stocks.db and upload/replace the latest GitHub release asset.
# Run manually: ./scripts/archive_and_release.sh

DB_FILE="database/master_stocks.db"
ZST_FILE="database/master_stocks.db.zst"

if [ ! -f "$DB_FILE" ]; then
    echo "Error: Database file '$DB_FILE' does not exist." >&2
    exit 1
fi

echo "==> Compressing $DB_FILE to $ZST_FILE using zstd -19..."
zstd -19 -T0 --keep -f "$DB_FILE" -o "$ZST_FILE"

echo "==> Detecting latest release tag from GitHub..."
LATEST_TAG=$(gh release view --json tagName -q .tagName)

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not retrieve latest release tag." >&2
    exit 1
fi

echo "==> Replacing asset in latest release '$LATEST_TAG'..."
gh release upload "$LATEST_TAG" "$ZST_FILE" --clobber

echo "==> Successfully compressed and replaced $ZST_FILE in release $LATEST_TAG!"
