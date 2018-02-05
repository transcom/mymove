import { CreateIssue } from 'shared/api.js';

// Types
export const CREATE_ISSUE = 'CREATE_ISSUE';
export const CREATE_ISSUE_SUCCESS = 'CREATE_ISSUE_SUCCESS';
export const CREATE_ISSUE_FAILURE = 'CREATE_ISSUE_FAILURE';

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

// Action creator
export function createIssue(value) {
  return function(dispatch, getState) {
    debugger;
    dispatch(createIssueRequest());
    CreateIssue(value)
      .then(item => dispatch(createIssueSuccess(item)))
      .catch(error => dispatch(createIssueFailure(error)));
    // Fix "things" - what value is being passed? This is a part I do not understand.
    // Does anything need to be passed here?
  };
}

// Reducer
export function feedbackReducer(
  state = { value: '', confirmationText: '' },
  action,
) {
  switch (action.type) {
    case CREATE_ISSUE_SUCCESS:
      return { value: action.item, confirmationText: 'Feedback submitted!' }; // need value passed up from child component, not empty string
    case CREATE_ISSUE_FAILURE:
      return { value: action.error, confirmationText: 'Submission error' }; // should this be the same as above, to preserve the value typed in?
    default:
      return state;
  }
}

// export default feedbackReducer;
