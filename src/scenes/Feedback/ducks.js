import { CreateIssue } from 'shared/api.js';

// Types
export const CREATE_ISSUE = 'CREATE_ISSUE';
export const CREATE_ISSUE_SUCCESS = 'CREATE_ISSUE_SUCCESS';
export const CREATE_ISSUE_FAILURE = 'CREATE_ISSUE_FAILURE';

// Actions
export const createIssueRequest = () => ({
  type: CREATE_ISSUE,
});

export const createIssueSuccess = () => ({
  type: CREATE_ISSUE_SUCCESS,
});

export const createIssueFailure = () => ({
  type: CREATE_ISSUE_FAILURE,
});

// Action creator
export function createIssue() {
  return function(dispatch, getState) {
    dispatch(createIssueRequest());
    CreateIssue(this.props.value)
      .then(thing => dispatch(createIssueSuccess(thing)))
      .catch(errorThing => dispatch(createIssueFailure(errorThing)));
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
      return { value: '', confirmationText: 'Feedback submitted!' }; // need value passed up from child component, not empty string
    case CREATE_ISSUE_FAILURE:
      return { value: '', confirmationText: 'Submission error' }; // should this be the same as above, to preserve the value typed in?
    default:
      return state;
  }
}

// export default feedbackReducer;
