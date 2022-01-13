---
layout: page
title: Local Overrides
nav_order: 7
---

# Local Overrides

A single configuration is designed to be shared between multiple machines. The configuration can be tailored to each machine by using *overrides*.

Using overrides involves setting module variables to custom values on a per-machine basis. This is done by adding these variables to a `local.yml` file in the peridot config directory. This file should be added to your peridot configurations `.gitignore` file if you commit your setup.

For a peridot config like the following which installs the `git` and `neovim` apt package:

```yaml
modules:
  - name: apt
    source: builtin:apt
    variables:
      packages:
        - git
        - neovim
```

A `local.yml` file could be used to prevent the installation of the `neovim` package on a certain machine:

```yaml
variables:
  apt:
    - git
```

The override variables are grouped by module name (`apt` in the example above).

If you're looking for a way to set overrides by OS/Linux distribution, it is recommended to instead use [filters](modules/filters).
