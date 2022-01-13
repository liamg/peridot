---
layout: page
title: Getting Started
nav_order: 3
---

# Getting Started

## Initialisation

After [installing]({{ 'install' | relative_url }}), the first action is usually to initialise peridot. If you already have a peridot configuration stored in a git repository, you can clone it to `~/.config/peridot`, otherwise, initialise a fresh configuration with:

```bash
peridot init
```

This will create a config file in `~/.config/peridot/config.yml` (or `$XDG_CONFIG_HOME/peridot/config.yml` if the environment variable is set.)

By itself, the config file will not do anything, it is empty by default.

You can see a [full breakdown of the config file format]({{ 'schema' | relative_url }}), but for now we'll delve into some conrete examples...

## Managing your first dotfile

In this example we'll configure a fictional app `my-app`, which we'll imagine uses a config file found in `~/.my-app.conf`.

First we'll modify our `config.yml` (see above) to target the `my-app` config file:

```yml
files:
    - target: "{% raw %}{{ .user_home_dir }}{% endraw %}/my-app.conf"
      source: ./my-app.conf.tmpl
```

You'll likely notice a few interesting things about the above example. Let's go through it...

Each item in the `files` list targets exactly one file. The `target` property refers to the location of the config file that should be written on the fileystem. In our case, we want to target `~/my-app.conf`, but you'll notice we're using `{% raw %}{{ .user_home_dir }}{% endraw %}/my-app.conf` instead - this is because we're using the `user_home_dir` variable to insert the full path to the current user's home directory. There are other [variables]({{ 'variables' | relative_url }}) available, and it's possible to define your own too.

The `source` property refers to the file to use as a template for the target file, relative to the directory of the `config.yml` it is defined in. It must always start with `./`. This means that in our example, if our peridot config is in `~/.config/peridot/config.yml`, the template file should be defined at `~/.config/peridot/my-app.conf.tmpl`. You don't have to use the `.tmpl` extension, but it's often useful to make it clear that the file in question is a template.

Let's say we want our config to switch on a special feature for all of our Linux machines - the template could look something like this:

```toml
{% raw %}
[config]
{{ if eq .sys_os "linux" }}enable_special_feature = true
{{ else }}enable_special_feature = false
{{ end }}
{% endraw %}
```

You can explore the functionality available in templates by reading up on the [Go text/template system](https://pkg.go.dev/text/template).

At this point, our configuration is complete. 

Peridot can tell us what changes it will make before it changes anything, to ensure our config is correct. We can do this with the `diff` command, like so:

```bash
peridot diff
```

You can run this command from anywhere on your system - Peridot knows where to find the config files.

We'll see some output that describes the change that our config files describe:

![diff example]({{ 'diff.png' | relative_url }})

If we're happy with the diff, we can go ahead and apply the change with the `apply` command.

```bash
peridot apply
```

If everything works, we'll see:

```
[Module root] Changes applied.

1 modules applied successfully.
```

And our `my-app` config file at `~/my-app.conf` will now contain (on Linux):

```
❯ cat ~/my-app.conf
[config]
enable_special_feature = true
```

The `--debug` flag is available for debugging this process. It provides stdout/stderr for all scripts run by peridot, and an internal log of peridot's actions too.

## Utilising your first built-in module

Several modules are included as part of the peridot binary, meaning you can get started using them without having to write any templates.

One of such modules is the [apt built-in]({{ 'modules/builtins/#builtinpacman' | relative_url }}). It enables you to specify one or more pacman packages for installation.

For example:

```yaml
modules:
  - builtin:pacman
    variables:
      packages:
        - neovim
        - firefox
```

When this configuration is applied, the `neovim` and `firefox` packages will be installed using pacman. 

Of course, `pacman` won;t be available on many systems, so we can add a [filter]({{ 'modules/filters' | relative_url }}) to ensure this module will only be applied on Arch Linux systems:

```yaml
modules:
  - builtin:pacman
    filters:
      os: linux
      distro: arch
    variables:
      packages:
        - neovim
        - firefox
```



