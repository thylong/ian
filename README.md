Ian
===

[![.github/workflows/github-actions-main.yml](https://github.com/thylong/ian/actions/workflows/github-actions-main.yml/badge.svg?branch=master)](https://github.com/thylong/ian/actions/workflows/github-actions-main.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/thylong/ian)](https://goreportcard.com/report/github.com/thylong/ian) [![Go Reference](https://pkg.go.dev/badge/github.com/thylong/ian.svg)](https://pkg.go.dev/github.com/thylong/ian)

Ian is a simple CLI tool to make your development environment portable.


Installing
==========

Ian comes as a binary. Once in your $PATH, you're good to go !

Linux
--------

```bash
go get github.com/thylong/ian
```

Mac OS X
--------

Ian requires Homebrew_.

```bash
go get -u github.com/thylong/ian
```

Usage
=====

```bash
ian

Ian is a simple tool to manage your development environment and repositories.

Learn more about Ian at http://goian.io

Usage:
    ian [command]

Default Commands:
    env         Manage development environment
    repo        Manage local repositories
    setup       Set up ian configuration

Use "ian [command] --help" for more information about a command.
```

Features
========

- Manage development environment (update OS package manager packages, etc).
- Self-Setup (import dotfiles and install related packages)
- Support [pip][pip], npm_, Homebrew_, RubyGems_, Cask_, apt_, yum_

Documentation
=============

Documentation can be seen here_. It was built thanks to the awesome Hugo project.
If you want to check the docs locally or to contribute to it, you can install hugo
and serve the static website locally using these commands:

```bash
brew update && brew install hugo
cd docs/
hugo server -t hugo-theme-learn --buildDrafts
```

Contributing
============

- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

Special thanks
==============
- [Devin Wilmot][devin]
- [Trisha Batchoo][trisha]
- [Carl Chinatomby][carl]

[template]: https://github.com/thylong/ian/blob/master/config/config_example.yml
[pip]: https://packaging.python.org/
[homebrew]:  http://brew.sh
[here]: https://goian.io
[cask]: https://caskroom.github.io
[rubygems]: https://rubygems.org/
[npm]: https://www.npmjs.com/
[apt]: https://wiki.debian.org/Apt
[yum]: https://fedoraproject.org/wiki/Yum
[devin]: mailto:devwilmot@gmail.com
[trisha]: https://github.com/tbat
[carl]: https://github.com/Carl-Chinatomby