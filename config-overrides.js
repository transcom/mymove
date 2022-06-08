// README: This is necessary due to the fact that Webpack 5 dropped support for
// Node Polyfills. We use the `process` object to do things in a couple of
// files. Search the project for `process.` to see where. To mitigate this
// update to Webpack 5 and not update our code, we can configure Webpack to
// support `process` again by using the `stream-http` package. We also override
// our react-scripts calls with a new package called `react-app-rewired` which
// wraps calls to `react-scripts` and allows for configuration updates to
// configuration files inside of Create React App.

// Webpack is a dependency of React-Scripts
const webpack = require('webpack');

module.exports = function override(config /* , env */) {
  config.resolve.fallback = {
    // This is the Node Polyfill for process/bowser
    http: require.resolve('stream-http'),
  };
  config.plugins.push(
    new webpack.ProvidePlugin({
      process: 'process/browser',
    }),
  );

  return config;
};
