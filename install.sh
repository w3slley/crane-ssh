#!/bin/bash

target_bin_dir="$HOME/.local/bin"
target_dir="$HOME/.crane-ssh"
repo_url="https://github.com/w3slley/crane-ssh.git"

# Check if the directory exists
if [ -d "$target_dir" ]; then
    echo "Directory $target_dir already exists."
    read -p "Do you want to overwrite it? [y/N]: " choice
    case "$choice" in 
        y|Y ) 
            echo "Overwriting $target_dir..."
            rm -rf "$target_dir"
            ;;
        * ) 
            echo "Aborting."
            exit 1
            ;;
    esac
fi

git clone "$repo_url" "$target_dir"

cd "$target_dir"

if [ -d "$target_bin_dir" ]; then
    echo "Directory $target_bin_dir already exists."
else
  mkdir "$target_bin_dir"
fi

# Build app
go build crane-ssh.go
mv crane-ssh "$target_bin_dir"

current_shell=$(basename "$SHELL")

if [ "$current_shell" = "zsh" ]; then
    shell_rc="$HOME/.zshrc"
elif [ "$current_shell" = "bash" ]; then
    shell_rc="$HOME/.bashrc"
else
    echo "Unsupported or unknown shell: $current_shell. Only Bash and Zsh are supported."
    exit 1
fi

if ! grep -q 'export PATH="$HOME/.local/bin:$PATH"' "$shell_rc"; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$shell_rc"
    echo "$shell_rc updated successfully!"
else
    echo "PATH already includes ~/.local/bin."
fi

echo "Please restart your terminal or run: source $shell_rc"

