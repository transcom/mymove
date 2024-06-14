export const LOA_VALIDATION_ACTIONS = {
  VALIDATION_RESPONSE: 'VALIDATION_RESPONSE',
};

export const reducer = (state, action) => {
  switch (action.type) {
    case LOA_VALIDATION_ACTIONS.VALIDATION_RESPONSE: {
      return {
        ...state,
        isValid: action.payload.isValid,
        loa: action.payload.loa,
      };
    }
    default:
      return state;
  }
};

export const initialState = () => {
  return {
    isValid: false,
    loa: null,
  };
};
