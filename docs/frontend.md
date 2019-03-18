# Front-end / React Guide

## Table of Contents

<!-- toc -->

* [Testing](#testing)
  * [Test Runners and Libraries](#test-runners-and-libraries)
  * [Writing Tests](#writing-tests)
  * [Browser Testing](#browser-testing)
* [Style](#style)
  * [Auto-formatting](#auto-formatting)
  * [Linting](#linting)
  * [File Layout & Naming](#file-layout--naming)
  * [Presentation vs. Container components](#presentation-vs-container-components)
  * [Function Declarations](#function-declarations)
  * [Ordering imports](#ordering-imports)
  * [Using Redux](#using-redux)
* [CSS](#css)
  * [BEM](#bem)
  * [USWDS](#uswds)
* [Tooling](#tooling)
  * [Sublime Plugins](#sublime-plugins)
  * [WebStorm](#webstorm)
  * [VS Code](#vs-code)
  * [vi](#vi)
  * [Browser Extensions](#browser-extensions)
* [Learning](#learning)
  * [JavaScript Concepts](#javascript-concepts)
  * [Resources](#resources)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

## Testing

### Test Runners and Libraries

* Jest - Testing framework
  * Provided by CRA, executes when you run `yarn test`.
  * Provides snapshot testing and DOM testing.
* Enzyme
  * Allows you to assert, and manipulate your rendered components with easy jQuery-like selectors. Read this nice intro guide.
  * Use Shallow rendering (`.shallow()`) as much as possible to limit the scope of testing to the component being tested and not its children.
  * Use Full rendering (`.mount()`) when you need access to component lifecycle methods.
  * Calling .debug() on a component is helpful to see what a shallow rendered component is composed of.

### Writing Tests

* React component should have a test.
  * At a minimum: does component render.
  * Container components have logic in them, and that logic should be tested.
* Redux Reducers
* Redux Action Creators?
  * TODO: Give guidance here.

### Browser Testing

* We use the [Cypress framework](https://www.cypress.io/) for most browser testing, both with chrome and headless chrome
* For testing on Windows 10 with IE 11 we have a [testing document](https://docs.google.com/document/d/1j04tGHTBpcdS8RSzlSB-dImLbIxsLpsFlCzZUWxUKxg/edit#)

## Style

Adhere to Airbnb's [JavaScript Style Guide](https://github.com/airbnb/javascript) unless they conflict with the project’s Prettier or Lint rules.

### Auto-formatting

* Prettier
  * Prefer single quotes for non-JSX code (CLI: `--single-quote` API: `singleQuote: true`)
  * Prefer trailing commas for cleaner PRs and error reduction (CLI: `--trailing-comma true` API: `trailingComma: true`)
  * A `.prettierrc` file is in the project for the above settings.
  * Make sure to [set up your editor](https://prettier.io/docs/en/editors.html) to format (and possibly autosave) with Prettier with the above configurations. You will need to install Prettier globally for this.

### Linting

* CRA runs ESLint on the dev server. We are using [create-app-rewired](https://github.com/timarney/react-app-rewired) to configure eslint to use a security package requested by the DOD.

### File Layout & Naming

* All front-end client code is kept within a subdirectory called `src`. This is an artifact of using `create-react-app`.
* Inside that directory:
  * `/src`
  * `/src/scenes` Group components by scene name
  * `/src/shared` Group shared components, like headers
* File naming
  * All component files should be named in `PascalCase`, component names should match the file names (Exception: Higher Order Components are named in `camelCase`)
  * Other files should be in `camelCase`
  * Component files should use the `.jsx` file extension
  * If there are multiple components for a feature, they should be in a folder with the primary component in a file named `index.jsx`.

### Presentation vs. Container components

* See this [blog post](https://medium.com/@dan_abramov/smart-and-dumb-components-7ca2f9a7c7d0), and this [GitHub gist](https://gist.github.com/chantastic/fc9e3853464dffdb1e3c). Personally I found the gist to convey the idea faster.
* The gist (ha!) of it is: React components should either have styling or logic, but not both.
* Presentational components should be declared with plain functions, not fat arrow functions.

### Function Declarations

* Use plain functions for stateless components and React component lifecycle methods. Use fat arrow functions for other class methods because it ensures the scope of the function will be the declaring component.
* Never create new functions in the render method (or return value for a stateless component). Functions should either be declared directly on a class, imported, or received as a prop.

### Ordering imports

* Imports should go in this order; group like with like.
  * React and Redux imports, anything primary to the framework
  * Other packages
  * Component imports
  * CSS files

### Using Redux

* Connect higher level components to Redux, pass down props to less significant children. (Avoid connecting everything to Redux.)
* Use [ducks](https://github.com/erikras/ducks-modular-redux) for organizing code.

## CSS

### BEM

* Where we need to write CSS, follow the BEM naming convention to increase readability & reusability.
  * BEM is short for Block, Element, Modifier which are the three components of classnames.
  * From [CSS Tricks](https://css-tricks.com/bem-101/): "In this CSS methodology a block is a top-level abstraction of a new component, for example a button: `.btn { }`. This block should be thought of as a parent. Child items, or elements, can be placed inside and these are denoted by two underscores following the name of the block like `.btn__price { }`. Finally, modifiers can manipulate the block so that we can theme or style that particular component without inflicting changes on a completely unrelated module. This is done by appending two hyphens to the name of the block just like `btn--orange`."
  * Expanding on this, a modified child class would have a class name like `.btn__price--orange`.

### USWDS

* Check the [USWDS Design Standards](https://standards.usa.gov/components/) for a component that matches your needs. Maximize the code view to see what classes to use to replicate the component styles.
* USWDS has a [Slack chat](https://chat.18f.gov/) you can go to for help. Get invited to it by filling out this form.

## Tooling

If you are using Sublime, Webpack, or VS Code, you may want to install plugins to support the following:

* Prettier
* ESLint
* React

Below are some suggestions for plugins. However, to get the plugins to work, you may need to install prettier and ESLint globally. You will have to make sure they are kept up to date with the project.

### Sublime Plugins

* PackageControl
* EditorConfig
* JsPrettier (you will need to configure it to auto-format on save)
* Babel (for JSX syntax--though looking for better option)
* Git

### WebStorm

Has plugins for most out the box, but setting up Prettier is tricky. See [the documentation](https://prettier.io/docs/en/webstorm.html).

### VS Code

* Prettier
* Path Intellisence

### vi

* [vim-prettier](github.com:prettier/vim-prettier)
* [vim-javascript](pangloss/vim-javascript.git)
* [editorconfig](editorconfig/editorconfig-vim.git)

### Browser Extensions

Install the following extensions to assist with debugging React and Redux applications:

* [React Developer Tools](https://github.com/facebook/react-devtools#installation)
* [Redux DevTools Extension](http://extension.remotedev.io/#redux-devtools-extension)

## Learning

### JavaScript Concepts

Important JS patterns and features to understand.

* Destructuring Assignment
  * [A Dead Simple Intro to Destructuring JavaScript Objects](http://wesbos.com/destructuring-objects/)
* Fat Arrow Functions
  * [ES5 Functions vs ES6 Fat Arrow Functions](https://medium.com/@thejasonfile/es5-functions-vs-es6-fat-arrow-functions-864033baa1a)
* Higher Order Components
  * [Higher Order Components: A React Application Design Pattern](https://www.sitepoint.com/react-higher-order-components/)
* Promises
  * [An incremental tutorial on promises](https://www.sohamkamani.com/blog/2016/08/28/incremenal-tutorial-to-promises/)
* Spread Operator/Rest Params
  * [JavaScript & The Spread Operator](https://codeburst.io/javascript-the-spread-operator-a867a71668ca)
  * [How Three Dots Changed JavaScript](https://dmitripavlutin.com/how-three-dots-changed-javascript/)
* Template Literals
  * [Template Literals](https://css-tricks.com/template-literals/)

### Resources

Various resources on React, Redux, etc, for a variety of learning styles.

* _Read_: [React Tutorial](https://reactjs.org/tutorial/tutorial.html) - Official tutorial from React. I (Alexi) personally found this cumbersome. If you stick with it you’ll learn the basics.
* _Watch_: [Getting Started with Redux](https://egghead.io/courses/getting-started-with-redux) - Free 30 video series by the author of Redux.
* _Watch_: [ReactJS / Redux Tutorial](https://www.youtube.com/playlist?list=PL55RiY5tL51rrC3sh8qLiYHqUV3twEYU_) - ~60 minutes of YouTube videos that will get you up and running with React and Redux. The content is useful, the guy’s voice can be a bit of a challenge.
* _Watch_: [This video](https://www.youtube.com/watch?list=PLb0IAmt7-GS188xDYE-u1ShQmFFGbrk0v&v=nYkdrAPrdcw) from the introduction of Flux can be useful for some high-level background about the pattern (the MVC bashing is overdone, but otherwise this video is useful.)
* _Do_: Roll your own React app! Make a little project of your own. This works well if you’re more hands-on. Here are some rough steps, but you’ll need to do a bit of filling-in-the-blanks:
  * Use [create-react-app](https://github.com/facebookincubator/create-react-app) to bootstrap a new React project.
  * Figure out how to run the app live (hint: yarn start)
  * Find and skim through some of the important files it made: `index.hmtl`, `index.js`, `App.js`. What do these look like they’re doing?
  * Change the page title to something of your choosing.
  * Create a new React [component](https://reactjs.org/docs/react-component.html) that has a `<button>` or something in it.
  * [import](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import) that component into `App.js`, and make sure you can see it!
  * Write a new test for your component. (Hint: `yarn test`). create-react-app gives you Jest for free, look at its manual.
  * Make the thing in your component clickable, even if it just does `alert(‘hey there!’)`
  * Add [React Router](https://github.com/ReactTraining/react-router) to your project.
  * Make a new component like the first one, and add routes so that they display depending on the URL. E.g:
    * `http://milmovelocal:3000/component1` shows the first one on the page;
    * `http://milmovelocal:3000/component2` shows the second one.
  * Add [Redux](https://redux.js.org/) to your project.
    * This is a rather big step. You’ll need to have some sort of state, so make a login button and “logged in” will be the state you are going to keep track of.
    * When the user is logged in, there should be a “log out” button shown.
