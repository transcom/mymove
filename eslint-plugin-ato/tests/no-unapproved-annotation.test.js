const { RuleTester } = require('eslint');

const rule = require('../no-unapproved-annotation');

const ruleTester = new RuleTester();

const ERRORS = {
  REQUIRES_APPROVAL_MSG:
    'Due to an added annotation, this PR requires approval from a codeowner. Once a codeowner has reviewed/approved your PR, you will need to change the RA Validator Status to CODEOWNER ACCEPTED. For more information, please visit https://dp3.atlassian.net/wiki/spaces/MT/pages/1920991340/Guide+to+Static+Analysis+Security+Workflow',
  REQUIRES_ANNOTATION_MSG:
    'Disabling of this rule requires an annotation. Please visit https://dp3.atlassian.net/wiki/spaces/MT/pages/1921122376/Guide+to+Static+Analysis+Annotations+for+Disabled+Linters',
  NO_INLINE_DISABLE_MSG: 'Please use eslint-disable-next-line instead of eslint-disable-line',
};
ruleTester.run('no-unapproved-annotation', rule, {
  valid: [
    '// RA Validator Status: Mitigated\n// eslint-disable no-console',
    '// RA Validator Status: Known Issue\n// eslint-disable security/detect-unsafe-regex',
    '// RA Validator Status: RA Accepted\n// eslint-disable-next-line no-console',
    '// RA Validator Status: Mitigated\n// eslint-disable-next-line no-console, prefer-const ',
    // approved to be bypassed, does not require annotation
    '// eslint-disable react/button-has-type',
    '// eslint-disable react/destructuring-assignment',
    '// eslint-disable react/forbid-prop-types',
    '// eslint-disable react/jsx-boolean-value',
    '// eslint-disable react/jsx-filename-extension',
    '// eslint-disable react/jsx-no-bind',
    '// eslint-disable react/jsx-props-no-spreading',
    '// eslint-disable react/no-string-refs',
    '// eslint-disable react/prefer-stateless-function',
    '// eslint-disable react/prefer-stateless-function',
    '// eslint-disable react/self-closing-comp',
    '// eslint-disable react/sort-comp',
    '// eslint-disable react/state-in-constructor',
    '// eslint-disable react/static-property-placement',
    '// eslint-disable import/extensions',
    '// eslint-disable import/newline-after-import',
    '// eslint-disable import/no-extraneous-dependencies',
    '// eslint-disable import/no-mutable-exports',
    '// eslint-disable import/no-named-as-default',
    '// eslint-disable import/order',
    '// eslint-disable import/prefer-default-export',
    '// eslint-disable camelcase',
    '// eslint-disable class-methods-use-this',
    '// eslint-disable func-names',
    '// eslint-disable lines-between-class-members',
    '// eslint-disable max-classes-per-file',
    '// eslint-disable new-cap',
    '// eslint-disable no-alert',
    '// eslint-disable no-extra-boolean-cast',
    '// eslint-disable no-nested-ternary',
    '// eslint-disable no-restricted-syntax',
    '// eslint-disable no-return-assign',
    '// eslint-disable no-return-await',
    '// eslint-disable no-underscore-dangle',
    '// eslint-disable no-unneeded-ternary',
    '// eslint-disable object-shorthand',
    '// eslint-disable one-var',
    '// eslint-disable prefer-const',
    '// eslint-disable prefer-destructuring',
    '// eslint-disable prefer-object-spread',
    '// eslint-disable prefer-promise-reject-errors',
    '// eslint-disable prefer-rest-params',
    '// eslint-disable prefer-template',
    '// eslint-disable spaced-comment',
    '// eslint-disable vars-on-top',
  ],

  invalid: [
    {
      code: '// RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}\n// eslint-disable no-console',
      errors: [{ message: ERRORS.REQUIRES_APPROVAL_MSG }],
    },
    {
      code: '// RA Validator Status: \n// eslint-disable-next-line no-console',
      errors: [{ message: ERRORS.REQUIRES_APPROVAL_MSG }],
    },
    {
      code: '// eslint-disable security/detect-unsafe-regex',
      errors: [
        {
          message: ERRORS.REQUIRES_ANNOTATION_MSG,
        },
      ],
    },
    {
      code: '// eslint-disable-line no-unused-vars',
      errors: [
        {
          message: ERRORS.NO_INLINE_DISABLE_MSG,
        },
        {
          message: ERRORS.REQUIRES_ANNOTATION_MSG,
        },
      ],
    },
  ],
});
