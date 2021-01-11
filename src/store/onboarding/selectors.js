export const selectConusStatus = (state) => {
  return state.onboarding.conusStatus;
};

export function selectPPMEstimateError(state) {
  return state.onboarding.ppmEstimateError || null;
}
