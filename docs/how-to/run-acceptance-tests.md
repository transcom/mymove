# How to Run Acceptance Tests

Before accepting work to merge into master different feature branches need to go through manual acceptance testing.
This document will get you set up for that testing and also show you how to run tests against feature branches.

## Setup

[Video Walk Through](https://drive.google.com/drive/folders/1VzFlHuJKnQ4V1TWL5taRi0QkjP6yVzNT)

This setup doc assumes you're new to development environments and that you also haven't set anything up on your
computer to use the terminal. You'll have to do a modified version of what developers typically do to set up their
machines. However, if you are a developer you likely have all this set up already.

**NOTE:** Many of these instructions should only be run once. If you run it and have an error then a few of these
commands will not work as intended the second time. Please reach out to engineering or infra for help.

* Work with the Infrastructure team to ensure you have Github access and an AWS user provisioned
* Install [Homebrew](https://brew.sh)
  * Use the following command `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
  * Update your `$PATH` in your `~/.bash_profile` with `echo "export PATH=$HOME/bin:/usr/local/bin:/usr/local/sbin:/usr/local/opt/openssl/bin:/sbin:$PATH" >> ~/.bash_profile`
  * Then update our terminal with the new changes `source ~/.bash_profile`
* Install `git` with `brew install git`
* Get the Project Code
  * If you do not have a directory for code then make one and move into it with `mkdir -p ~/Projects && cd ~/Projects`
  * Visit the [Github Personal Access Tokens](https://github.com/settings/tokens) page to generate a new token. Name the token "MY_NAME Truss Repo Token" and select the `repo` permissions and generate. Then copy this token and put it into your 1Password (You will not be able to see it again). You will use this token instead of a password when cloning the repositories. For more on 2FA and Github you can visit [Two-factor Authentication with Github](https://github.blog/2013-09-03-two-factor-authentication/#how-does-it-work-for-command-line-git).
  * Clone the `transcom/mymove` repository to your laptop with `git clone https://github.com/transcom/mymove.git`. No username/password should be required.
  * Clone the `transcom/ppp-infra` repository to your laptop with `git clone https://github.com/transcom/ppp-infra.git`. Your username and password will be needed. The password is the generated personal access token you made previously.
* You will need to modify your `/etc/hosts` file to include the hosts required for this project.
  * Run `make check_hosts` and follow any of the instructions that it presents. Those instructions will likely look like:

  ```sh
  echo "127.0.0.1 milmovelocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 officelocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 tsplocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 orderslocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 adminlocal" | sudo tee -a /etc/hosts
  ```

* Change directories into the `transcom/mymove` directory and install other dependencies with `cd mymove && make docker_compose_setup`
* Update your `~/.bash_profile` to install `direnv` correctly.
  * With this command: `echo "if command -v direnv >/dev/null; then eval \"\$(direnv hook bash)\"; fi" >> ~/.bash_profile`
  * Then update our terminal with the new changes `source ~/.bash_profile`
* Edit your `.envrc.local` file and set the `PPP_INFRA_PATH` to point at the directory for `transcom/ppp-infra`
  * With this command: `echo "export PPP_INFRA_PATH=$HOME/Projects/transcom/ppp-infra" >> ~/Projects/mymove/.envrc.local`
  * Now set up everything by running `direnv allow`
* Confirm with the Infrastructure team that you can use the `aws` command. Try `aws s3 ls`.

At this point you should be ready to start running Acceptance Tests.

## Running Acceptance Tests

The first step in running acceptance tests is getting the branch name for the feature you wish to test.
A developer should be able to link you to the branch name in Pivotal or you can ask them directly. For this
set of instructions we'll use `branch_name`. From the terminal run:

```sh
cd ~/Projects/mymove
git pull
git checkout branch_name
direnv allow
make docker_compose_up
```

At this point the server should be running and a lot of text will be scrolling by the screen. This is expected.

To log into the server you must browse to one of these websites:

* [Service Member login](http://milmovelocal:4000/devlocal-auth/login)
* [Office login](http://officelocal:4000/devlocal-auth/login)
* [TSP login](http://tsplocal:4000/devlocal-auth/login)
* [Admin login](http://adminlocal:4000/devlocal-auth/login)

**NOTE:** Unlike in development there is no `Local Sign In` button. This is because the production builds are
specifically disallowed from compiling that button into what we would deliver to production as a safety measure.

At this point you can run through any user flows that allow you test the feature and accept that it meets the
requirements needed for acceptance.

Finally, remember to shut down the working server:

```sh
make docker_compose_down
```

This should clean up the docker images downloaded to your computer and stop any running processes.
