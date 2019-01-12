# *CSS Tooling*

Currently there is not consistency in how we are using CSS in our React code.  We are using vanilla CSS.  This causes problems with name collisions and overlap of class names.  We also have styles that are repeated that could be shared for various components and elements.  Our CSS code ends up difficult to read.

It would be helpful to have the option to use variables and calculations and to scope classes locally.

## Decision Drivers

* Allows use of variables
* Allows local scoping of classes
* Ease of implementation
* High usage and support
* Likelihood for familiarity for client inheriting the app
* Quick learning curve

## Decision Outcome

Chosen Alternative: **Sass with CSS Modules**

* **Justification:** Sass and CSS Modules are easy to learn and implement with minimal learning curve
* They are built into create-react-app 2
* Sass allows variables, mixins, and calculations
* CSS Modules allows local scoping of classes
* Sass and CSS Modules come in the box with Create React App 2
* Sass has high usage and support
* Likelihood for familiarity with Sass for client inheriting the app
* In combination, we fill two of our most pressing needs: use of variables and local scoping of classes

* **Consequences:** We need to update to create-react-app 2 first before implementing
* Engineers not familiar with Sass will need to learn it

## Considered Alternatives

* *CSS-in-JS*
* *LESS*

Resources:

* [Modular CSS with React](https://medium.com/@pioul/modular-css-with-react-61638ae9ea3e)
* [CSS Preprocessors – Sass vs LESS](https://www.keycdn.com/blog/sass-vs-less/)
* [How to use Sass and CSS Modules with create-react-app](https://blog.bitsrc.io/how-to-use-sass-and-css-modules-with-create-react-app-83fa8b805e5e)
* [All You Need to KNow about CSS-in-JS](https://hackernoon.com/all-you-need-to-know-about-css-in-js-984a72d48ebc)

## Pros and Cons of the Alternatives

### *CSS-in-JS*

* `+` Scoped locally by default: Styles can live in separate, modular files that get imported into JS modules using regular import statements
* `-` Steeper learning curve
* `-` Less likely client will be familiar and be able to easily pick up

### *CSS Modules*

* `+` all class names are scoped locally by default
* `-` variables require additional plugin

### *Sass*

* `+` Cleaner code with reusable pieces and variables
* `+` Saves you time
* `+` Easier to maintain code with snippets and libraries
* `+` Calculations and logic
* `+` More organized and easy to setup
* `+` Integrated into Create React App 2
* `+` can use variables and mixins
* `+` backwards compatible (can change .css files to .sass or .scss file)
* `+` Quick learning curve
* `+` growing in popularity (and support)
* written in Ruby

### *LESS*

* `+` Cleaner code with reusable pieces and variables
* `+` Saves you time
* `+` Easier to maintain code with snippets and libraries
* `+` Calculations and logic
* `+` More organized and easy to setup
* `+` Can use variables
* `+` Backwards compatible (can change .css files to .less file)
* `+` written in Javascript, so easy to extend
* `-` Not automatically accessible in Create React App 2