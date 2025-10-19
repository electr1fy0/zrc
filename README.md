# zrc
A command-line tool to manage your shell configuration (`.zshrc`, `.bashrc`, etc.) with ease.

## Install



### From source
```bash
git clone https://github.com/your-username/zrc.git
cd zrc
go build -o zrc

# Move to PATH
sudo mv zrc /usr/local/bin/
```

## Usage

`zrc` helps you manage aliases and your PATH by modifying your shell's config file (e.g., `.zshrc`). On first run, it will prompt you to specify which config file to use.

### Add a raw line
```bash
zrc add 'export CUSTOM_VAR="hello"'
```

### Manage Aliases
```bash
# Add an alias
zrc alias add gs "git status"

# Remove an alias
zrc alias remove gs

# List all aliases
zrc list aliases
```

### Manage PATH
```bash
# Add a directory to your PATH
zrc path add ~/.local/bin

# Remove a directory from your PATH
zrc path remove ~/.local/bin
```

## Notes
- `zrc` directly modifies your shell configuration file.
- It stores the name of your shell config file in `~/.zrcc`.
