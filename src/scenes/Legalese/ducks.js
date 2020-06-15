const SHOW_SUCCESS_BANNER = 'SHOW_SUCCESS_BANNER';
const REMOVE_SUCCESS_BANNER = 'REMOVE_SUCCESS_BANNER';

export const showSubmitSuccessBanner = () => {
  return {
    type: SHOW_SUCCESS_BANNER,
  };
};

export const removeSubmitSuccessBanner = () => {
  return {
    type: REMOVE_SUCCESS_BANNER,
  };
};

// Reducer
const initialState = {
  moveSubmitSuccess: false,
  error: null,
};
export function signedCertificationReducer(state = initialState, action) {
  switch (action.type) {
    case SHOW_SUCCESS_BANNER:
      return Object.assign({}, state, {
        moveSubmitSuccess: true,
      });
    case REMOVE_SUCCESS_BANNER:
      return Object.assign({}, state, {
        moveSubmitSuccess: false,
      });
    default:
      return state;
  }
}
