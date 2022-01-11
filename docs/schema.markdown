---
layout: page
title: Schema
nav_order: 4
---

# Config Schema

Modules can be nested inside one another indefintely, using a child directory per module. The below is a commented example to demonstrate the schema for the root module config file.

```yaml
{% raw %}
files: # a list of files which should be managed
  - target: "{{ .user_home_dir }}/example.conf" # the path to the config file to manage
    source: "./example.conf.tmpl" # the path to the template for the targeted template file 
    disable_templating: false # enable this to include "{{" characters in your config files without having to escape them.
modules: # a list of child modules which should be managed
  - name: "my module" # a name to uniquely identify the module within the config file
    source: "./submodule" # a relative path to a submodule directory (or a url, or a builtin module identifier)
    depends_on: other_module_name # another module name within the config file that should be applied before this one
    filters: # a set of filters which must be satisfied to apply this module
      os: linux # operating system
      distro: ubuntu # linux distribution
      arch: amd64 # architecture
    variables: # a map of variable to pass into the module (see individual module docs)
      x: 1
      y: 2
scripts: # a set of scripts to manage the lifecycle of a module
  should_install: # a script that dictates whether a module should be installed (i.e. whether the install script should be run)
    command: ./should_install.sh # a path to the script to run. if the script exits with status '0', the module should be installed
    sudo: false # whether the script should be launched with sudo
  install:
    command: ./install.sh
    sudo: true
  should_update:
    command: ./should_update.sh # a path to the script to run. if the script exits with status '0', the module should be updated
    sudo: false
  update:
    command: ./update.sh
    sudo: true
  after_file_change: # a script that should run after a file within the "files" list is written
    command: ./after.sh
    sudo: false
{% endraw %}
```

## Module Schema

Non-root modules/submodules use all of the above schema, but can also take input variables via the `variables` section, as demonstrated below:

```yaml
variables:
  - name: x
    default: 100
    required: false
  - name: y
    default: 200
    required: false
```

This allows a module to be customised using these inputs. For example, a file managed by a module with the above variables could use `{% raw %}{{ .x }}{% endraw %}` and `{% raw %}{{ .y }}{% endraw %}` to access these variables within a template.

