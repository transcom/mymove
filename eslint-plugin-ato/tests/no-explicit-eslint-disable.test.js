const { RuleTester } = require('eslint');

const rule = require('../no-explicit-eslint-disable');

const ruleTester = new RuleTester();

ruleTester.run('no-explicit-eslint-disable', rule, {
  valid: [
    '// eslint-disable no-console',
    '// eslint-disable react/jsx-props-no-spreading',
    '/* eslint-disable no-console */',
    '/* eslint-disable no-console, prefer-const */',
    '// eslint-disable no-console',
  ],

  invalid: [
    {
      code: '// eslint-disable',
      errors: [{ message: 'Please specify the rule(s) you want to disable.' }],
    },
    {
      code: '// eslint-disable-next-line',
      errors: [{ message: 'Please specify the rule(s) you want to disable.' }],
    },
    {
      code: '/* eslint-disable */',
      errors: [{ message: 'Please specify the rule(s) you want to disable.' }],
    },
    {
      code: '/* eslint-disable-next-line */',
      errors: [{ message: 'Please specify the rule(s) you want to disable.' }],
    },
  ],
});
