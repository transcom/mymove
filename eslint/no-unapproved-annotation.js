const MESSAGE_ID = 'no-unapproved-annotation';

const messages = {
  [MESSAGE_ID]: 'Requires annotation approval from an ISSO',
};

const disableRegex = /^eslint-disable(?:-next-line|-line)?(?<ruleId>$|(?:\s+(?:@(?:[\w-]+\/){1,2})?[\w-]+)?)/;

const validatorStatusOptions = new Set([
  'RA Accepted',
  'Return to Developer',
  'Known Issue',
  'Mitigated',
  'False Positive',
  'Bad Practice',
]);

const rulesRequiringAnnotation = new Set(['no-console', 'security', 'react']);

const VALIDATOR_LABEL = 'RA Validator Status:';
const create = (context) => ({
  Program: (node) => {
    node.comments.forEach((comment) => {
      const commentValue = comment.value.trim();
      const result = disableRegex.exec(commentValue);

      if (
        result && // It's a eslint-disable comment
        rulesRequiringAnnotation.has(result.groups.ruleId) && // Specifies a rule that we need annotations for
        context.getCommentsBefore(comment).length
      ) {
        console.log('what', result.groups);
        const validatorComment = context
          .getCommentsBefore(comment)
          .map(({ value }) => value.trim())
          .filter((value) => value.includes(VALIDATOR_LABEL))[0];
        const [, validatorStatusValue] = validatorComment.split(':');

        if (!validatorStatusOptions.has(validatorStatusValue.trim())) {
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
            messageId: MESSAGE_ID,
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
