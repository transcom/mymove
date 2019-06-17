# How to Deploy to Experimental

Experimental is the MilMove environment that can be used to test out a code change in a branch in a safe way. This is especially good if you'd like Product, Design, or the client to test out the change before you merge the code in. It's also a good decision to use experimental to test out particularly risky changes like changes to containers, app startup, connection to data stores, secure migrations and data loads.

## How to do it

To deploy to experimental, you'll need to trigger the CircleCI workflow to include the deploy to experimental steps.
Edit [the config file](https://github.com/transcom/mymove/blob/master/.circleci/config.yml) by replacing all instances of `placeholder_branch_name` with your branch name. You'll also want to uncomment any of these lines (and related ones), like this:

```sh
# if testing on experimental, you can disable these tests by using the commented block below.
  filters:
    branches:
      ignore: your_branch_name_goes_here
```

Best practices mean you should announce deploys to experimental in the `#dp3-experinental-env` slack channel at least 20 minutes before you intend to deploy. Try to get a 👍 from someone who is commonly using experimental. If no one comments in that time-frame feel free to deploy. Only one person can use experimental at the same time. Any deploys will overwrite what is currently on experimental. When you push up to Github, it'll trigger the CircleCI workflow to start immediately. You can view its progress on [CircleCI's UI](https://circleci.com/gh/transcom/workflows/mymove) and clicking on your branch. If it has succeeded, then it should be available immediately on [SM experimental](https://my.experimental.move.mil/), [Office experimental](https://office.experimental.move.mil/), and [TSP experimental](https://tsp.experimental.move.mil/).

## I've got a server-side feature flag

You'll need to add it to [experimental config file](https://github.com/transcom/mymove/blob/master/config/env/experimental.env). In [the container config file](https://github.com/transcom/mymove/blob/master/config/app.container-definition.json) you'll need to add it by using mustache syntax.

## Don't forget to clean up your branch for your pr review

All of the [how to do it instructions](#How-to-do-it) need to be reverted so all future branches don't get deployed to experimental. Also, be sure to announce in `#dp3-experimental-env` that experimental is available for someone else to use!
