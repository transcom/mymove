# Use ASDF To Manage Golang Versions In Development

There are many tools for managing versions of developer tools on developer machines. [brew](https://brew.sh/), [nodenv](https://github.com/nodenv/nodenv), [g](https://github.com/stefanmaric/g), etc. Historically MilMove has used brew for many things, but for node and golang this has lead to problems. Because of the issues around requiring specific versions of node and golang brew has caused more headaches than it is worth dealing with. This lead to [ADR 0046 Use nodenv](0046-use-nodenv.md), which solved the problem for node. However we don't have one for golang. So this ADR aims to provide a recommendation towards managing golang release versions in development.

## Considered Alternatives

* Do nothing, keep using brew
* Use a golang specific tool (g, goenv, gvm, etc.)
* Use [asdf](https://asdf-vm.com/) to manage golang, node, and yarn
* Use [asdf](https://asdf-vm.com/) to manage only golang for now

## Decision Outcome

* Chosen Alternative: Use [asdf](https://asdf-vm.com/) to manage only golang for now
  * `+` asdf supports a `.tool-version` config file within the project to define what version of golang are required and can be checked into the repo
  * `+` will allow us to define what version of golang all developers are using so it is consistent
  * `+` removes dependence on brew which regularly only has the latest and greatest version of these tools.
  * `+` allows us to switch node and yarn to be managed by asdf at a later time
  * `-` asdf is yet another tool to be familiar with
  * `=` Our use and dependence on brew will not be removed, just it's use to install yarn and golang

To elaborate some more on why ASDF verses a golang specific option. In googling there were not many discussions around version managers. Polling among Truss revealed use of [g](https://github.com/stefanmaric/g) and `asdf`. I did like the tool `g` and the maintainer's [reasons for preferring it](https://github.com/stefanmaric/g#the-alternatives-and-why-i-prefer-g) to other such tools. However, the only negative listed to `asdf-golang` was it's dependence on `asdf`, and didn't account for the fact that having one tool version manager simplifies one set of tools required for keeping development flowing smoothly. Admittedly we are not yet taking advantage of `asdf` for other tools, but I would like to see that happen.

Another advantage of these tools is they also runs on Linux so we can utilize it within docker image creation to manage versions of our tooling that is installed and keep it in sync with what developers are using locally from a single tool configuration file that is checked into the repo. This ADR advocates changing our tooling to rely on `asdf` locally to manage golang and where it makes sense in our [circleci-docker](https://github.com/transcom/circleci-docker) images.

## Pros and Cons of the Alternatives

### Do nothing, keep using brew

* `+` Easy, nothing to do here
* `+` It got us here
* `-` It's often confusing to fix if a newer version of node, golang, or yarn is installed
* `-` By default brew doesn't save old versions, so if that version is gone from the brew repo there is no rolling back
* `-` Pinning brew packages leads to issues installing other brew packages that depend on them, which may be unrelated to the project
* `=` Our use and dependence on brew will not be removed, just it's use to install yarn and golang

### Use a golang specific tool (g, goenv, gvm, etc.)

* `+` Solves the version control issues introduced by brew
* `+` Similar to how we handled node
* `-` Yet another tool
* `-` Yet another configuration
* `=` Our use and dependence on brew will not be removed, just it's use to install yarn and golang

### Use asdf to manage golang, node, and yarn

* `+` Solves the version control issues introduced by brew
* `+` Similar to how we handled node
* `+` Allows us to use one tool to manage node, yarn, and golang so we can have matching versions in development and deployed environments
* `+` One configuration file for all these dependencies
* `+` Clear documentation
* `-` Pain of switching tools around for all engineers
* `-` Possible conflicts if developers have multiple of these tools installed to manage any one of the tools mentioned
* `=` Our use and dependence on brew will not be removed, just it's use to install yarn and golang

A disadvantage of this approach, and part of why it was not selected is that many of these tool rely on your shell path to work. So installing node, yarn, or golang with asdf and another one of the similar tools can cause conflicts depending on which of the version manager tools is first in a particular path. This risk exists today, but could be exacerbated by choosing tooling that is less common. I think using `direnv` or other method to control the path can allow the tools to co-exist for those that need it though. Another way to mitigate this is turn on asdf's legacy configuration feature as it will pickup config files for other tools line nodenv that way. See [this documentation](https://asdf-vm.com/#/core-configuration?id=homeasdfrc)

### Use asdf to manage only golang for now

* `+` asdf supports a `.tool-version` config file within the project to define what version of golang are required and can be checked into the repo
* `+` will allow us to define what version of golang all developers are using so it is consistent
* `+` removes dependence on brew which regularly only has the latest and greatest version of these tools.
* `+` allows us to switch node and yarn to be managed by asdf at a later time
* `+` Supports configuration file for setting the golang version and checking in
* `+` Clear documentation
* `-` asdf is yet another tool to be familiar with
* `=` Our use and dependence on brew will not be removed, just it's use to install yarn and golang
