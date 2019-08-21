# How to Use and Run Storybook

## What is Storybook

Storybook is a user interface development environment and playground for UI components. The tool enables developers to create components independently and showcase components interactively in an isolated development environment. [Read more here](https://storybook.js.org/docs/basics/introduction/)

## Basics

Storybook expects _stories_ to be defined for each component to be showcased. These stories are defined in the stories directory `src/stories`

### How to run storybook server locally

To see the components locally simply run `make storybook` and the server will start and automatically open a browser window. If not open [http://localhost:6006](http://localhost:6006)

### How to generate static storybook site files

If you wish to generate the static version of storybook run `make build_storybook` and the command will generate the files in `storybook-static`
