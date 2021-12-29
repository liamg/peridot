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

## Filters

There will often be modules which should only be applied in a certain scenario. For example, a module that installs an `apt` package should not run on OSX, as Aptitude will not be available. To solve this problem, Peridot uses the concept of *filters* when a module is listed in a config file.

A filter is a list of values for which a match indicates the module should be applied. If there is no match, the module is simply ignored.

For example, the following module will only be applied on Linux, when the architecture is either `amd64` or `arm64`.

```yaml
modules:
    - name: filter example
      source: ./my-module
      variables:
        message: this module only runs on linux
      filters:
        arch:
          - amd64
          - arm64
        os: 
          - linux
```

The `filters` property allows the following:

| Filter | Description | Possible Values |
|--------|-------------|-----------------|
| os     | Operating System | `aix, android, darwin, dragonfly, freebsd, hurd, illumos, ios, js, linux, nacl, netbsd, openbsd, plan9, solaris, windows, zos` |
| distro | Linux Distribution (empty if not Linux) | `arch, centos, debian, fedora, ubuntu, rhel` (and others, dependant on content of `ID` in `/etc/os-release`) |
| arch   | System Architecture | `386, amd64, amd64p32, arm, arm64, arm64be, armbe, loong64, mips, mips64, mips64le, mips64p32, mips64p32le, mipsle, ppc, ppc64, ppc64le, riscv, riscv64, s390, s390x, sparc, sparc64, wasm`