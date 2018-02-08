import { CreateIssue } from 'shared/api.js';

// Types
export const CREATE_ISSUE = 'CREATE_ISSUE';
export const CREATE_ISSUE_SUCCESS = 'CREATE_ISSUE_SUCCESS';
export const CREATE_ISSUE_FAILURE = 'CREATE_ISSUE_FAILURE';
export const CREATE_PENDING_ISSUE_VALUE = 'CREATE_PENDING_ISSUE_VALUE';

// Actions
export const createIssueRequest = () => ({
  type: CREATE_ISSUE,
});

export const createIssueSuccess = item => ({
  type: CREATE_ISSUE_SUCCESS,
  item,
});

export const createIssueFailure = error => ({
  type: CREATE_ISSUE_FAILURE,
  error,
});

export const updateIssueValue = value => ({
  type: CREATE_PENDING_ISSUE_VALUE,
  value,
});

// Action creator
export function createIssue(value) {
  return function(dispatch, getState) {
    dispatch(createIssueRequest());
    CreateIssue(value)
      .then(item => dispatch(createIssueSuccess(item)))
      .catch(error => dispatch(createIssueFailure(error)));
  };
}

export function updatePendingIssueValue(value) {
  return updateIssueValue(value);
}

// Reducer
export function feedbackReducer(
  state = { pendingValue: '', confirmationText: '' },
  action,
) {
  switch (action.type) {
    case CREATE_ISSUE_SUCCESS:
      return {
        pendingValue: '',
        confirmationText: 'Feedback submitted!',
      };
    case CREATE_ISSUE_FAILURE:
      return {
        pendingValue: state.pendingValue,
        confirmationText: 'Submission error.',
      };
    case CREATE_PENDING_ISSUE_VALUE:
      return {
        pendingValue: action.value,
        confirmationText: '',
      };
    default:
      return state;
  }
}

// export default feedbackReducer;
