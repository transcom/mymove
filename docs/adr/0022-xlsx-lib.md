# Chose Excelize package to parse XLSX files

**User Story:** *[Pivotal Task](https://www.pivotaltracker.com/story/show/156427513)*

The Rate Engine relies on data from the 400NG Tariff baseline rates excel doc which is published yearly. We need a way to import this data into the application to use in look ups for the rate engine.
Importing this from the excel document is the most reliable and reproducible way to get the data in the database.

## Considered Alternatives

* *[XLSX](https://github.com/tealeg/xlsx)*
* *[Excelize](https://github.com/360EntSecGroup-Skylar/excelize)*

## Decision Outcome

* Chosen Alternative: *Excelize*
* Both alternatives seem to have similar base functionality, but this one has better documentation, more engaged contributors, and the other top hit refers to this as worth checking out. Because we're using for a very discrete task (parsing and then hoovering up data for our database), I think having a package that has a relatively limited scope is appropriate. Used [this](https://docs.google.com/document/d/1Z_R9mFRo4n-rvxLD1Cbz-KXy32H5Z50_J3j4x15hq0Y/edit) as guidance for how to make this decision.

## Pros and Cons of the Alternatives

### *XLSX*

* `+` *Most used on [GoDocs](https://godoc.org/github.com/tealeg/xlsx)*
* `+` *Easy to get started
* `+` *Has some nice offshoots such as an xlsx to csv project and a streaming xlsx parser*
* `+` *Has CI and decent test coverage*
* `-` *Less in-depth docs*
* `-` *Less frequently updated*
* `-` *Long waits on PRs*

### *Excelize*

* `+` *Second most frequently used on [GoDocs](https://godoc.org/github.com/360EntSecGroup-Skylar/excelize)*
* `+` *Easy to get started
* `+` *Has CI and excellent test coverage*
* `+` *Better documentation*
* `-` *Fewer references on the internet that I've found*
