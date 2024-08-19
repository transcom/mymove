import { LOA_TYPE } from 'shared/constants';

export const LOA_VALIDATION_ACTIONS = {
  VALIDATION_RESPONSE: 'VALIDATION_RESPONSE',
};

export const reducer = (state, action) => {
  switch (action.type) {
    case LOA_VALIDATION_ACTIONS.VALIDATION_RESPONSE: {
      return {
        ...state,
        [action.payload.loaType]: {
          isValid: action.payload.isValid,
          longLineOfAccounting: action.payload.longLineOfAccounting,
          loa: action.payload.loa,
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
      isValid: false,
      longLineOfAccounting: '',
      loa: null,
    },
    [LOA_TYPE.NTS]: {
      isValid: false,
      longLineOfAccounting: '',
      loa: null,
    },
  };
};
