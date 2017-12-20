#! /bin/bash

set -eux -o pipefail

# Install utilities

# Throwing away stdout logs because they were too plentiful and
# Circle couldn't display them in their web interface.
# Errors should still post to the console.
sudo apt-get update

# Install Node and Yarn
# Throwing away stdout logs because they were too plentiful and
# Circle couldn't display them in their web interface.
# Errors should still post to the console.
sudo curl -sL https://deb.nodesource.com/setup_6.x -o nodesource_setup.sh
sudo bash nodesource_setup.sh
sudo apt-get -qq install -y nodejs > /dev/null
npm install -g yarn
