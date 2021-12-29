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

Variables for a module come from the module defaults, the parent modules variables for that module, or from the global overrides.

Global overrides are looked at first, then configured values, then defaults. If no value is found, a template error will occur and execution will fail.
