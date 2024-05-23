// Reducer created to store needed information in state
import { SET_MOVE_ID } from './actions';

export const initialState = {
  // Select the moveId that is set from clicking on Go To Move on the MultiMoveLandingPage
  moveId: '',
};

const generalStateReducer = (state = initialState, action = {}) => {
  switch (action?.type) {
    // Action is fired clicking on Go To Move on the MultiMoveLandingPage
    case SET_MOVE_ID: {
      return {
        ...state,
        moveId: action.payload,
      };
    }

    default:
      return state;
  }
};

export default generalStateReducer;
