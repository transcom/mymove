import { IssuesIndex } from 'shared/api';

// Types
export const SHOW_ISSUES = 'SHOW_ISSUES';
export const SHOW_ISSUES_SUCCESS = 'SHOW_ISSUES_SUCCESS';
export const SHOW_ISSUES_FAILURE = 'SHOW_ISSUES_FAILURE';

// Actions
export const createShowIssuesRequest = () => ({
  type: SHOW_ISSUES,
});

export const createShowIssuesSuccess = items => ({
  type: SHOW_ISSUES_SUCCESS,
  items,
});

export const createShowIssuesFailure = error => ({
  type: SHOW_ISSUES_FAILURE,
  error,
});

// Action Creator
export function loadIssues() {
  // Interpreted by the thunk middleware:
  return function(dispatch, getState) {
    dispatch(createShowIssuesRequest());
    IssuesIndex()
      .then(items => dispatch(createShowIssuesSuccess(items)))
      .catch(error => dispatch(createShowIssuesFailure(error)));
  };
}

// Reducer
function issuesReducer(state = { issues: null, hasError: false }, action) {
  switch (action.type) {
    case SHOW_ISSUES_SUCCESS:
      return { issues: action.items, hasError: false };
    case SHOW_ISSUES_FAILURE:
      return { issues: null, hasError: true };
    default:
      return state;
  }
}

export default issuesReducer;
