const editBeginType = 'EDIT_BEGIN';

export function editBegin() {
  return function(dispatch, getState) {
    dispatch({ type: editBeginType });
  };
}

const editSuccessfulType = 'EDIT_SUCCESS';

export function editSuccessful() {
  return function(dispatch, getState) {
    dispatch({ type: editSuccessfulType });
  };
}

export function reviewReducer(state = {}, action) {
  switch (action.type) {
    case editBeginType:
      return Object.assign({}, state, {
        editSuccess: false,
      });
    case editSuccessfulType:
      return Object.assign({}, state, {
        editSuccess: true,
      });
    default:
      return state;
  }
}
