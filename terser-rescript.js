//  import/no-extraneous-dependencies
const TerserPlugin = require('terser-webpack-plugin');

// modify Terser settings in webpack config
// defaults taken from https://github.com/facebook/create-react-app/blob/master/packages/react-scripts/config/webpack.config.js
module.exports = (webpackConfig) => {
  return {
    ...webpackConfig,
    optimization: {
      ...webpackConfig.optimization,
      minimizer: [
        new TerserPlugin({
          terserOptions: {
            parse: { ecma: 8 },
            compress: {
              ecma: 5,
              warnings: false,
              comparisons: false,
              inline: 2,
            },
            mangle: { safari10: true },
            output: {
              ecma: 5,
              comments: false,
              ascii_only: true,
            },
          },
          parallel: false,
          sourceMap: false,
          cache: true,
        }),
      ],
    },
  };
};
