# Organizing and Hosting MilMove Technical Docs

## Problem Statement

MilMove does not have a consistent documentation strategy. Currently, all of the following resources are being used to
house project documentation:

* Google Docs, in multiple different Google Drives
* The `mymove` repository's wiki
* Markdown files within the `mymove` codebase
* Markdown files in external repositories
* Comments within the `mymove` codebase
* Descriptions within the API specification files in `mymove`'s codebase
* Confluence

This has made locating specific documents a huge struggle. It's unclear where new docs should go, and it's impossible to
search on a topic to find relevant docs - you must know exactly what you're looking for and where. As a result, we have a
lot of duplicated documentation, unknown but useful docs, and out-of-date docs that are still being referenced.

There's also no way to establish standards or content guidelines because it is impossible to review all these docs, or
to create standards across so many different formats. Ultimately, we need a solution that can help us achieve the
following:

* The ability to search our docs based on topics. We should be able to, at the very least, create a link farm.
* The ability to review our updates and establish standards.
* A hosted resource that can be easily shared with and used by external users.
* A tool that can integrate all of our docs. We need a home base. We should ideally be able to host our API docs using
  this framework.
* The ability to involve non-engineers in the doc writing process.

Critically, it is important that we are consistent with whatever we decide upon, otherwise we will find ourselves with
just another documentation resource and the exact same problem.

## Considered Solutions

* React-based static-site generator in new repo
* React-based static-site generator in `mymove` repo
* Non-React static-site generator
* Confluence
* `mymove` wiki
* Google docs
* A little bit of everything (status quo)

## Decision Outcome

### Chosen Alternative: *React-based static-site generator in new repo*

This choice hits all of our main requirements, namely:

* Documentation updates have oversight and a well-defined review process.
* It is well-organized and searchable.
* It integrates well with other tools and resources, most importantly our API docs. The `openapi`/Redocly CLI tools also
  use React as a base.
* It can be hosted, branded, and made easily accessible to external clients.

We will, however, need to decide on a specific framework. Some good options are
[Gatsby](https://www.gatsbyjs.com/docs/), [Next.js](https://nextjs.org/docs/getting-started), and
[Docusaurus](https://docusaurus.io/docs).

Furthermore, it will take time to set up the new repo, establish standards, and migrate existing documentation over.
This would have been a challenge with any of the alternatives as well, but it is something we need to critically
evaluate. Getting design on board with the process would also be extremely valuable in making sure we are setting up a
system that is as non-engineer accessible as possible.

## Pros and Cons of the Alternatives

### Static-site generator

* `+` Articles can be written in markdown.
* `+` Can be hosted using GitHub Pages. Easy to share with external clients, devs and non-devs.
* `+` Highly customizable. Can be styled and branded to match MilMove.
* `+` Can integrate with API documentation. Can actually become an information hub.
* `+` Version control with `git`. Documentation updates must be reviewed in PRs.
* `+` PRs make standards enforceable.
* `-` Not very accessible for non-engineers to edit.
* `-` Would need to publicize the link.

#### React-based (ex: Gatsby, Next.js, Docusaurus) in new repo

* `+` The framework is a new dependency, but the environment and tech-stack will be familiar.
* `+` Putting it in an separate repo will let us update our PR requirements for documentation-specific needs. Can have
  wider permissions for non-engineers as well, making it marginally more accessible.
* `+` A new repo gives us a clean slate with regards to setup, organization, etc.
* `-` Documentation updates can't be easily incorporated into feature PRs (but PRs can be easily linked).

#### React-based in `mymove` repo

* `+` The framework is a new dependency, but the environment and tech-stack will be familiar.
* `+` Putting it in the `mymove` repo makes it easy to incorporate documentation updates into our feature PRs.
* `-` Even less accessible to non-engineers with `mymove`'s repository restrictions.
* `-` Putting it in the `mymove` repo will create an organizational nightmare.
* `-` Setting up two distinct React apps (that should be built separately and have different dependencies) in one
  project is complex and painful.
* `-` Hosting on GitHub Pages would be tricky with all the other code in the repo. Might raise security questions.

#### Non-React (ex: Jekyll) in any repo

* `-` Adds a new dependency to the project, and potentially a new tech-stack to learn. No clear benefits over a
  React-based static-site generator.

### Confluence

* `+` We have access through Jira and have to use it for other documentation, so it would help reduce our project
  dependencies.
* `+` Search functionality is decent.
* `+` Organization is decent, although there is a learning curve to setting it up effectively.
* `+` Looks nice and official. Easily able to set view-only permissions.
* `+` Easy for non-engineers to edit.
* `-` No defined process for reviewing documentation. There are alerts for when docs are changed, but they are easily bypassed.
* `-` Some standards, but they are difficult to review and enforce. Very little oversight on changes.
* `-` Cannot integrate with API docs.

### `mymove` wiki

* `+` Articles are written in markdown.
* `+` In the same repo as the codebase. Easy to reference, and searching the codebase includes wiki articles.
* `+` External devs have easy access to these docs.
* `-` Poor organization. Creating link trees and navigating from one document to another is an adhoc process and easily
  gets out of date.
* `-` No defined process for reviewing documentation. We _can_ have PRs on the wiki, but they're hard to set up, cloning
  and editing the wiki that way is awkward.
* `-` Some standards, but they are difficult to review and enforce. Very little oversight on changes.
* `-` Not very usable for external users who aren't developers. No branding.
* `-` Poor support for images.
* `-` Cannot integrate with API docs.

### Google docs

* `+` Docs can be well-formatted and readable.
* `+` Easy to add comments. Can tag folks directly on docs.
* `+` Easy for non-engineers to edit.
* `-` Not easily searchable. Practically impossible to find something unless you know the exact title.
* `-` Poor organization. Google Drive file structure is opaque and unhelpful, and we have no defined organizational
  structure.
* `-` No defined process for reviewing documentation.
* `-` Some standards, but they are difficult to review and enforce.
* `-` Hard to distribute to external folks. Very internal-focused.
* `-` Constantly a WIP. There's no uneditable version to use as a reference.
* `-` Cannot integrate with API docs.

### A little bit of everything (status quo)

* `-` Not searchable. Finding the correct document is next to impossible if you don't know which platform it's on.
* `-` No organization. New documentation is hard to promote and circulate if some folks aren't looking in that location.
* `-` No defined process for reviewing documentation.
* `-` No standards. It's impossible to enforce (or even define) content standards with so many different formats.
* `-` No set of links pointing people to the correct docs. Decent documentation is forgotten, falls out of date, and
  gets duplicated.
* `-` Nothing to point external clients/users/developers to. No (good) informational resources for outside folks.
