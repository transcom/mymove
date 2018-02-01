import { CreateIssue } from 'shared/api.js';

// Types
export const CREATE_ISSUE = 'CREATE_ISSUE';
export const CREATE_ISSUE_SUCCESS = 'CREATE_ISSUE_SUCCESS';
export const CREATE_ISSUE_FAILURE = 'CREATE_ISSUE_FAILURE';

// Actions
export const createIssueRequest = () = ({
  type: CREATE_ISSUE,
});

export const createIssueSuccess = () = ({
  type: CREATE_ISSUE_SUCCESS,
});

export const createIssueFailure = () = ({
  type: CREATE_ISSUE_FAILURE,
});

// Action creator
export function createIssue() {
  return function(dispatch, getState) {
    dispatch(createIssueRequest());
    CreateIssue(this.props.value)
      .then(thing => dispatch(createShowIssuesSuccess(thing)))
      .catch(errorThing => dispatch(createIssueFailure(errorThing)));
      // Fix "things" - what value is being passed? This is a part I do not understand.
  }
}

// Reducer
function feedbackReducer(state = { value: '', confirmationText: '' }, action) {
  switch (action.type) {
    case CREATE_ISSUE_SUCCESS:
      return { value: '', confirmationText: 'text of confirmation' };
    case CREATE_ISSUE_FAILURE:
      return { value: '', confirmationText: 'text of failure' };
    default:
      return state;
  }
}

export default feedbackReducer;
