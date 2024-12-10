export function selectGetCurrentUserIsSuccess(state) {
  return state.auth.hasSucceeded;
}

export const selectIsLoggedIn = (state) => {
  return state.auth.isLoggedIn;
};

export function selectGetCurrentUserIsLoading(state) {
  return state.auth.isLoading;
}

export function selectGetCurrentUserIsError(state) {
  return state.auth.hasErrored;
}

export const selectCacValidated = (serviceMember) => {
  return serviceMember?.cac_validated || false;
};

export const selectUnderMaintenance = (state) => {
  return state.auth.underMaintenance;
};
