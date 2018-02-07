import { GetSpec, CreateForm1299 } from 'shared/api';
import { getUiSchema } from './uiSchema';

// Types
export const LOAD_SCHEMA = 'LOAD_SCHEMA';
export const LOAD_SCHEMA_SUCCESS = 'LOAD_SCHEMA_SUCCESS';
export const LOAD_SCHEMA_FAILURE = 'LOAD_SCHEMA_FAILURE';

export const REQUEST_CREATE = 'REQUEST_CREATE';
export const CREATE_SUCCESS = 'CREATE_SUCCESS';
export const CREATE_FAILURE = 'CREATE_FAILURE';
export const CREATE_RESET = 'CREATE_RESET';

// Actions
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

export const createRequestCreate = () => ({
  type: REQUEST_CREATE,
});
export const createCreateSuccess = responseData => ({
  type: CREATE_SUCCESS,
  responseData,
});

export const createCreateFailure = error => ({
  type: CREATE_FAILURE,
  error,
});

export const createCreateReset = () => ({
  type: CREATE_RESET,
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

export function createForm(formData) {
  return function(dispatch) {
    dispatch(createRequestCreate());
    CreateForm1299(formData)
      .then(result => dispatch(createCreateSuccess(result)))
      .catch(error => dispatch(createCreateFailure(error)));
  };
}

export function resetSuccess() {
  return createCreateReset();
}
// Reducer
//todo: we may want to reorganize this after we have implemented more reports
// for instance it may make sense to have the whole schema (and uiSchema) in the store and use selectors to get the pieces we need for reports
const initialState = {
  schema: {},
  uiSchema: getUiSchema(),
  hasError: false,
  hasCreateError: false,
  hasCreateSuccess: false,
};
function dd1299Reducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_SCHEMA_SUCCESS:
      return Object.assign({}, state, {
        schema: action.schema.definitions.CreateForm1299Payload,
        hasError: false,
      });
    case LOAD_SCHEMA_FAILURE:
      return Object.assign({}, state, { schema: {}, hasError: true });
    case CREATE_SUCCESS:
      return Object.assign({}, state, { hasCreateSuccess: true });
    case CREATE_FAILURE:
      return Object.assign({}, state, {
        hasCreateSuccess: false,
        hasCreateError: true,
      });
    case CREATE_RESET:
      return Object.assign({}, state, {
        hasCreateError: false,
        hasCreateSuccess: false,
      });
    default:
      return state;
  }
}

export default dd1299Reducer;
