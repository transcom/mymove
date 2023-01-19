module.exports = {
  plugins: ['prettier', 'security', 'no-only-tests', 'you-dont-need-lodash-underscore', 'ato', 'import'],
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
    'ato/no-unapproved-annotation': process.env.NODE_ENV === 'production' ? 'error' : 'off',
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
    'react/function-component-definition': [
      'error',
      {
        namedComponents: ['function-declaration', 'function-expression', 'arrow-function'],
        unnamedComponents: ['function-expression', 'arrow-function'],
      },
    ],
    'react/jsx-no-target-blank': 'error',
    'react/jsx-curly-brace-presence': 'error',
    'react/require-default-props': 'warn',
    // TODO: The following rules were introduced in the dependency update to
    // `eslint-plugin-react` and require some additiona refactor of things in
    // our app to enable. They're disabled for now, but we should address these
    // in future changes to the affected files. {
    'react/prop-types': 'off',
    'react/no-array-index-key': 'off',
    'react/forbid-prop-types': 'off',
    // }
    'no-unused-vars': 'warn',
    'import/order': [
      'error',
      {
        'newlines-between': 'always',
      },
    ],
    'import/named': 'error',
    // TODO: The folowing rules were introduced from exisiting eslint preset libraries we use.
    // Disabling them for now until we have a conversation to see which rules we would like to keep from
    // the import library.
    'import/no-named-as-default-member': 'off',
    'import/no-named-as-default': 'off',
    //
    'import/no-extraneous-dependencies': [
      'error',
      {
        devDependencies: [
          '**/*.stories.js*',
          // '**/*.stories.js',
          '**/*.test.js*',
          // '**/*.test.jsx',
          '**/setupTests.js',
          '**/testUtils.jsx',
          '**/test/factories/**',
          // playwright
          '**/playwright.config.js',
          '**/playwright/**/*.js',
        ],
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
        'react/jsx-props-no-spreading': 'off',
        'react/destructuring-assignment': 'off',
      },
    },
    {
      files: ['*.test.jsx', 'testUtils.jsx'],
      rules: {
        'react/jsx-props-no-spreading': 'off',
      },
    },
    {
      files: ['src/utils/test/factories/**'],
      rules: {
        'no-param-reassign': ['error', { props: false }],
      },
    },
    {
      files: ['playwright/**/*.js*'],
      rules: {
        'no-restricted-syntax': 'off',
        'no-await-in-loop': 'off',
      },
    },
  ],
  settings: {
    'import/resolver': {
      node: {
        moduleDirectory: ['src', 'node_modules'],
        extensions: ['.js', '.jsx'],
      },
    },
    'import/ignore': ['.coffee$'],
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
