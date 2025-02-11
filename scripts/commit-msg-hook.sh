#!/bin/sh
commit_msg_file="$1"

# Run the linter with the commit message file
git-commit-linter lint-file "$commit_msg_file" || exit 1