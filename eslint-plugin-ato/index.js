const noExplicitEslintDisable = require('./no-explicit-eslint-disable');
const noUnapprovedAnnotation = require('./no-unapproved-annotation');

module.exports = {
  rules: {
    'no-explicit-eslint-disable': noExplicitEslintDisable,
    'no-unapproved-annotation': noUnapprovedAnnotation,
  },
};
