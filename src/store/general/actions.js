export const SET_MOVE_ID = 'SET_MOVE_ID';

// Action for setting moveId when clicking on Go To Move on the MultiMoveLandingPage
export const setMoveId = (moveId) => ({
  type: SET_MOVE_ID,
  payload: moveId,
});

export const SET_CAN_ADD_ORDERS = 'SET_CAN_ADD_ORDERS';

export const setCanAddOrders = (value) => ({
  type: SET_CAN_ADD_ORDERS,
  payload: value,
});

export const SET_REFETCH_QUEUE = 'SET_REFETCH_QUEUE';

export const setRefetchQueue = (value) => ({
  type: SET_REFETCH_QUEUE,
  payload: value,
});
export const SET_SHOW_LOADING_SPINNER = 'SET_SHOW_LOADING_SPINNER';

export const setShowLoadingSpinner = (showSpinner, loadingSpinnerMessage) => {
  return {
    type: SET_SHOW_LOADING_SPINNER,
    showSpinner,
    loadingSpinnerMessage,
  };
};
