const path = require('path');
module.exports = {
  stories: ['../src/**/*.stories.jsx'],
  addons: [
    '@storybook/preset-create-react-app',
    '@storybook/addon-a11y',
    {
      name: '@storybook/addon-docs',
      options: {
        configureJSX: true,
      },
    },
    '@storybook/addon-essentials',
    '@storybook/addon-knobs',
    '@storybook/addon-links',
    '@dump247/storybook-state',
  ],
  webpackFinal: async (config) => {
    config.resolve.modules = config.resolve.modules || [];
    config.resolve.modules.push(path.resolve(__dirname, '../src'));
    config.resolve.modules.push('node_modules');
    return config;
  },
  refs: {
    'design-system': {
      title: 'ReactUSWDS',
      url: 'https://trussworks.github.io/react-uswds/',
    },
  },
  framework: {
    name: '@storybook/react-webpack5',
    options: {},
  },
  docs: {
    autodocs: true,
  },
};
