export const initialState = {
  hasRecentError: false,
  timestamp: 0,
};

const interceptorReducer = (state = initialState, action) => {
  switch (action?.type) {
    default:
      return state;
  }
};

export default interceptorReducer;
