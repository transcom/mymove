# Configuration

## env files

The `/env/*.env` files define how the ECS containers are configured. They are used to set non-secret environment
variables. These files are read in when using the `ecs-deploy task-def` tool. For example, to view a task definition:

```sh
go run github.com/transcom/mymove/cmd/ecs-deploy task-def \
  --aws-account-id "${AWS_ACCOUNT_ID}" \
  --aws-region us-west-2 \
  --service app \
  --environment experimental \
  --image ${AWS_ACCOUNT_ID}.dkr.ecr.us-west-2.amazonaws.com/app:git-e2b6c625368d05b9bc24a5a58a04350278d31ad9 \
  --variables-file config/env/experimental.app.env \
  --entrypoint "/bin/milmove serve" \
  --dry-run
```

## TLS cert/key (optional)

The `devlocal-https.*` files are a self-signed TLS cert/key pair. They are a [snake oil](https://en.wikipedia.org/wiki/Snake_oil_(cryptography)) certificate used to locally run the webserver during development. They are included as a convenience so engineers don't have to generate their own.
