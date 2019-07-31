# How to Run Acceptance Tests

Before accepting work to merge into master different feature branches need to go through manual acceptance testing.
This document will get you set up for that testing and also show you how to run tests against feature branches.

## Setup

This setup doc assumes you're new to development environments and that you also haven't set anything up on your
computer to use the terminal. You'll have to do a modified version of what developers typically do to set up their
machines. However, if you are a developer you likely have all this set up already.

* Work with the Infrastructure team to ensure you have Github access and an AWS user provisioned
* Install [Homebrew](https://brew.sh)
  * Use the following command `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
* Install `git` with `brew install git`
* Clone the `transcom/mymove` repository to your laptop with `git clone git@github.com:transcom/mymove.git`
* Clone the `transcom/ppp-infra` repository to your laptop with `git clone git@github.com:transcom/ppp-infra.git`
* Change directories into the `transcom/mymove` directory and install other dependencies with `cd mymove && make docker_compose_setup`
* Work with someone to edit your `.envrc.local` file and set the `PPP_INFRA_PATH` to point at the directory for `transcom/ppp-infra`
* Work with someone to edit your `/etc/hosts` file to include the hosts required for this project.
* Confirm with the Infrastructure team that you can use the `aws` command. Try `aws s3 ls`.

At this point you should be ready to start running Acceptance Tests.

## Running Acceptance Tests

The first step in running acceptance tests is getting the branch name for the feature you wish to test.
A developer should be able to link you to the branch name in Pivotal or you can ask them directly. For this
set of instructions we'll use `branch_name`. From the terminal run:

```sh
git checkout branch_name
make docker_compose_up
```

At this point the server should be running and a lot of text will be scrolling by the screen. This is expected.

To log into the server you must browse to one of these websites:

* [Service Member login](http://milmovelocal:5000/devlocal-auth/login)
* [Office login](http://officelocal:5000/devlocal-auth/login)
* [TSP login](http://tsplocal:5000/devlocal-auth/login)
* [Admin login](http://adminlocal:5000/devlocal-auth/login)

**NOTE:** Unlike in development there is no `Local Sign In` button. This is because the production builds are
specifically disallowed from compiling that button into what we would deliver to production as a safety measure.

At this point you can run through any user flows that allow you test the feature and accept that it meets the
requirements needed for acceptance.

Finally, remember to shut down the working server:

```sh
make docker_compose_down
```

This should clean up the docker images downloaded to your computer and stop any running processes.
