# Use [Pop](https://github.com/markbates/pop) as the ORM for 3M

3M will have include evolving data structures that we anticipate will change significantly throughout development.
We'd like to have a way to handle such changes with as little pain as possible. One of the things orms do well is to help with all the boilerplate when the data model changes..

Prior to the start of the contract, other Go ORMs had been explored. Pop is considered one of the more mature ORMs currently available.

## Considered Alternatives

* No ORM
* [Pop](https://github.com/markbates/pop)
* Other Go ORMs

## Decision Outcome

* Chosen Alternative: *[Pop](https://github.com/markbates/pop)*
* Pop is one of the more mature Go ORMs around
* Written by the same author as [Buffalo](https://gobuffalo.io/) framework and used in Buffalo
* Don't have to write our own SQL migrations and models; instead, we can check in the generated Fizz migrations that come with Pop
* Pop handles simple loading and saving models to the db
* Pop does not include features for doing joins /magically/ which can be a dangerous ORM feature
* We have not spent significant time looking into all our options--Pop seemed good enough for our purposes.

## Pros and Cons of the Alternatives

### No ORM

* `+` Don't have to rely on another tool over your database layer
* `-` Have to write SQL for migrations etc
* `-` Have to write models and migrations rather than having them generated.

### Other Go ORM

* `+` Don't have to write SQL for migrations etc
* `-` Not as mature as [Pop](https://github.com/markbates/pop)
* `-` Not used in [Buffalo](https://gobuffalo.io/) framework
