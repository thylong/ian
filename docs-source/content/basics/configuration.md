---
date: "2017-03-21T16:36:26-04:00"
title: Configuration
prev: /basics/getting-started
weight: 5
---

Ian relies a lot on configuration files as they describe your environments.
It's possible to edit your entire configuration using ian commands (recommended) but if you're stuck
or if you simply prefer to edit manually your configuration, it's possible !

Ian configuration files are written with Yaml and are relatively simple to read.
If you don't have specific ian files, no worries some are going to be generated
for you when you execute it for the first time.


## Configuration through CLI

### Config

config.yml contains all the settings that are not specific to Ian's feature.
This file is generated at the first run of Ian and can be edited at any time using Ian setup.

It contains currently 2 variables like in this example:

- **dotfiles_repository**: *thylong/dotfiles*
- **repositories_path**: */Users/thylong/www/repositories*

**dotfiles_repository** is the path to your dotfiles repository.
This variable is used to manage your dotfiles configuration (learn more about it here: https://github.com/webpro/awesome-dotfiles).
if you don't have one, you will not be able to use `ian env save` feature.

{{% notice info %}}
For now, only Github is supported but stories to add support for Gitlab and Bitbucket are already in the pipe.
If you wish to support any new repositories solution, don't hesitate to open an issue and/or make a pull request.
{{% /notice %}}

**repositories_path** is the fullpath to the directory that contains all your repositories.
This variable is use by a lot of Ian's commands to interact with your repositories,
by env commands to display stats, by the setup, etc.

{{% notice info %}}
For now, we support only having a single repositories path but as many languages have their specificities, I'm thinking about an easy way to have a more granular configuration if needed.
{{% /notice %}}

### Env

env.yml contains all the packages to be installed when setting up Ian on a new device.
This file is your best garantee to don't loose your configuration during a migration and
to keep it consistent when working on several computers.

The file content looks like the following example:
```yaml
    brew:
        - httpie
        - mongodb

    cask:
        - atom
        - caffeine
        - dash
        - google-chrome
        - iterm2
        - libreoffice
```

{{% notice note %}}
This file can contains packages that are not compatible with the current OS
you're working on, during setup Ian will simply ignore them.
{{% /notice %}}

## Editing yaml files

Ian configuration files can be found in `$HOME/.config/ian`.
