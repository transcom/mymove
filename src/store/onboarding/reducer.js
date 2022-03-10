import { SET_CONUS_STATUS, SET_PPM_ESTIMATE_ERROR } from './actions';

export const initialState = {
  conusStatus: null,
  ppmEstimateError: null,
};

const onboardingReducer = (state = initialState, action = {}) => {
  switch (action?.type) {
    case SET_CONUS_STATUS: {
      const { moveType } = action;

      return {
        ...state,
        conusStatus: moveType,
      };
    }

    case SET_PPM_ESTIMATE_ERROR: {
      const { error } = action;

      return {
        ...state,
        ppmEstimateError: error,
      };
    }

    default:
      return state;
  }
};

export default onboardingReducer;
