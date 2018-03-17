import { CreateIssue } from './api.js';
import { getUiSchema } from './uiSchema';

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
  return function(dispatch) {
    dispatch(createIssueRequest());
    CreateIssue(value)
      .then(item => dispatch(createIssueSuccess(item)))
      .catch(error => dispatch(createIssueFailure(error)));
  };
}

// Reducer
const initialState = {
  schema: {},
  uiSchema: getUiSchema(),
  hasSchemaError: false,
  hasSubmitError: false,
  hasSubmitSuccess: false,
  confirmationText: '',
};
export function feedbackReducer(state = initialState, action) {
  switch (action.type) {
    case CREATE_ISSUE_SUCCESS:
      return Object.assign({}, state, {
        hasSubmitSuccess: true,
        hasSubmitError: false,
        confirmationText: 'Feedback submitted!',
      });
    case CREATE_ISSUE_FAILURE:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
        confirmationText: 'Submission error.',
      });
    default:
      return state;
  }
}

// export default feedbackReducer;
