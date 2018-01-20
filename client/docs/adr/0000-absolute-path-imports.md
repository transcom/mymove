# Using Absolute rather than Relative Path for Imports

Imports should be as easy to use and consistent as possible across the project.

Create-react-app allows for either relative or absolute imports.
Developers should maintain consistent stylistic patterns, best set early on.

## Considered Alternatives

* Relative Paths
* Absolute Paths

## Decision Outcome

* Chosen Alternative: *Absolute Paths*
* Absolute Paths allow developers to immediately understand where imports are coming from. They also allow files to be moved without changing all the local import statements. The clarity of these imports will become more valuable should our project structure become more nested.

* Consequence: Our implementation results in adding the "NODE_PATH" variable to each import.

* Further notes: Any imports outside of the .src dir must still be imported with relative paths.

## Pros and Cons of the Alternatives

### Relative Paths

* `+` Can use without setting a NODE_PATH var
* `+` Can use for all paths--even those outside of src dir
* `-` Become chaotic and sloppy if imports traverse many dirs
* `-` Developers cannot immediately see the origin of the import
* `-` Should structure of files change, all imports have to be changed as well
* `-` Cannot copy/paste imports between files
