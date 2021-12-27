# peridot

User environment management across multiple machines, operating systems, users etc.

Manage your dotfiles and more!

## Getting Started

*Guide coming soon*

## Commands

Below are some of the most immediately useful peridot commands.

Use `peridot --help` for more information on all available commands and options.

### `init`

Initialises a new peridot config for the local user environment.

### `use`

Uses a remote config as a basis for managing the local user environment.

This takes a GitHub `owner/repo` repository alias of a repository which you have read/write access for. It's recommended to create a dedicated GitHub repository for your peridot config files.

### `diff`

Compares the desired state as dictated by your peridot templates and config files with the actual local state.

### `apply`

Applies changes to the local state to conform to your peridot templates and config files. You can preview these changes before making them using the `diff` command.

## Usage

## Variables

| Name          | Description |
|---------------|-------------|
| config_dir    | `$XDG_CONFIG_HOME` if it exists, otherwise `~/.config`
| 