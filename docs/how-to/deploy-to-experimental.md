# How to Deploy to Experimental

Experimental is the MilMove environment that can be used to test out a branch out. This is especially good if you'd like Product, Design, or the client to test out a change before you merge the code in. It's also a good decision to use experimental to test out particularly risky changes.

## How to do it

To deploy to experimental, you'll need to trigger the CircleCI workflow to include the deploy to experimental steps.
Edit [the config file](https://github.com/transcom/mymove/blob/master/.circleci/config.yml) by replacing all instances of `placeholder_branch_name` with your branch name. You'll also want to uncomment any of these lines (and related ones), like this:

```sh
# if testing on experimental, you can disable these tests by using the commented block below.
  filters:
    branches:
      ignore: your_branch_name_goes_here
```

When you push up to Github, it'll trigger the CircleCI workflow to start. You can view it's progress on [CircleCI's UI](https://circleci.com/gh/transcom/workflows/mymove) and clicking on your branch. If it's succeeded, then it should be up on [SM experimental](https://my.experimental.move.mil/), [Office experimental](https://office.experimental.move.mil/), and [TSP experimental](https://tsp.experimental.move.mil/). Best practices mean you should announce deploys to experimental in the `#dp3-engineering` slack channel, so only one person is trying to use it at the same time.

## I've got a server-side feature flag

You'll need to add it to [the container definition](https://github.com/transcom/mymove/blob/master/config/app.container-definition.json) in the `environment` array with the name of the server-side flag and the value you'd like it to be on experimental.

## Don't forget to clean up your branch for your pr review

All of the [how to do it instructions](#How-to-do-it) need to be reverted so all future branches don't get deployed to experimental. Also, be sure to announce in `#dp3-engineering` that experimental is available for someone else to use!
