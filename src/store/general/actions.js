export const SET_MOVE_ID = 'SET_MOVE_ID';

// Action for setting moveId when clicking on Go To Move on the MultiMoveLandingPage
export const setMoveId = (moveId) => ({
  type: SET_MOVE_ID,
  payload: moveId,
});
