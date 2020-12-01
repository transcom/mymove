// TODO: remove this after refactor is implemented for creating a flash message

export const SHOW_SUCCESS_BANNER = 'SHOW_SUCCESS_BANNER';
export const REMOVE_SUCCESS_BANNER = 'REMOVE_SUCCESS_BANNER';

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
