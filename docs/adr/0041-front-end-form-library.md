# Front End Form Library

The current form library we are using with React, [Redux Form](https://redux-form.com/8.2.2/), is not serving us well. Here is why:

* Difficult to unit test
* Lots of overhead if we want to add forms to Storybook
* Having issues dealing with some of the more nuanced features, such as asynchronously checking form values

Therefore, we want to replace it with a new library that fixes these problems and, hopefully, future-proofs us a bit.

## Considered Options

* [Formik](https://github.com/jaredpalmer/formik)
* [React Hook Form](https://react-hook-form.com/)
* Implementing our own pattern

## Decision Outcome

* Chosen Option: Formik
* Forms are implemented using plain JSX (no need for higher-order components)
* Easy to unit test
* Forms can be added to Storybook without needing a fake redux store
* Uses current React patterns that should be easy for devs on the project to pick up
* Meets all of our needs for our form validation pattern

## Pros and Cons of the Alternatives <!-- optional -->

### React Hook Form

* `+` Uses newer React hook style of code, which future-proofs us a bit
* `+` No library dependencies
* `+` Highly performant
* `-` Uses newer React hook style of code, which might be harder for project devs to pick up
* `-` Does not support our form validation pattern
* `-` Documentation is not as well written as we would like

### Implementing our own pattern

* `+` Available to customize to meet whatever product needs we have
* `-` We need to build and support the code
* `-` Off-the-shelf solutions already do what we need to, with no need to spend dev time on it