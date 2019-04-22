// SINGLE-RESOURCE ACTION CREATORS

// Reducer
const initialState = {
  storageInTransitIsCreating: false,
  storageInTransitHasCreatedSuccess: false,
  storageInTransitHasCreatedError: null,
  storageInTransitsAreLoading: false,
  storageInTransitsHasLoadSuccess: false,
  storageInTransitsHasLoadError: null,
  flashMessage: false,
};

export function tspReducer(state = initialState, action) {
  switch (action.type) {
    default:
      return state;
  }
}
