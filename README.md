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

| Name               | Description |
|--------------------|-------------|
| user_config_dir    | `$XDG_CONFIG_HOME` if it exists, otherwise `~/.config`.
| user_home_dir      | Current user's home directory.

Variables for a module come from the module defaults, the parent modules variables for that module, or from the global overrides.

Global overrides are looked at first, then configured values, then defaults. If no value is found, a template error will occur and execution will fail.

Load config -> read innermodules and vars. 
For each inner module, parse the input variables from the actual module, merge these defaults with the inner module values, and finally merge with the global overrides (these are namespaced by module.)
