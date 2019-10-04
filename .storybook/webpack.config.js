// remove eslint from the webpack configuration to use the create-react-app eslint config
const removeEsLint = exports => {
  return (
    exports.module.rules
      .filter(e => e.use && e.use.some(e => e.options && void 0 !== e.options.useEslintrc))
      .forEach(s => {
        exports.module.rules = exports.module.rules.filter(e => e !== s);
      }),
    exports
  );
};

module.exports = async ({ config, mode }) => {
  // remove eslint from config
  config = removeEsLint(config);

  // return the new config
  return config;
};
