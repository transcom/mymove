const MESSAGE_ID = 'no-ruleless-eslint-disable';
const messages = {
  [MESSAGE_ID]: 'Please specify the rule(s) you want to disable.',
};

// eslint-disable-next-line security/detect-unsafe-regex
const disableRegex = /^eslint-disable(?:-next-line|-line)?(?<ruleId>$|(?:\s+(?:@(?:[\w-]+\/){1,2})?[\w-]+)?)/;

const create = (context) => ({
  Program: (node) => {
    node.comments.forEach((comment) => {
      const value = comment.value.trim();
      const result = disableRegex.exec(value);
      if (
        result && // It's a eslint-disable comment
        !result.groups.ruleId // But it did not specify any rules
      ) {
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
