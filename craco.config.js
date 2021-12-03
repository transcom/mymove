// eslint-disable-next-line import/no-extraneous-dependencies
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin');
// eslint-disable-next-line import/no-extraneous-dependencies
// const { addBeforeLoaders, loaderByName } = require('@craco/craco');

module.exports = {
  plugins: [
    {
      plugin: {
        overrideWebpackConfig: ({ webpackConfig }) => {
          const isCircleCi = process.env.CIRCLECI || false;
          // CircleCI doesn't report the right number of CPUs
          // https://ideas.circleci.com/cloud-feature-requests/p/provide-environment-variables-for-resource-limits

          // we have to fake it
          const parallel = isCircleCi ? 4 : true;

          const minimizerIndex = webpackConfig.optimization.minimizer.findIndex((item) => item.options.terserOptions);

          // set parallel for circleci
          // eslint-disable-next-line no-param-reassign, security/detect-object-injection
          webpackConfig.optimization.minimizer[minimizerIndex].options.terserOptions.parallel = parallel;

          // const threadLoader = {
          //   loader: require.resolve('thread-loader'),
          // };
          // addBeforeLoaders(webpackConfig, loaderByName('sass-loader'), threadLoader);
          // addBeforeLoaders(webpackConfig, loaderByName('babel-loader'), threadLoader);

          if (process.env.WEBPACK_SPEED_MEASURE) {
            const smp = new SpeedMeasurePlugin();
            return smp.wrap(webpackConfig);
          }

          return webpackConfig;
        },
      },
    },
  ],
};
