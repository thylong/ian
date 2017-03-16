Ian
===

Ian is a simple tool to manage your development environment, repositories, and projects.


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
    Ian is a simple tool to manage your development environment, repositories, and projects.

    Learn more about Ian at http://goian.io

    Usage:
      ian [command]

    Default Commands:
      env         Manage development environment
      project     Interact with local projects
      repo        Manage local repositories
      setup       Set up ian configuration

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

Special thanks
==============
- `Devin Wilmot`_
- `Trisha Batchoo`_

.. _`template`: https://github.com/thylong/ian/blob/master/config/config_example.yml
.. _Homebrew: http://brew.sh
.. _Cask: https://caskroom.github.io
.. _RubyGems: https://rubygems.org/
.. _pip: https://packaging.python.org/
.. _npm: https://www.npmjs.com/
.. _apt: https://wiki.debian.org/Apt
.. _yum: https://fedoraproject.org/wiki/Yum
.. _`Devin Wilmot`: mailto:devwilmot@gmail.com
.. _`Trisha Batchoo`: https://github.com/tbat
