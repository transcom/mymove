const REQUIRES_APPROVAL_MESSAGE_ID = 'no-unapproved-annotation';
const NO_ANNOTATION_MESSAGE_ID = 'no-annotation';
const messages = {
  [REQUIRES_APPROVAL_MESSAGE_ID]: 'Requires annotation approval from an ISSO',
  [NO_ANNOTATION_MESSAGE_ID]:
    'Disabling of this rule requires an annotation. Please visit https://docs.google.com/document/d/1qiBNHlctSby0RZeaPzb-afVxAdA9vlrrQgce00zjDww/edit?usp=sharing',
};

// eslint-disable-next-line security/detect-unsafe-regex
const disableRegex = /^eslint-disable(?:-next-line|-line)?(?<ruleId>$|(?:\s+(?:@(?:[\w-]+\/){1,2})?[\w-]+)?)/;

const validatorStatusOptions = new Set([
  'RA ACCEPTED',
  'RETURN TO DEVELOPER',
  'KNOWN ISSUE',
  'MITIGATED',
  'FALSE POSITIVE',
  'BAD PRACTICE',
]);

const approvedBypassableRules = new Set([
  'react/button-has-type',
  'react/destructuring-assignment',
  'react/forbid-prop-types',
  'react/jsx-boolean-value',
  'react/jsx-filename-extension',
  'react/jsx-no-bind',
  'react/jsx-props-no-spreading',
  'react/no-string-refs',
  'react/prefer-stateless-function',
  'react/prefer-stateless-function',
  'react/self-closing-comp',
  'react/sort-comp',
  'react/state-in-constructor',
  'react/static-property-placement',
  'import/extensions',
  'import/newline-after-import',
  'import/no-extraneous-dependencies',
  'import/no-mutable-exports',
  'import/no-named-as-default',
  'import/order',
  'import/prefer-default-export',
  'camelcase',
  'class-methods-use-this',
  'func-names',
  'lines-between-class-members',
  'max-classes-per-file',
  'new-cap',
  'no-alert',
  'no-extra-boolean-cast',
  'no-nested-ternary',
  'no-restricted-syntax',
  'no-return-assign',
  'no-return-await',
  'no-underscore-dangle',
  'no-unneeded-ternary',
  'object-shorthand',
  'one-var',
  'prefer-const',
  'prefer-destructuring',
  'prefer-object-spread',
  'prefer-promise-reject-errors',
  'prefer-rest-params',
  'prefer-template',
  'spaced-comment',
  'vars-on-top',
]);

const VALIDATOR_LABEL = 'RA Validator Status:';

const hasAnnotation = (commentsArr) => {
  if (!commentsArr.length) {
    return false;
  }

  return commentsArr.filter((str) => str.startsWith('RA')).length > 0;
};

const getValidatorStatus = (commentsArr) =>
  commentsArr.reduce((accum, curr) => {
    if (curr.startsWith(VALIDATOR_LABEL)) {
      // eg. RA Validator Status: Mitigated
      return curr.split(':')[1].trim();
    }
    return accum;
  }, '');

/*
  List of false positives:
  - comments directly in render func
    - ex.
    {// RA ... }
    {// eslint-disable some-rule }
  - inline eslint disables (eslint-disable-line)
    - RA Validator status: ...
    - someCode() // eslint-disable-line
*/

const create = (context) => ({
  Program: (node) => {
    node.comments.forEach((comment) => {
      const commentValue = comment.value.trim();
      const result = disableRegex.exec(commentValue);

      if (
        result && // It's a eslint-disable comment
        result.groups.ruleId // disabling a specific rule
      ) {
        const [, rule] = result.input.split(' ');
        const commentsBefore = context.getCommentsBefore(comment).map(({ value }) => value.trim());
        const validatorStatus = getValidatorStatus(commentsBefore);
        if (!approvedBypassableRules.has(rule) && !hasAnnotation(commentsBefore)) {
          context.report({
            // Can't set it at the given location as the warning
            // will be ignored due to the disable comment
            loc: {
              start: {
                ...comment.loc.start,
                column: -1,
              },
              end: comment.loc.end,
            },
            messageId: NO_ANNOTATION_MESSAGE_ID,
          });
        } else if (!approvedBypassableRules.has(rule) && !validatorStatusOptions.has(validatorStatus.toUpperCase())) {
          context.report({
            // Can't set it at the given location as the warning
            // will be ignored due to the disable comment
            loc: {
              start: {
                ...comment.loc.start,
                column: -1,
              },
              end: comment.loc.end,
            },
            messageId: REQUIRES_APPROVAL_MESSAGE_ID,
          });
        }
      }
    });
  },
});

module.exports = {
  create,
  meta: {
    type: 'suggestion',
    messages,
  },
};
