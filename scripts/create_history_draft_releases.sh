#!/bin/bash
set -euo pipefail

HISTORY_DIR="database/history"

if [ ! -d "$HISTORY_DIR" ]; then
    echo "Error: Directory $HISTORY_DIR not found."
    exit 1
fi

# Find and sort all monthly db files (e.g. 2019-07.db)
FILES=$(find "$HISTORY_DIR" -maxdepth 1 -name "????-??.db" | sort)
TOTAL=$(echo "$FILES" | wc -l | tr -d ' ')

echo "Found $TOTAL monthly history database files to compress and release."

COUNT=0
for file in $FILES; do
    COUNT=$((COUNT + 1))
    basename=$(basename "$file")
    month="${basename%.db}"
    tag="v${month}"
    zst_file="${file}.zst"

    echo "=========================================="
    echo "[$COUNT/$TOTAL] Processing: $month (Tag: $tag)"

    # Compress using zstd -19 if not already compressed
    if [ ! -f "$zst_file" ] || [ "$file" -nt "$zst_file" ]; then
        echo "[$COUNT/$TOTAL] Compressing $file -> $zst_file with zstd -19..."
        zstd -19 -f "$file" -o "$zst_file"
    else
        echo "[$COUNT/$TOTAL] Compressed file $zst_file already exists."
    fi

    # Check if draft release already exists
    if gh release view "$tag" >/dev/null 2>&1; then
        echo "[$COUNT/$TOTAL] Release $tag already exists. Skipping release creation."
        continue
    fi

    # Create draft release and attach .db.zst file
    echo "[$COUNT/$TOTAL] Creating draft release $tag and uploading $zst_file..."
    gh release create "$tag" "$zst_file" \
        --draft \
        --title "$tag" \
        --notes "Historical stock database for ${month} (compressed with zstd -19)"

    echo "[$COUNT/$TOTAL] Successfully created draft release $tag."
done

echo "=========================================="
echo "Done! All $TOTAL draft releases processed."
