# use ./nix/update.sh to install
#
# use <https://ahobson.github.io/nix-package-search> to find a package version

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
      rev = "253aecf69ed7595aaefabde779aa6449195bebb7";
    }) {}).bash_5

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "nodejs-14.18.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "14ee52edff84f16f6268ebb9f87380cd86c433da";
    }) {}).nodejs-14_x

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "yarn-1.22.11";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "253aecf69ed7595aaefabde779aa6449195bebb7";
    }) {}).yarn

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "go-1.17.5";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "0acfd0c1e179c1fa276fd931bdd22f207e7d5a48";
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
      rev = "9c3de9dd586506a7694fc9f19d459ad381239e34";
    }) {}).pre-commit

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "circleci-cli-0.1.15663";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "23cedc3088a628e1f5454cab6864f9b1a059e1ba";
    }) {}).circleci-cli

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "jq-1.6";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "253aecf69ed7595aaefabde779aa6449195bebb7";
    }) {}).jq

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "shellcheck-0.7.2";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "391f93a83c3a486475d60eb4a569bb6afbf306ad";
    }) {}).shellcheck

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "opensc-0.22.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "253aecf69ed7595aaefabde779aa6449195bebb7";
    }) {}).opensc

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "entr-5.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "253aecf69ed7595aaefabde779aa6449195bebb7";
    }) {}).entr

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "aws-vault-6.3.1";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "253aecf69ed7595aaefabde779aa6449195bebb7";
    }) {}).aws-vault

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "watchman-4.9.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "253aecf69ed7595aaefabde779aa6449195bebb7";
    }) {}).watchman

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "awscli2-2.2.14";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "14b0f20fa1f56438b74100513c9b1f7c072cf789";
    }) {}).awscli2

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "chamber-2.10.2";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "725ef07e543a6f60b534036c684d44e57bb8d5de";
    }) {}).chamber

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "nss-cacert-3.66";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "e579e93b6517c29af8b84641d2ec16a064d34c8b";
    }) {}).cacert

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "hadolint-2.8.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "8afc4e543663ca0a6a4f496262cd05233737e732";
    }) {}).hadolint

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "go-swagger-0.29.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "5efc8ca954272c4376ac929f4c5ffefcc20551d5";
    }) {}).go-swagger

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "schemaspy-6.1.0";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "9c3de9dd586506a7694fc9f19d459ad381239e34";
    }) {}).schemaspy

    (import (builtins.fetchGit {
      # Descriptive name to make the store path easier to identify
      name = "postgresql-jdbc-42.2.20";
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "refs/heads/nixpkgs-unstable";
      rev = "9c3de9dd586506a7694fc9f19d459ad381239e34";
    }) {}).postgresql_jdbc

];

  # the pre-commit hooks expects the binary to be `circleci`
  postBuild = ''
  ln -s $out/bin/circleci-cli $out/bin/circleci
  '';
}
