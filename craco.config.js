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

          return webpackConfig;
        },
      },
    },
  ],
};
