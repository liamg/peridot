---
layout: page
title: Variables
nav_order: 5
---

# Variables

| Name               | Description |
|--------------------|-------------|
| user_config_dir    | `$XDG_CONFIG_HOME` if it exists, otherwise `~/.config`.
| user_home_dir      | Current user's home directory.
| sys_os             | Operating system, e.g. `linux`, `macos`, `windows`
| sys_distro         | Linux distribution e.g. `ubuntu`, `arch` etc. Empty if non-Linux
| sys_arch           | System architecture, e.g. amd64

Variables for a module come from the module defaults, the parent modules variables for that module, or from the global overrides.

Global overrides are looked at first, then configured values, then defaults. If no value is found, a template error will occur and execution will fail.
