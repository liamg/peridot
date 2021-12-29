---
layout: page
title: Modules
nav_order: 4
has_children: true
---

# Modules

Peridot modules are a useful way of organising configuration into reusable, simple building blocks.

A module is typically a single directory containing a `config.yml`, and optionally some supporting files. Modules can be nested inside one another indefinitely, so you can break down a module into sensible submodules where required, and so on.

There are three types of module: [built-in](builtins), [custom](custom) and [community](community). Whilst you can achieve anything by writing your own custom modules, built-in and community modules often provide off-the-shelf setups that make progress much faster.

You can use a module by adding an item to the `modules` list in your main `config.yml`.

For example:

```yaml
modules:
    - name: git example
      source: builtin:git
      variables:
        username: Bobby Tables
        email: bobby@tabl.es
```

The important parameter here is `source`, as it tells peridot which module to load. Built-in modules have a `builtin:` prefix, locally defined custom modules have a `./` prefix, and community modules use URLs to a `.tar.gz` file containing module source code.

Modules can typically be configured using [variables](../variables) (using the `variables` parameter seen in the example above). To see which variables are available, you'll need to check out the documentation of the specific module you're using.
