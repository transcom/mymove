# How to Automatically add JIRA ID to Commit Message

## Context

Jira is integrated with GitHub. As we create commits we put the ticket id in them and Jira will automatically detect that and associate the commits and branch with that ticket. This how to will show you how to enable a commit message hook to automatically strip the Jira id from your branch name and insert it at the  beginning of the commit message. If your branch doesn't match the pattern then nothing will be changed.

## Branch Format

`initials_MB-123_branch_name_description`

## Enable

To enable the automatic addition simply run the following command

```sh
ln -s ~/projects/dod/mymove/scripts/commit-msg .git/hooks/commit-msg
```

## Disable

To disable the automatic addition temporarily use the `--no-verify` flag. To disable it permanently run the following command

```sh
rm .git/hooks/commit-msg
```
