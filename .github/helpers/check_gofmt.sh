#!/bin/bash

# Check for unformatted files using gofmt
unformatted_files=$(gofmt -l .)

# If there are unformatted files, print an error message and exit with a non-zero status
if [ -n "$unformatted_files" ]; then
    echo "The following files are not formatted:"
    echo "$unformatted_files"
    exit 1
fi