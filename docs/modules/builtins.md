---
layout: page
title: Built-in Modules
nav_order: 1
parent: Modules

---

# Built-in Modules

Built-in modules are available by default with a Peridot install. They get their name because they are compiled into the Peridot binary (you can see the source code for them [on GitHub](https://github.com/liamg/peridot/tree/main/internal/pkg/builtins)).

You can use a built-in module by adding a `modules` item to your main config file with the `source` parameter set to `builtin:` followed by the module name, for example:

```yaml
modules:
    - name: example
      source: builtin:git
```

The following set of built-in modules in currently available.

## `builtin:git`

Manages your Git configuration (`~/.gitconfig`) and your global git ignore patterns (`~/.gitignore`).

| Name        | Required | Default    | Description |
|-------------|----------|------------|-------------|
| email       | yes      |            | Git email. Used to set `user.email` in `.gitconfig`.
| username    | yes      |            | Git user name. Used to set `user.name` in `.gitconfig`.
| editor      | no       | vim        | Editor to use for commit messages etc. Used to set `core.editor` in `.gitconfig`.
| signingkey  | no       |            | GPG signing key used to sign commits
| aliases     | no       |            | List of git aliases
| ignores     |          |            | List of ignore patterns to apply globally
| extra       |          |            | Any extra configuration as a single multiline string

## `builtin:fonts`

Install one or more fonts.

| Name        | Required | Default    | Description |
|-------------|----------|------------|-------------|
| files       | yes      |            | A list of font files to install. Can be local file paths relative to your config (must start with `./`) or URLs to font files.
| dir         | no       | `~/.local/share/fonts` | 

## `builtin:apt`

Install one or more `apt` packages.

| Name        | Required | Default    | Description |
| ----------- | -------- | ---------- | ----------- |
| packages    | yes      |            | A list of packages to install.

## `builtin:pacman`

Install one or more `pacman` packages.

| Name        | Required | Default    | Description |
| ----------- | -------- | ---------- | ----------- |
| packages    | yes      |            | A list of packages to install.

## `builtin:yay`

Install one or more `aur` packages to install with `yay`.

| Name        | Required | Default    | Description |
| ----------- | -------- | ---------- | ----------- |
| packages    | yes      |            | A list of packages to install.

*Can't find what you're looking for? PRs are very welcome!*

