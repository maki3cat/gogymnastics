#!/bin/bash
set -e

# Build the Go program
go build -o main main.go

# List global symbols in the binary, highlight 'global' in color
echo -e "Global symbols containing '\033[31mexample\033[0m':"
go tool nm main | GREP_COLOR='01;31' grep --color=always -E 'example'

echo -e " \033[31mB\033[0m: BSS (Uninitialized Data Segment)global or package-level variables that are declared but not explicitly initialized."
echo -e " \033[31mD\033[0m: Data Segment (Initialized Globals)global or package-level variables that are declared and explicitly initialized."

# # Show the size of the binary
# echo "Binary size:"
# ls -lh main

# Run the program to show addresses, making the addr green
echo "Program output:"
echo -e "\033[32mthe addresses are changing because of ASLR\033[0m"
./main

# Clean up
rm -f main
