# How to Switch Local Node Version

Switching the local node version should be something that happens as project dependencies and libraries change. The MilMove
project uses [NVM](https://github.com/nvm-sh/nvm) for switching between different local versions of node.

## Installing NVM

1. Install `NVM` via `brew install nvm`.
1. Follow the instructions given by brew:
   You should create NVM's working directory if it doesn't exist:

  ```mkdir ~/.nvm```

   Add the following to `~/.bash_profile` or your desired shell configuration file:

    export NVM_DIR="$HOME/.nvm"
    [ -s "/usr/local/opt/nvm/nvm.sh" ] && . "/usr/local/opt/nvm/nvm.sh"  # This loads nvm
    [ -s "/usr/local/opt/nvm/etc/bash_completion" ] && . "/usr/local/opt/nvm/etc/bash_completion"  # This loads nvm bash_completion

1. `source ~/.bash_profile`
1. `nvm install <NODE_VERSION> && NVM use <NODE_VERSION>`

## Remove NVM

1. To remove, delete, or uninstall `nvm` - just remove the `$NVM_DIR` folder, `~/.nvm`.
1. Remove the `nvm` references from `~/.bash_profile`
