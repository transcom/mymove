module.exports = {
  rules: {
    'no-template-literals': {
      create(context) {
        return {
          TemplateLiteral(node) {
            context.report(node, 'Do not use template literals ayayayayaya');
          },
        };
      },
    },
  },
};
