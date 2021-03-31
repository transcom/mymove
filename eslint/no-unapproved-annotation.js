const MESSAGE_ID = 'no-unapproved-annotation';
const NO_ANNOTATION_MESSAGE_ID = 'no-annotation';
const messages = {
  [MESSAGE_ID]: 'Requires annotation approval from an ISSO',
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

// const rulesRequiringAnnotation = new Set(['no-console', 'security', 'react']);

/*
Lint reqs:
- (x) for disabling of a specific rule, we're checking to see if it has annotations at all (show err if not and link to documentation)
  - doc link: https://docs.google.com/document/d/1qiBNHlctSby0RZeaPzb-afVxAdA9vlrrQgce00zjDww/edit?usp=sharing
- if it has an annotation, check the validator status and ensure that it is not empty
- if it has an annotation and validator status is not empty, check if that status is a single value
*/
const VALIDATOR_LABEL = 'RA Validator Status:';

const hasAnnotation = (context, comment) => {
  const possibleAnnotation = context.getCommentsBefore(comment);
  if (!possibleAnnotation.length) {
    return false;
  }
  const containsAnnotationBlock =
    possibleAnnotation.map(({ value }) => value.trim()).filter((str) => str.startsWith('RA')).length > 0;

  return containsAnnotationBlock;
};
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
        if (!approvedBypassableRules.has(rule) && !hasAnnotation(context, comment)) {
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
        }
        // const validatorComment = context
        //   .getCommentsBefore(comment)
        //   .map(({ value }) => value.trim())
        //   .filter((value) => value.includes(VALIDATOR_LABEL))[0];
        // if (validatorComment) {
        // const [, validatorStatusValue] = validatorComment.split(':');
        //   if (!validatorStatusOptions.has(validatorStatusValue.trim())) {
        //     context.report({
        //       // Can't set it at the given location as the warning
        //       // will be ignored due to the disable comment
        //       loc: {
        //         start: {
        //           ...comment.loc.start,
        //           column: -1,
        //         },
        //         end: comment.loc.end,
        //       },
        //       messageId: MESSAGE_ID,
        //     });
        //   }
        // }
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
