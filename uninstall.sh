target_bin_dir="$HOME/.local/bin"
target_repo_dir="$HOME/.crane-ssh"

echo "Removing executable in $target_bin_dir/crane-ssh"
rm -f "$target_bin_dir"/crane-ssh
echo "File $target_bin_dir/crane-ssh removed!"

echo "Removing executable in $target_bin_dir/crane-ssh"
rm -rf "$target_repo_dir"
echo "Folder $target_repo_dir removed!"

echo "crane-ssh was removed from your system!"
