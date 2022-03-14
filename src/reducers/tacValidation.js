import { LOA_TYPE } from 'shared/constants';

export const TAC_VALIDATION_ACTIONS = {
  VALIDATION_RESPONSE: 'VALIDATION_RESPONSE',
};

export const reducer = (state, action) => {
  switch (action.type) {
    case TAC_VALIDATION_ACTIONS.VALIDATION_RESPONSE: {
      return {
        ...state,
        [action.loaType]: {
          isValid: action.isValid,
          tac: action.tac,
        },
      };
    }

    default:
      return state;
  }
};

export const initialState = () => {
  return {
    [LOA_TYPE.HHG]: {
      isValid: true,
      tac: '',
    },
    [LOA_TYPE.NTS]: {
      isValid: true,
      tac: '',
    },
  };
};
