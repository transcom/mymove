module.exports = {
  stories: ['../src/**/*.stories.jsx'],
  addons: [
    '@storybook/preset-create-react-app',
    '@storybook/addon-a11y',
    '@storybook/addon-actions',
    '@storybook/addon-essentials',
    '@storybook/addon-knobs',
    '@storybook/addon-links',
    '@storybook/addon-viewport',
    'storybook-addons-abstract',
    {
      name: '@storybook/addon-docs',
      options: {
        configureJSX: true,
      },
    },
  ],
};
