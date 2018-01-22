# Using Absolute rather than Relative Path for Imports

Imports should be as easy to use and consistent as possible across the project.

Create-react-app allows for either relative or absolute imports.
Developers should maintain consistent patterns, best set early on.

## Considered Alternatives

* *Relative Paths*: depending on where in the app the import is, to get to the same file (feedback.jsx) could look like either of the following
  * `./feedback.jsx`
  * `../../../src/scenes/feedback/feedback.jsx`
* *Absolute Paths*: For both of the above examples, the absolute path is the same
  * `scenes/feedback/feedback.jsx`
  * **Implementation**:
    * Add `NODE_PATH` to .env file--this file contains environment variables that are sourced into the shell that runs the app. Typically, it's not checked into Git. However, in this case, it would be, because we would want all devs to be able to use this same variable. So we'd be somewhat altering the .env file's intended application.
    * Add `NODE_PATH` to build scripts in `packages.json`. Avoid using the .env file altogether, and set the necessary variable before each script is run.
* *Combination*: Use whichever option is cleanest and easiest to understand.
  * `./feedback.jsx`
  * `scenes/feedback/feedback.jsx`

## Decision Outcome

* Chosen Alternative: *Combination*
* Absolute Paths allow developers to immediately understand where imports are coming from. They also allow files to be moved without changing all the local import statements. The clarity of these imports will become more valuable should our project structure become more nested.
* Meanwhile, using relative paths allows developers not to be dogmatic when importing a module from within the same dir, and acknowledges that all imports outside of the .src dir requier relative paths.

* Consequence: Our implementation sets the NODE_PATH environment variable before each build script in packages.json

## Pros and Cons of the Alternatives

### Absolute Paths

* `+` Are consistent no matter the complexity of relationship between import and module
* `+` Developers can immediately see the origin of the import
* `+` Should structure of files change, all imports can remain the same
* `+` Can copy/paste imports between files
* `-` Cannot use without setting a NODE_PATH var
* `-` Cannot use for all paths--such as those outside of src dir
* `-` In cases of import from a close module, unnecessarily lengthy

### Relative Paths

* `+` Can use without setting a NODE_PATH var
* `+` Can use for all paths--even those outside of src dir
* `-` Become chaotic and sloppy if imports traverse many dirs
* `-` Developers cannot immediately see the origin of the import
* `-` Should structure of files change, all imports have to be changed as well
* `-` Cannot copy/paste imports between files
