import { GetSpec, CreateIssue } from './api.js';
import { getUiSchema } from './uiSchema';

// Types
export const LOAD_SCHEMA = 'LOAD_SCHEMA';
export const LOAD_SCHEMA_SUCCESS = 'LOAD_SCHEMA_SUCCESS';
export const LOAD_SCHEMA_FAILURE = 'LOAD_SCHEMA_FAILURE';

export const CREATE_ISSUE = 'CREATE_ISSUE';
export const CREATE_ISSUE_SUCCESS = 'CREATE_ISSUE_SUCCESS';
export const CREATE_ISSUE_FAILURE = 'CREATE_ISSUE_FAILURE';

// Actions
// loading schema
export const createLoadSchemaRequest = () => ({
  type: LOAD_SCHEMA,
});

export const createLoadSchemaSuccess = schema => ({
  type: LOAD_SCHEMA_SUCCESS,
  schema,
});

export const createLoadSchemaFailure = error => ({
  type: LOAD_SCHEMA_FAILURE,
  error,
});

// creating issue
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
export function loadSchema() {
  // Interpreted by the thunk middleware:
  return function(dispatch) {
    dispatch(createLoadSchemaRequest());
    GetSpec()
      .then(spec => dispatch(createLoadSchemaSuccess(spec)))
      .catch(error => dispatch(createLoadSchemaFailure(error)));
  };
}

export function createIssue(value) {
  return function(dispatch, getState) {
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
    case LOAD_SCHEMA_SUCCESS:
      console.log('WOEINOWNFOWIENFOIWENFIIIIII');
      return Object.assign({}, state, {
        schema: action.schema.definitions.CreateIssuePayload,
        hasSchemaError: false,
      });
    case LOAD_SCHEMA_FAILURE:
      return Object.assign({}, state, {
        schema: {},
        hasSchemaError: true,
      });
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
