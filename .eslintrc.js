module.exports = {
  plugins: ['prettier', 'security', 'no-only-tests', 'you-dont-need-lodash-underscore', 'ato'],
  extends: [
    'react-app',
    'airbnb',
    'prettier',
    'prettier/prettier',
    'plugin:security/recommended',
    'plugin:you-dont-need-lodash-underscore/compatible',
  ],
  root: true,
  rules: {
    'prettier/prettier': 'error',
    'no-only-tests/no-only-tests': 'error',
    'ato/no-explicit-eslint-disable': 'error',
    'ato/no-unapproved-annotation': 'off',
    'no-debugger': 'warn',
    'no-console': 'warn',
    'jsx-a11y/anchor-is-valid': ['warn', { aspects: ['noHref', 'preferButton'] }],
    'jsx-a11y/label-has-associated-control': [
      'error',
      {
        labelAttributes: ['htmlFor'],
        assert: 'either',
        depth: 2,
      },
    ],
    'react/button-has-type': 'error',
    'react/jsx-no-target-blank': 'error',
    'react/jsx-curly-brace-presence': 'error',
    'react/require-default-props': 'warn',
    'no-unused-vars': 'warn',
    'import/order': [
      'error',
      {
        'newlines-between': 'always',
      },
    ],
    'you-dont-need-lodash-underscore/capitalize': 'off',
    'you-dont-need-lodash-underscore/clone': 'off',
    'you-dont-need-lodash-underscore/cloneDeep': 'off',
    'you-dont-need-lodash-underscore/findKey': 'off',
    'you-dont-need-lodash-underscore/findLast': 'off',
    'you-dont-need-lodash-underscore/mapValues': 'off',
    'you-dont-need-lodash-underscore/memoize': 'off',
    'you-dont-need-lodash-underscore/snakeCase': 'off',
    'you-dont-need-lodash-underscore/startCase': 'off',
    'you-dont-need-lodash-underscore/sum': 'off',
    'you-dont-need-lodash-underscore/union': 'off',
    'you-dont-need-lodash-underscore/uniqueId': 'off',
    'you-dont-need-lodash-underscore/get': 'off',
  },
  overrides: [
    {
      files: ['*.stories.js', '*.stories.jsx', 'setupTests.js'],
      rules: {
        'import/no-extraneous-dependencies': 'off',
        'react/jsx-props-no-spreading': 'off',
        'react/destructuring-assignment': 'off',
      },
    },
    {
      files: ['*.test.jsx'],
      rules: {
        'react/jsx-props-no-spreading': 'off',
      },
    },
  ],
  settings: {
    'import/resolver': {
      node: {
        paths: ['src'],
        extensions: ['.js', '.jsx', '.ts', '.tsx'],
      },
    },
    linkComponents: [
      // Components used as alternatives to <a> for linking, eg. <Link to={ url } />
      {
        name: 'Link',
        linkAttribute: 'href',
      },
      {
        name: 'Link',
        linkAttribute: 'to',
      },
    ],
  },
};
