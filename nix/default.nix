# to install
# nix-env -p /nix/var/nix/profiles/mymove -f nix -i
#
# use
#
# https://lazamar.co.uk/nix-versions/
# to find rev for specific package version

let
  pkgs = import <nixpkgs> {};
  inherit (pkgs) buildEnv;
in buildEnv {
  name = "mymove-packages";
  paths = [

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "bash-5.1-p4";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).bash_5

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "nodejs-14.17.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "75916fb375eb571dceebef84263a6cb942372769";
    }) {}).nodejs-14_x

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "yarn-1.22.10";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).yarn

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "go-1.16.7";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "08c66abd89d9be404c4decd360746095a4bad0a5";
    }) {}).go

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "postgresql-12.7";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "860b56be91fb874d48e23a950815969a7b832fbc";
    }) {}).postgresql_12

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "pre-commit-2.13.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "f62847d5655a6954fc946e5d9c81424171cc281e";
    }) {}).pre-commit

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "circleci-cli-0.1.11540";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).circleci-cli

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "jq-1.6";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).jq

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "shellcheck-0.7.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).shellcheck

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "opensc-0.21.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).opensc

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "entr-4.6";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).entr

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "aws-vault-6.2.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).aws-vault

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "watchman-4.9.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).watchman

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "awscli2-2.1.7";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).awscli2

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "chamber-2.9.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).chamber

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "dbmate-1.11.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).dbmate

    (import (builtins.fetchGit {
         # Descriptive name to make the store path easier to identify
         name = "nss-cacert-3.60";
         url = "https://github.com/NixOS/nixpkgs/";
         ref = "refs/heads/nixpkgs-unstable";
         rev = "559cf76fa3642106d9f23c9e845baf4d354be682";
    }) {}).cacert

    (import (builtins.fetchGit {
         # Descriptive name to make the store path easier to identify
         name = "hadolint-2.1.0";
         url = "https://github.com/NixOS/nixpkgs/";
         ref = "refs/heads/nixpkgs-unstable";
         rev = "68e901c68446539165b19b227f51f7fc2b3655f9";
    }) {}).hadolint
];

  # the pre-commit hooks expects the binary to be `circleci`
  postBuild = ''
  ln -s $out/bin/circleci-cli $out/bin/circleci
  '';
}
