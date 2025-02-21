// Reducer created to store needed information in state
import { SET_CAN_ADD_ORDERS, SET_MOVE_ID, SET_REFETCH_QUEUE } from './actions';

export const initialState = {
  // Select the moveId that is set from clicking on Go To Move on the MultiMoveLandingPage
  moveId: '',
  canAddOrders: false,
  refetchQueue: false,
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
    case SET_CAN_ADD_ORDERS: {
      return {
        ...state,
        canAddOrders: action.payload,
      };
    }
    case SET_REFETCH_QUEUE: {
      return {
        ...state,
        refetchQueue: action.payload,
      };
    }
    default:
      return state;
  }
};

export default generalStateReducer;
