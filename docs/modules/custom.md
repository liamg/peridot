---
layout: page
title: Custom Modules
nav_order: 2
parent: Modules
---

# Custom Modules

Custom modules can manage config files, run scripts, and/or include other modules.

The simplest possible module consists of a directory containing a `config.yml` file.

A project tree would look like the following when defining and using custom modules named `a` and `b`:

```
.
├── a
│   └── config.yml
├── b
│   └── config.yml
└── config.yml
```

The root config file would include the two custom modules using the following:

```yml
modules:
  - name: module A
    source: ./a
  - name: module B
    source: ./b
```

Check the [module config schema]({{ 'schema' | relative_url }}) for more information about building modules.

