import { SET_CONUS_STATUS } from './actions';

export const initialState = {
  conusStatus: null,
};

const onboardingReducer = (state = initialState, action) => {
  switch (action?.type) {
    case SET_CONUS_STATUS: {
      const { moveType } = action;

      return {
        ...state,
        conusStatus: moveType,
      };
    }

    default:
      return state;
  }
};

export default onboardingReducer;
