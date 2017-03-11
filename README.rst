Ian
===

Ian is a CLI interface to setup, manage and save your dev environment.


Installing
==========

Ian comes as a binary. Once in your $GOPATH, you're good to go !

Linux
--------

.. code-block:: console

    $ go get github.com/thylong/ian


Mac OS X
--------

Ian requires Homebrew_.

.. code-block:: console

    $ go get github.com/thylong/ian


Usage
=====

.. code-block:: console

    $ ian
    Ian is a very simple automation tool for developer with dev environment.

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

- Manage dev environment (update OS package manager packages, Fetch diff on repos, etc).
- Interact with git repositories.
- Manage projects (deploy, rollback, healthcheck, etc).
- Self-Setup (import dotfiles and install related packages)
- Support pip_, npm_, Homebrew_, RubyGems_, Cask_, apt_, yum_


Contributing
============

- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request

.. _`template`: https://github.com/thylong/ian/blob/master/config/config_example.yml
.. _Homebrew: http://brew.sh
.. _Cask: https://caskroom.github.io
.. _RubyGems: https://rubygems.org/
.. _pip: https://packaging.python.org/
.. _npm: https://www.npmjs.com/
.. _apt: https://wiki.debian.org/Apt
.. _yum: https://fedoraproject.org/wiki/Yum

Special thanks
==============
- Devin Willmot
- Trisha Batchoo
