# Front-end / React Guide

## Table of Contents

<!-- toc -->

* [Design + Engineering Process for new components](#design--engineering-process-for-new-components)
  * [Design delivers component design](#design-delivers-component-design)
  * [Engineering](#engineering)
  * [Update Loki tests accordingly](#update-loki-tests-accordingly)
* [Testing](#testing)
  * [Test Runners and Libraries](#test-runners-and-libraries)
  * [Writing Tests](#writing-tests)
  * [Browser Testing](#browser-testing)
  * [Storybook Testing](#storybook-testing)
* [Code Style](#code-style)
  * [Auto-formatting](#auto-formatting)
  * [Linting](#linting)
  * [File Layout & Naming](#file-layout--naming)
  * [Presentation vs. Container components](#presentation-vs-container-components)
  * [Function Declarations](#function-declarations)
  * [Ordering imports](#ordering-imports)
  * [Using Redux](#using-redux)
  * [CSS Styling Standards](#css-styling-standards)
    * [Using Sass and CSS Modules](#using-sass-and-css-modules)
    * [Classnames](#classnames)
    * [rem vs. em](#rem-vs-em)
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

Regenerate with "pre-commit run -a markdown-toc"

<!-- tocstop -->

## Design + Engineering Process for new components

MilMove has defined a process for taking a new component from concept to design to implementation. This section of the doc will describe this process. We use [Storybook](https://storybook.js.org/) for showing the finished components and you can view all current ones on master by going to our [public storybook site](https://storybook.move.mil/). If you want to see things locally please check out the [How To Run Storybook](how-to/run-storybook.md) document.

### Design delivers component design

After the research and initial prototypes are made a designer will create a full design for a new component, card, or page. Once the design has passed the design team's review process the designer will deliver a link to the [Abstract](https://www.abstract.com/) design. Since engineers are not likely to have an Abstract account the designers will ensure that this link is a publicly viewable version. For example here is the link we used for the [TabNav](https://app.abstract.com/share/39907fe2-a5c6-4063-ac68-71bae522e296?mode=build&selected=3210965808-139C6AE4-167B-4B24-B583-C1F45CC3493D) component.

We have added the github `@transcom/truss-design` as code owners of `src/stories` thus requiring their approval for these changes in addition to normal engineering review.

### Engineering

Once an engineer has the Abstract design for a new component they can begin to implement it. The new process requires that all components have a [Storybook](https://storybook.js.org/) story created or updated for it. Storybook stories require approval from someone on the design team before they can be merged, preferable the designer who created the original Abstract design. We are following the [USWDS](#uswds) standard for design and implementation, so please review that section of this document. Be sure to use [USWDS mixins](https://designsystem.digital.gov/utilities/) and any components that are available in `react-uswds`. If there is a USWDS component not already in `react-uswds` please add it to that package and then make use of it.

### Update Loki tests accordingly

We currently use [Loki](https://loki.js.org/) for ensuring our storybook components do not regress as the project goes on. Please ensure you run the tests and add or update new reference images as you create or update components. See [How to Run Loki tests against Storybook](how-to/run-loki-tests-against-storybook.md) document for more details.

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

### Storybook Testing

* We use the [Loki](https://loki.js.org/) package for visually testing storybook.
* For details on how to run, add, or update these tests see [How to Run Loki tests against Storybook](how-to/run-loki-tests-against-storybook.md)

## Code Style

Adhere to Airbnb's [JavaScript Style Guide](https://github.com/airbnb/javascript) unless they conflict with the project’s Prettier or Lint rules.

### Auto-formatting

* Prettier
  * Prefer single quotes for non-JSX code (CLI: `--single-quote` API: `singleQuote: true`)
  * Prefer trailing commas for cleaner PRs and error reduction (CLI: `--trailing-comma true` API: `trailingComma: true`)
  * A `.prettierrc` file is in the project for the above settings.
  * Make sure to [set up your editor](https://prettier.io/docs/en/editors.html) to format (and possibly autosave) with Prettier with the above configurations. You will need to install Prettier globally for this.
  * We currently pin the prettier dependency to a specific version to avoid frequent formatting churn.

### Linting

* CRA runs ESLint on the dev server, you can execute `yarn run lint` to execute linting on all files, otherwise pre-commit will run it for you on files changed.
* We are using [rescripts](https://github.com/harrysolovay/rescripts) to configure eslint to use a security package, [eslint-plugin-security](https://github.com/nodesecurity/eslint-plugin-security), requested by the DOD.

### File Layout & Naming

* All front-end client code is kept within a subdirectory called `src`. This is an artifact of using `create-react-app` and common React best practice.
* Inside that directory:
  * `/src/components` Low-level React Components that are more about rendering UI than handling application logic, and should typically not be connected to providers directly. Aim for abstract, generic components that can be shared across the application easily, and will usually have corresponding stories files for viewing in Storybook.
  * `/src/config` High-level configuration definitions relevant to the whole application (such as routes).
  * `/src/constants` Define constants here instead of using string literals for any values with specific meaning in the context of the application. For example, data that comes back from the API that may be used in UI logic (such as a user role or payment request status).
  * `/src/containers` React Components that are primarily concerned with connecting UI to containers or providers (such as Redux), and sharing behavior or patterns via hooks or higher-order components.
  * `/src/helpers` Miscellaneous utilities that implement logic, data handling, and other common functions used throughout the application. These should not include React-specific code such as JSX, and they should generally be purely functional and well-tested.
  * `/src/layout` React components used to render common layout elements, such as header, footer, page content, etc. Similar to the components located in /src/components, they should focus on rendering UI rather than application logic or connecting to providers. However, they are designed such that there should only ever be one instance on each page.
  * `/src/pages` React components that correspond to actual routes (URLs). These are responsible for assembling the UI components for a page, and hooking them up with the necessary providers such as Redux. Queries should be co-located with page components, since pages are explicitly dependent on them.
  * `/src/stories` Storybook stories for components live here.
  * ***NOTE: The code style recommendations above are strictly enforced in the above directories***
  * `/src/shared/styles` Global or shared styles
* Previous layout of components PPM and HHG work was done in the following structure and will remain concurrently until it can be migrated to the new recommendations above. ***No new files should go in the following directories.*** New files should be put into the above structure. If there are significant changes to components in the below directories please migrate them to the new structure.
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

### CSS Styling Standards

MilMove is transitioning from anarchistic styling to more organized and standardized styling, so much of the existing code is not yet organized to the current standards.  You can find an example of refactored code styling of `InvoicePane.jsx` in `InvoicePanel.module.scss` and its child components and corresponding stylesheets.  All new components/styling should utilize the below standards. When we touch an existing component, we should try to adjust the styling to follow the standards.

#### Using Sass and CSS Modules

* All new components should utilize [Sass](https://sass-lang.com/documentation/file.SASS_REFERENCE.html) and [CSS modules](https://github.com/css-modules/css-modules) (See [ADR](https://github.com/transcom/mymove/blob/master/docs/adr/0031-css-tooling.md)).  To apply Sass and CSS Modules, name files  with this syntax: `<component>.module.scss`
* Use CSS Modules' [`composes`](https://bambielli.com/til/2017-08-11-css-modules-composes/#) to build a new class out of other pre-defined classes
* Where to put styles?
  * For global styles, such as colors and themes, utilize global variables in files such as `colors.scss`
  * Styles shared across the app should go into `src/shared/styles`.  Most common styles can fit into `common.module.scss`, though additional files may be appropriate.
  * Styles unique to a component should go into a corresponding component style file (`<component>.module.sccss`)
  * Sibling components that share styles: share styles through the parent component's stylesheet (ex. `StorageInTransitApproveForm.jsx` and `StorageInTransitDenyForm.jsx` could share styles that would go into their parent `StorageInTransit.jsx`'s stylesheet `StorageInTransit.module.scss`)
* Syntax
  * Import styles from a component's stylesheet as `import styles from 'InvoicePanel.module.scss'`.  If more than one stylesheet needs to be imported, use `styles` for the component's style, and another name for the secondary stylesheets (ex. `iconStyles`)
  * Access the styles with dot notation `styles.myclassname`.  If the class name uses a hyphen, like `invoice-panel`, it must be accessed with this notation: `styles['invoice-panel']`
  * If fewer than 50% of the styles are used from a stylesheet, import only the styles used (ex. `import { myclassname } from 'MyComponent.module.scss'`)
* Tests
  * Add `data-cy` as an attribute on the elements you want to identify in your tests, and use that attribute instead of the class name to identify an element in your test (see [Cypress best practices](https://docs.cypress.io/guides/references/best-practices.html#Selecting-Elements)).  CSS modules tags on a unique identifier to create local classes.  This means that tests that use just the class name will fail.
  * Ex. Rather than `cy.get('.invoice-detail')`, use `cy.get('[data-cy="invoice-detail"]')` when your element looks like `<div classname=styles["invoice-detail"] data-cy="invoice-detail"></div>`
* Resources
  * [Beginners Guide to Sass](https://coolestguidesontheplanet.com/guide-beginning-sass-css/)
  * CSS Modules with React: [The Complete Guide](https://blog.yipl.com.np/css-modules-with-react-the-complete-guide-a98737f79c7c)

#### Classnames

* Use [`classnames`](https://github.com/JedWatson/classnames) package for assigning classes based on boolean values

#### rem vs. em

Understand the [difference between rem and em](https://zellwk.com/blog/rem-vs-em/). Which one you use can impact styling elsewhere on the webpage.

#### BEM

* Where we need to write CSS, follow the BEM naming convention to increase readability & reusability.
  * BEM is short for Block, Element, Modifier which are the three components of classnames.
  * From [CSS Tricks](https://css-tricks.com/bem-101/): "In this CSS methodology a block is a top-level abstraction of a new component, for example a button: `.btn { }`. This block should be thought of as a parent. Child items, or elements, can be placed inside and these are denoted by two underscores following the name of the block like `.btn__price { }`. Finally, modifiers can manipulate the block so that we can theme or style that particular component without inflicting changes on a completely unrelated module. This is done by appending two hyphens to the name of the block just like `btn--orange`."
  * Expanding on this, a modified child class would have a class name like `.btn__price--orange`.

#### USWDS

* Check the [Truss USWDS React package](https://github.com/trussworks/react-uswds) for a component that matches your needs. Maximize the code view to see what classes to use to replicate the component styles.
* If there isn't a component there already Check the [Truss USWDS React package](https://standards.usa.gov/components/) for a component that matches your needs. Please add it to the USWDS React code and then import the new version for use in MilMove.
* USWDS has a [Slack chat](https://chat.18f.gov/) you can go to for help. Get invited to it by filling out [this form](https://chat.18f.gov/).

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
* _Read_: [Modern JavaScript Tutorial](https://javascript.info/) - A site with tutorials covering many modern javascript concepts
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
