# Programming Guide

The intention for this document is to share our collective knowledge on best practices and allow everyone working on the DOD MyMove project to write code in compatible styles.

If you are looking to understand choices made in this project, see the list of [ADRs](https://github.com/transcom/mymove/tree/master/docs/adr).

## Frontend / React

### Testing

#### Test Runners and Libraries

* Jest - Testing framework
  * Provided by CRA, executes when you run `yarn test`.
  * Provides snapshot testing and DOM testing.
* Enzyme
  * Allows you to assert, and manipulate your rendered components with easy jQuery-like selectors. Read this nice intro guide.
  * Use Shallow rendering (`.shallow()`) as much as possible to limit the scope of testing to the component being tested and not its children.
  * Use Full rendering (`.mount()`) when you need access to component lifecycle methods.
  * Calling .debug() on a component is helpful to see what a shallow rendered component is composed of.

#### Writing Tests

* React component should have a test.
  * At a minimum: does component render.
  * Container components have logic in them, and that logic should be tested.
* Redux Reducers
* Redux Action Creators?
  * TODO: Give guidance here.

### Style

Adhere to AirBnB's [Javascript Style Guide](https://github.com/airbnb/javascript) unless they conflict with the project’s Prettier or Lint rules.

#### Auto-formatting

* Prettier
  * Prefer single quotes for non-JSX code (CLI: `--single-quote` API: `singleQuote: true`)
  * Prefer trailing commas for cleaner PRs and error reduction (CLI: `--trailing-comma true` API: `trailingComma: true`)
  * A `.prettierrc` file is in the project for the above settings.
  * Make sure to [set up your editor](https://prettier.io/docs/en/editors.html) to format (and possibly autosave) with Prettier with the above configurations. You will need to install Prettier globally for this.

#### Linting

* CRA runs ESLint on the dev server. No additional configuration is available unless the app is ejected.

#### References

* File naming
  * All component files should be named in PascalCase, component names should match the file names (Exception: Higher Order Components are named in camelCase)
  * Other files should be in camelCase
  * Component files should use the .jsx file extension
  * If there are multiple components for a feature, they should be in a folder with the primary component in a file named `index.jsx`.

#### File layout

* All frontend client code is kept within a subdirectory called ‘client’.
* Inside that directory:
  * `client/`
  * `client/src`
  * `client/src/page` Group components by page name
  * `client/src/shared` Group shared components, like headers

#### Presentation vs. Container components

* See this blog post, and this GitHub gist. Personally I found the gist to convey the idea faster.
* The gist (ha!) of it is: React components should either have styling or logic, but not both.
* Presentational components should be declared with plain functions, not fat arrow functions.

#### Function Declarations

* Use plain functions for stateless components and React component lifecycle methods. Use fat arrow functions for other class methods because it ensures the scope of the function will be the declaring component.
* Never create new functions in the render method (or return value for a stateless component). Functions should either be declared directly on a class, imported, or received as a prop.

#### Ordering imports

* Imports should go in this order; group like with like.
  * React and Redux imports, anything primary to the framework
  * Other packages
  * Component imports
  * CSS files

#### Using Redux

* Connect higher level components to Redux, pass down props to less significant children. (Avoid connecting everything to Redux.)
* Use [ducks](https://github.com/erikras/ducks-modular-redux) for organizing code.

### CSS

#### BEM

* Where we need to write CSS, follow the BEM naming convention to increase readability & reusability.
  * BEM is short for Block, Element, Modifier which are the three components of classnames.
  * From [CSS Tricks](https://css-tricks.com/bem-101/): “In this CSS methodology a block is a top-level abstraction of a new component, for example a button: .btn { }. This block should be thought of as a parent. Child items, or elements, can be placed inside and these are denoted by two underscores following the name of the block like .btn__price { }. Finally, modifiers can manipulate the block so that we can theme or style that particular component without inflicting changes on a completely unrelated module. This is done by appending two hyphens to the name of the block just like btn--orange.”
  * Expanding on this, a modified child class would have a class name like .btn__price--orange.

#### USWDS

* Check the [USWDS Design Standards](https://standards.usa.gov/components/) for a component that matches your needs. Maximize the code view to see what classes to use to replicate the component styles.
* USWDS has a [Slack chat](https://chat.18f.gov/) you can go to for help. Get invited to it by filling out this form.
