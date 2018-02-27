import { GetSpec, CreateForm1299 } from './api';
import { getUiSchema } from './uiSchema';

// Types
export const LOAD_SCHEMA = 'LOAD_SCHEMA';
export const LOAD_SCHEMA_SUCCESS = 'LOAD_SCHEMA_SUCCESS';
export const LOAD_SCHEMA_FAILURE = 'LOAD_SCHEMA_FAILURE';

export const REQUEST_SUBMIT = 'REQUEST_SUBMIT';
export const SUBMIT_SUCCESS = 'SUBMIT_SUCCESS';
export const SUBMIT_FAILURE = 'SUBMIT_FAILURE';
export const SUBMIT_RESET = 'SUBMIT_RESET';

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

export const createRequestSubmit = () => ({
  type: REQUEST_SUBMIT,
});

//submitting form
export const createSubmitSuccess = responseData => ({
  type: SUBMIT_SUCCESS,
  responseData,
});

export const createSubmitFailure = error => ({
  type: SUBMIT_FAILURE,
  error,
});

export const createSubmitReset = () => ({
  type: SUBMIT_RESET,
});

// Action Creator
export function loadSchema() {
  // Interpreted by the thunk middleware:
  return function(dispatch) {
    dispatch(createLoadSchemaRequest());
    GetSpec()
      .then(spec => dispatch(createLoadSchemaSuccess(spec)))
      .catch(error => dispatch(createLoadSchemaFailure(error)));
  };
}

export function submitForm(formData) {
  return function(dispatch, getState) {
    if (!formData) {
      // HACK: since we are using redux-thunk, have access to other state
      formData = getState().form.DD1299.values;
    }
    dispatch(createRequestSubmit());
    CreateForm1299(formData)
      .then(result => dispatch(createSubmitSuccess(result)))
      .catch(error => dispatch(createSubmitFailure(error)));
  };
}

export function resetSuccess() {
  return createSubmitReset();
}
// Reducer
//todo: we may want to reorganize this after we have implemented more reports
// for instance it may make sense to have the whole schema (and uiSchema) in the store and use selectors to get the pieces we need for reports
const initialState = {
  schema: {},
  uiSchema: getUiSchema(),
  hasSchemaError: false,
  hasSubmitError: false,
  hasSubmitSuccess: false,
};
function dd1299Reducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_SCHEMA_SUCCESS:
      return Object.assign({}, state, {
        schema: action.schema.definitions.CreateForm1299Payload,
        hasSchemaError: false,
      });
    case LOAD_SCHEMA_FAILURE:
      return Object.assign({}, state, {
        schema: {},
        hasSchemaError: true,
      });
    case SUBMIT_SUCCESS:
      return Object.assign({}, state, {
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case SUBMIT_FAILURE:
      return Object.assign({}, state, {
        hasSubmitSuccess: false,
        hasSubmitError: true,
      });
    case SUBMIT_RESET:
      return Object.assign({}, state, {
        hasSubmitError: false,
        hasSubmitSuccess: false,
      });
    default:
      return state;
  }
}

export default dd1299Reducer;
