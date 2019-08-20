# How To Manage Docker Locally

There is a lot to know about Docker and how to use it.  You can always read [Docker Documentation](https://docs.docker.com/)
at the source.  However, here are some helpful tips.

## Managing the system

Docker on Mac starts you off with a 64GB disk image by default. As you use docker it may begin to complain about
having no space or running out of space. To fix this you will need to prune:

```sh
$ docker system prune
WARNING! This will remove:
        - all stopped containers
        - all networks not used by at least one container
        - all dangling images
        - all dangling build cache
Are you sure you want to continue? [y/N]
```

This will clean up your system of the typical things that cause you to run out of space. However, docker also
uses volumes to manage its storage and the volumes take up space as well. You can fix this with:

```sh
docker system prune --volumes
WARNING! This will remove:
        - all stopped containers
        - all networks not used by at least one container
        - all volumes not used by at least one container
        - all dangling images
        - all dangling build cache
Are you sure you want to continue? [y/N]
```

This free up GB of space on the disk image.
