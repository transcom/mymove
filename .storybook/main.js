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
    // 'storybook-addons-abstract',
    '@dump247/storybook-state',
  ],
  refs: {
    'design-system': {
      title: 'ReactUSWDS',
      url: 'https://trussworks.github.io/react-uswds/',
    },
  },
};
