// Select the moveId that is set from clicking on Go To Move on the MultiMoveLandingPage
export function selectCurrentMoveId(state) {
  return state.generalState.moveId;
}

export function selectRefetchQueue(state) {
  return state.generalState.refetchQueue;
}

export default {
  selectCurrentMoveId,
  selectRefetchQueue,
};
