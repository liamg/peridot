# peridot

User environment management for Linux/OSX. Think Terraform for your dotfiles and local setup that can be replicated across your machines, regardless of differences in distro, window manager etc.

Allows the templating of any configuration file to provide automatic tweaking of configuration on apply.

If you struggle with maintaining your dotfiles in one big git repository and/or sharing them across multiple machines, Peridot is for you.

## Documentation

Please see the [official documentation](https://www.liam-galvin.co.uk/peridot/guide).

## TODO

- add `depends_on` to modules list
- support sudo
- add more useful builtins
- update/publish mapped to git pull && apply/git commit && git push
- add a gif to the readme + docs
- finish documentation
- write contribution guide
- add `upgrade` command to update peridot binary
