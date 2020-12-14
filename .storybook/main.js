module.exports = {
  stories: ['../src/**/*.stories.jsx'],
  addons: [
    // '@storybook/addon-a11y',
    '@storybook/preset-create-react-app',
    '@storybook/addon-essentials',
    {
      name: '@storybook/addon-docs',
      options: {
        configureJSX: true,
      },
    },
    '@storybook/addon-knobs',
    '@storybook/addon-links',
    // 'storybook-addons-abstract',
    '@dump247/storybook-state',
  ],
};
