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
      project     Interact with listed project
      repo        Manage stored repositories
      setup       Setup ian working environment
      version     Print the version information


    Use "ian [command] --help" for more information about a command.

Features
========

- Self-Setup (import dotfiles and install related packages)
- Interact with local git repositories.
- Get environment infos.
- Manage projects (deploy, rollback, healthcheck, etc).
- Manage dev environment (update OS package manager packages, Fetch diff on repos, etc).
- Support pip, Npm, Brew, RubyGems, Cask, Apt, Yum

# TODO v0.1
- TESTS (env, projects, setup)
- Remove logic from Cobra commands
- Fix english typos
- Ensure commands architecture
- Travis workflow
- tag v0.1

# TODO v0.2
- Get started experience (website + form on site to determine the profile)
- Presets by profiles (backend, frontend, fullstack)
- Customize presets easily http://getbootstrap.com/customize/ (UI local/remote? Hash copy/paste in CLI ?)

# TODO v0.3
- Add tests
- Add a project set <field_name> <value> to set the config from the terminal

# TODO v0.4
- Cyphering / Decyphering of the configuration (--encrypted option to env save)
- Safe export functionnality (export only the non sensible infos) ("--safe" option to env save)

# TODO v0.5
- Specific package management with Ian


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
