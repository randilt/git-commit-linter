#!/bin/bash

# Check if inside a git directory
if [ ! -d ".git" ]; then
  echo "Not a git directory. Please run this inside a Git repository."
  exit 1
fi

# Check if commit-msg hook already exists
HOOK_PATH=".git/hooks/commit-msg"
if [ -f "$HOOK_PATH" ]; then
  # Ask the user if they want to overwrite the existing commit-msg hook
  echo "commit-msg hook already exists."
  echo "Options: "
  echo "  y/Y  - Overwrite the commit-msg hook"
  echo "  n/N  - Cancel the operation"
  echo "  v/V  - View the current commit-msg hook content"
  read -p "Select an option (y/n/v): " choice

  case "$choice" in
    y|Y)
      echo "Overwriting commit-msg hook..."
      ;;
    n|N)
      echo "Operation canceled. No changes made."
      exit 0
      ;;
    v|V)
      echo "Displaying current commit-msg hook content:"
      cat "$HOOK_PATH"
      read -p "Do you still want to overwrite? (y/n): " overwrite_choice
      case "$overwrite_choice" in
        y|Y)
          echo "Overwriting commit-msg hook..."
          ;;
        n|N)
          echo "Operation canceled. No changes made."
          exit 0
          ;;
        *)
          echo "Invalid choice. Operation canceled."
          exit 1
          ;;
      esac
      ;;
    *)
      echo "Invalid choice. Operation canceled."
      exit 1
      ;;
  esac
else
  # Rename commit-msg.sample to commit-msg if it exists
  if [ -f ".git/hooks/commit-msg.sample" ]; then
    mv .git/hooks/commit-msg.sample .git/hooks/commit-msg
    echo "Renamed commit-msg.sample to commit-msg"
  else
    echo "No commit-msg.sample file found."
    exit 1
  fi
fi

# Add the script to commit-msg hook
echo -e "#!/bin/sh\ncommit_msg_file=\"\$1\"\n\n# Run the linter with the commit message file\ngit-commit-linter lint-file \"\$commit_msg_file\" || exit 1" > "$HOOK_PATH"

# Make the hook executable
chmod +x "$HOOK_PATH"

echo "commit-msg hook created successfully!"
