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
  'RA Accepted',
  'Return to Developer',
  'Known Issue',
  'Mitigated',
  'False Positive',
  'Bad Practice',
]);

const approvedBypassableRules = new Set([
  'no-underscore-dangle',
  'prefer-object-spread',
  'object-shorthand',
  'camelcase',
  'react/jsx-props-no-spreading',
  'react/destructuring-assignment',
  'react/forbid-prop-types',
  'react/prefer-stateless-function',
  'react/sort-comp',
  'import/no-extraneous-dependencies',
  'import/order',
  'import/prefer-default-export',
  'import/no-named-as-default',
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
        } else if (!approvedBypassableRules.has(rule) && !validatorStatusOptions.has(validatorStatus)) {
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
