---
layout: page
title: Ordering
nav_order: 5
parent: Modules
---

# Ordering Modules

Modules often depend on each other, and therefore it is required to dictate the order in which they are applied. This can be achieved using the `depends_on` attribute in your module usage.

```yaml
modules:
  - name: install
    source: builtin:apt
    variables:
      packages:
        - git
  - name: clone
    source: ./cloner
    depends_on:
      - install
```

The `depends_on` attribute is a list of all modules (by `name`) that the module depends on being applied before it can be applied itself.

