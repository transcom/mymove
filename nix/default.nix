# use ./nix/update.sh to install
#
# use <https://ahobson.github.io/nix-package-search> to find a package version

let
  pkgs = import <nixpkgs> { };
  inherit (pkgs) buildEnv;
in
buildEnv {
  name = "mymove-packages";
  paths = [

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "bash-5.2-p15";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).bash_5

    (import
      (builtins.fetchGit {
        # Descriptive name to make the store path easier to identify
        name = "nodejs-16.15.0";
        url = "https://github.com/NixOS/nixpkgs/";
        ref = "refs/heads/nixpkgs-unstable";
        rev = "0b45cae8a35412e461c13c5037dcdc99c06b7451";
      })
      { }).nodejs-16_x

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "yarn-1.22.19";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).yarn

    (import
      (builtins.fetchGit {
        # Descriptive name to make the store path easier to identify
        name = "go-1.19.3";
        url = "https://github.com/NixOS/nixpkgs/";
        ref = "refs/heads/nixpkgs-unstable";
        rev = "c4ba130a43d716a2e042222231471e2d60790aa6";
      })
      { }).go_1_19

    (import
      (builtins.fetchGit {
        # Descriptive name to make the store path easier to identify
        name = "postgresql-12.7";
        url = "https://github.com/NixOS/nixpkgs/";
        ref = "refs/heads/nixpkgs-unstable";
        rev = "860b56be91fb874d48e23a950815969a7b832fbc";
      })
      { }).postgresql_12

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "python3.10-pre-commit-2.20.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).pre-commit

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "circleci-cli-0.1.22924";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).circleci-cli

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "jq-1.6";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).jq

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "shellcheck-0.8.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "293a28df6d7ff3dec1e61e37cc4ee6e6c0fb0847";
    }) {}).shellcheck

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "opensc-0.23.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).opensc

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "entr-5.2";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).entr

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "aws-vault-6.6.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).aws-vault

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "watchman-4.9.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).watchman

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "awscli2-2.9.13";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).awscli2

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "chamber-2.11.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).chamber

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "nss-cacert-3.86";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).cacert

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "hadolint-2.8.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "316e4a4ca184512a1620baf014452c97bc11b025";
    }) {}).hadolint

    (import
      (builtins.fetchGit {
        # Descriptive name to make the store path easier to identify
        name = "schemaspy-6.1.0";
        url = "https://github.com/NixOS/nixpkgs/";
        ref = "refs/heads/nixpkgs-unstable";
        rev = "9c3de9dd586506a7694fc9f19d459ad381239e34";
      })
      { }).schemaspy

    (import
      (builtins.fetchGit {
        # Descriptive name to make the store path easier to identify
        name = "postgresql-jdbc-42.2.20";
        url = "https://github.com/NixOS/nixpkgs/";
        ref = "refs/heads/nixpkgs-unstable";
        rev = "9c3de9dd586506a7694fc9f19d459ad381239e34";
      })
      { }).postgresql_jdbc

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "diffutils-3.8";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "cfb43ad7b941d9c3606fb35d91228da7ebddbfc5";
    }) {}).diffutils
  ];

}
