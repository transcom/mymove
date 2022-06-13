const webpack = require('webpack');

module.exports = function override(config /* , env */) {
  config.resolve.fallback = {
    http: require.resolve('stream-http'),
  };
  config.plugins.push(
    new webpack.ProvidePlugin({
      process: 'process/browser',
    }),
  );

  return config;
};
