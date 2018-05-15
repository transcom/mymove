export const no_op = () => undefined;
export const no_op_action = () => {
  return function(dispatch) {
    dispatch({
      type: 'NO_OP_TYPE',
      item: null,
    });
  };
};
