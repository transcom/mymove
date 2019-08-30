# How to Use and Run Storybook

## What is Storybook

Storybook is a user interface development environment and playground for UI components. The tool enables developers to create components independently and showcase components interactively in an isolated development environment. [Read more here](https://storybook.js.org/docs/basics/introduction/)

## Basics

Storybook expects _stories_ to be defined for each component to be showcased. These stories are defined in the stories directory `src/stories`

### Dependencies

If this is your first time running storybook you should run `make client_deps` first to ensure storybook packages are installed

### How to run storybook server locally

To see the components locally simply run `make storybook` and the server will start and automatically open a browser window. If not open [http://localhost:6006](http://localhost:6006)

### How to generate static storybook site files

If you wish to generate the static version of storybook run `make build_storybook` and the command will generate the files in `storybook-static`

## Adding Stories

To showcase a component add the _stories_ to the `src/stories` folder in an appropriate file. The storybook documentation on [Writing Stories](https://storybook.js.org/docs/basics/writing-stories/) is a good place to start with how to create ones. If there is not an appropriate file you need to create a new file in the pattern `componentName.stories.js` in the src/stories directory, and then modify the `.storybook/config.js` file to include your new file in the generated site.

### Addons

Stories in Storybook can use addons to extend the features of Storybook. Some addons already included are the [actions](https://github.com/storybookjs/storybook/tree/master/addons/actions) and [knobs](https://github.com/storybookjs/storybook/tree/master/addons/knobs) addons. The controls for each of these addons shows up in a pane at the bottom of the page where they are used. If you cannot find the pane try toggling the _Change addons orientation_ from the ellipsis menu next to the logo in the upper left, or using the **D** keyboard short cut to toggle it.

#### Actions

Storybook Addon Actions can be used to display data received by event handlers in Storybook. See [the documentation](https://github.com/storybookjs/storybook/tree/master/addons/actions) for more details.

#### Knobs

Storybook Addon Knobs allow you to edit props dynamically using the Storybook UI. You can also use Knobs as a dynamic variable inside stories in Storybook. See [the documentation](https://github.com/storybookjs/storybook/tree/master/addons/knobs) for more details.
