---
layout: page
title: Commands
nav_order: 5
---

# Commands

Below are some of the most immediately useful peridot commands.

Use `peridot --help` for more information on all available commands and options.

## `init`

Initialises a new peridot config for the local user environment.

## `validate`

Validates the peridot config files, printing any errors to the terminal and exiting with a non-zero status if any are found.

## `diff`

Compares the desired state as dictated by your peridot templates and config files against the actual local state.

## `apply`

Applies changes to the local state to conform to your peridot templates and config files. You can preview these changes before making them using the `diff` command.

## `system`

Print local system information to the terminal, to be used with the [filters](modules/filters) feature.

