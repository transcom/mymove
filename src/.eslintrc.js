module.exports = {
  rules: {
    // okay to disable because the threat actor (web application user) already controls the execution environment (web browser)
    'security/detect-object-injection': 'off',
  },
};
