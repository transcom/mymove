# How to revert a change

## On Staging app

Go to Github and find the pr with the change you'd like to revert. Click the revert button. This will generate a pr for you that reverts the change. You'll still need a reviewer to approve the pr.

Note that this will only revert code changes (and not migration changes), so you'll need to consider if reverting the code, but not the db migrations will cause any problems.

Once CI passes, it should automatically deploy to staging.

## On Production app

Still to be discovered...
