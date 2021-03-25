const noRulelessEslintDisable = require('./no-ruleless-eslint-disable');
const noUnapprovedAnnotation = require('./no-unapproved-annotation');

module.exports = {
  rules: {
    'no-ruleless-eslint-disable': noRulelessEslintDisable,
    'no-unapproved-annotation': noUnapprovedAnnotation,
  },
};
