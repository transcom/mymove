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
      name = "nodejs-14.18.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/staging";
      rev = "574ab51aa96701a1a6cec7435c6945a181fe2d63";
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
      name = "go-1.17.2";
      url = "https://github.com/NixOS/nixpkgs/";
      # Using master branch since 1.17.2 hasn't made it to nixpkgs-unstable yet
      ref = "refs/heads/master";
      rev = "db3aa421df73f43c03ad266619e22ce7c5354d92";
    }) {}).go_1_17

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "postgresql-12.7";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "860b56be91fb874d48e23a950815969a7b832fbc";
    }) {}).postgresql_12

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "pre-commit-2.14.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "229ff549e611d33d6b8ef2c6a2ce8879e9e75545";
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
      name = "aws-vault-6.3.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "54c1e44240d8a527a8f4892608c4bce5440c3ecb";
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

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "golangci-lint-1.42.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "1c804c475a6267dbe4c3f38a8f8d98067802daa0";
    }) {}).golangci-lint
];

  # the pre-commit hooks expects the binary to be `circleci`
  postBuild = ''
  ln -s $out/bin/circleci-cli $out/bin/circleci
  '';
}
