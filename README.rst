Ian
===

Ian is CLI interface to interact with you Mac environment.


Installing
==========

(Required) To work properly, Ian requires Homebrew and Golang to be installed
with a valid $GOPATH.

.. code-block:: console

    $ go get github.com/thylong/ian
    $ ian setup


Usage
=====

.. code-block:: console

    $ ian
    Ian is a very simple automation tool for developer with Mac environment.

    Usage:
      ian [command]

    Default Commands:
      env         Get infos about the local environment
      news        Retrieve last news from Google News
      packages    packages allows you to manage ian extensions
      project     Interact with listed project
      repo        Manage stored repositories
      setup       Setup ian working environment
      version     Print the version information

    Package Commands:
      baily-cli baily-cli type:npm


    Use "ian [command] --help" for more information about a command.

Features
========

- Self-Setup (import dotfiles, install brew, cask and related packages)
- Interact with local git repositories.
- Get environment infos.
- Manage projects (deploy, rollback, healthcheck, etc).
- Manage dev environment (update OS package manager packages, Fetch diff on repos, etc).
- Possibility to extend with other packages
- Support pip, Npm, Brew\

# TODO v0.1
- ian env save
- Cyphering / Decyphering of the configuration (--encrypted option to env save)
- Safe export functionnality (export only the non sensible infos) ("--safe" option to env save)
- Improve CLI architecture (KISS)

# TODO v0.2
- Specific package management with Ian
- Add a project set <field_name> <value> to set the config from the terminal
- Add tests

Contributing
============

- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

.. _`template`: https://github.com/thylong/ian/blob/master/config/config_example.yml
.. _Brew: http://brew.sh
.. _Cask: https://caskroom.github.io
