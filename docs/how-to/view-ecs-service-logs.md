# How to View ECS Service Logs

To view service logs you'll want to use the `ecs-service-logs` binary. You can build it with `make bin/ecs-service-logs`.
Here are some examples of usage:

For running containers:

```sh
ecs-service-logs show --cluster app-staging --service app --git-branch "placeholder_branch_name" --status "RUNNING" --verbose
```

You can even narrow it down to the specific commit:

```sh
ecs-service-logs show --cluster app-staging --service app --git-branch "placeholder_branch_name" --git-commit "git_commit_hash" --status "RUNNING" --verbose
```

If you want to see logs from something that is stopped you can change to `--status "STOPPED"` in the above commands.

You can do a lot more with this command. Check out the rest of the information in the
[ecs-service-logs README.md](https://github.com/transcom/mymove/blob/master/cmd/ecs-service-logs/README.md).
