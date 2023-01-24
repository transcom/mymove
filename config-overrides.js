// README: This is necessary due to the fact that webpack 5 dropped support for
// Node polyfills. We use the `process` object to do things in a couple of
// files. Search the project for `process.` to see where. To mitigate this
// update to webpack 5 and not update our code, we can configure webpack to
// support `process` again by using the `stream-http` package. We also override
// our react-scripts calls with a new package called `react-app-rewired` which
// wraps calls to `react-scripts` and allows for configuration updates to
// configuration files inside of Create React App.

// webpack is a dependency of React-Scripts
const webpack = require('webpack');

module.exports = {
  webpack: (config) => {
    config.resolve.fallback = {
      // This is the Node polyfill for process/browser
      http: require.resolve('stream-http'),
    };
    config.plugins.push(
      new webpack.ProvidePlugin({
        process: 'process/browser',
      }),
    );

    // Loop through all the plugins in the Array and update the
    // MiniCssExtractPlugin one with the ignoreOrder flag as True
    // README: This has to be done because of the fact that Create-React-App
    // already includes the mini-css-extract-plugin and configures it. This
    // means that just pushing a new configuration won't modify the original
    // one. So we need to find the original and override just the properties we
    // need. There may be better ways to do this but this works for now.
    config.plugins.forEach((p) => {
      if (p.options) {
        // Ignore warnings about ordering from mini-css-extract-plugin
        if (p.options.hasOwnProperty('ignoreOrder')) {
          p.options.ignoreOrder = true;
        }
      }
    });

    return config;
  },
  jest: (config) => {
    config.collectCoverageFrom = [
      '**/src/**/*.{js,jsx,ts,tsx}',
      '!**/src/**/*.stories.{js,jsx,ts,tsx}',
      '!**/node_modules/**',
    ];

    config.coverageThreshold = {
      global: {
        branches: 60,
        functions: 40,
        lines: 60,
        statements: 60,
      },
    };

    return config;
  },
};
