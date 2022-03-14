# Use a consistent file structure for front-end code

Currently, the front end code is in a state of mid-reorganization from one file structure to another. However, without a specific plan in place, this reorganization risks losing momentum, and over time could result in the code being indefinitely split between two different file structures. This makes future work on this codebase more confusing and increases the risk for mistakes to be made. Therefore, the purpose of this ADR is to establish a plan and certain guidelines for moving the reorganization forward, hopefully at an increased pace.

## Considered Alternatives

- Actively move code into the new structure as it is worked on
- Move everything into the new structure at once
- Do nothing, continue work in two different structures

## Decision Outcome

_Chosen Alternative:_ Actively move code into the new structure as it is worked on

For the purpose of distinguishing between the two file structures, I will use **"legacy"** to mean the previous structure (in which all files were located in `src/shared` or `src/scenes`), and **"new"** to mean the new structure.

For reference, the new file structure is described [here](https://transcom.github.io/mymove-docs/docs/dev/contributing/frontend/frontend#file-layout--naming).

I am proposing we start operating under the following guidelines:

- Any new front-end files should be created in one of the new directories.
  - If you are creating a new file and feel that it doesn't fit in one of the existing new directories, _or_ you aren't sure where it belongs, ask in #prac-frontend.
- If you are making significant changes to a file in a legacy directory, take the opportunity to move it into one of the new directories.
  - For moving files, usually you can do a global search of the codebase to find where it is imported to other files, searching by the filename and/or its named exports. Some IDEs (including VSCode) will also offer to automatically update import paths for you.
  - While this can be a delicate process, it is also usually easily tested by building the application and viewing some implementation of the code that has been moved. The build process and Jest tests will both throw an error if a missing file is referenced somewhere. If you don't feel comfortable doing this though, or if you would like help, ask in #prac-frontend.
- When creating or editing React components, they should be structured so they are encapsulated in their own directory, and contain all related files (the component source code, tests, stories, CSS, etc.).
  - For example, this would be the file structure for a component called `MyComponent`:
    - `MyComponent/`
      - `MyComponent.jsx` (component source code)
      - `MyComponent.test.jsx` (unit tests, will be run in Jest)
      - `MyComponent.stories.jsx` (Storybook stories, UI test cases)
      - `index.jsx` (optional - used to connect component code with required providers such as Redux or the `withRouter` HOC)
      - `MyComponent.module.scss` (optional, SCSS module code related to the component)
  - Component files should be PascalCased (following the naming convention of the exported component itself).
  - Some components may composed of multiple, smaller components. Sometimes it will make sense for all of these components to be defined in the same file, or for them to be in multiple files within the one component directory (especially if they are only ever used within that one component. That's okay too!
  - These are guidelines, not strict rules, and there may be exceptions. Use your best judgment and always ask in #prac-frontend if you aren't sure or would like another opinion!
- Front-end files that are _not_ React components (such as JS helpers, utilities, constants, etc.) can just be organized within a top-level directory (such as `src/helpers`), and should be camelCased.
- Most of the time, all JS files should also have a corresponding test file.

## Pros and Cons of the Alternatives

### Move everything into the new structure at once

- `+` All of the frontend code will be organized according to the new file structure right away
- `+` We will no longer be in between two ways of organizing files
- `-` Such a significant change to the code would be difficult to resolve promptly, and would be prone to conflicts and introducing bugs

### Do nothing, continue work in two different structures

- `+` No changes needed to our current process
- `-` It is unclear to the team where new files should be created, and when to move existing files
- `-` It may become more and more confusing as time passes and team members lose context for this work
