---
date: "2017-03-21T16:36:26-04:00"
title: Contributing
weight: 0
icon: "<b>2. </b>"
---

Any contribution to this project is more than welcomed ! As small or big as they are, contributions are what will make this project better. There is still a lot to do : reporting bugs, proposing fixes, new features, helping with the documentation...

If it’s your first contribution, no worries. As long as you follow the guidelines, everything is going to be fine :)

## Open an issue

If you spot a bug or want to propose an enhancement, you can open an issue on Github. Please make sure to avoid duplicates and explain clearly what is your need.

## Coding guidelines

This project wants to stay close to the standard and was inspired by Kubernetes CLI, Cobra and Requests code quality.

My advice would be to follow effective Go guide and to listen to golinter suggestions.

## Creating a pull request

- Fork the project
- Install dev dependencies and package
- Checkout a new branch based on dev
- Make your change and run the test suite
- Be sure to respect Go conventions
- Run the entire test suite before commiting
- If your commit fixes an open issue, reference it in the commit message
- Commit with a proper commit message
- Open a PR using this description structure
- Travis (CI) will run on your branch
- If the PR is accepted, it will be merged into the dev branch and then released.
