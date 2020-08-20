# Use ASDF To Manage Node, Yarn, and Golang Versions In Development

**This supersedes [ADR 0046 Use nodenv](0046-use-nodenv.md)**

There are many tools for managing versions of developer tools on developer machines. [brew](https://brew.sh/), [nodenv](https://github.com/nodenv/nodenv), [g](https://github.com/stefanmaric/g), etc. Historically MilMove has used brew for many things, but for node and golang this has lead to problems. Because of the issues around requiring specific versions of node and golang brew has caused more headaches than it is worth dealing with. This lead to [ADR 0046 Use nodenv](0046-use-nodenv.md), which solved the problem for node. However we don't have one for golang. So this ADR aims to provide a recommendation towards managing golang release versions in development.

## Considered Alternatives

* Do nothing, keep using brew
* Use a golang specific tool
* Use [asdf](https://asdf-vm.com/) to manage golang, node, and yarn

## Decision Outcome

* Chosen Alternative: Use [asdf](https://asdf-vm.com/) to manage golang, node, and yarn
  * `+` asdf supports multiple language version management, thus simplifying the amount of tools we need to install
  * `+` asdf supports a `.tool-version` config file within the project to define what versions are required and can be checked into the repo
  * `+` will allow us to define what version of node, yarn, and golang all developers are using so it is consistent
  * `+` removes dependence on brew which regularly only has the latest and greatest version of these tools.
  * `-` asdf is yet another tool to be familiar with

## Pros and Cons of the Alternatives

### Do nothing, keep using brew

* `+` Easy, nothing to do here
* `+` It got us here
* `-` It's often confusing to fix if a newer version of node, golang, or yarn is installed
* `-` By default brew doesn't save old versions, so if that version is gone from the brew repo there is no rolling back
* `-` Pinning brew packages leads to issues installing other brew packages that depend on them, which may be unrelated to the project

### Use a golang specific tool

* `+` Solves the version control issues introduced by brew
* `+` Similar to how we handled node
* `-` Yet another tool
* `-` Yet another configuration

### Use asdf to manage golang, node, and yarn

* `+` Solves the version control issues introduced by brew
* `+` Similar to how we handled node
* `+` Allows us to use one tool to manage node, yarn, and golang so we can have matching versions in development and deployed environments
* `+` One configuration file for all these dependencies
* `-` Pain of switching tools around for all engineers
