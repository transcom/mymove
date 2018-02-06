import { GetSpec } from 'shared/api';
import { getUiSchema } from './uiSchema';

// Types
export const LOAD_SCHEMA = 'LOAD_SCHEMA';
export const LOAD_SCHEMA_SUCCESS = 'LOAD_SCHEMA_SUCCESS';
export const LOAD_SCHEMA_FAILURE = 'LOAD_SCHEMA_FAILURE';
export const LOAD_UI_SCHEMA = 'LOAD_UI_SCHEMA';

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

export const createLoadUiSchema = uiSchema => ({
  type: LOAD_UI_SCHEMA,
  uiSchema,
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

export function loadUiSchema() {
  return createLoadUiSchema(getUiSchema());
}
// Reducer
//todo: we may want to reorganize this after we have implemented more reports
// for instance it may make sense to have the whole schema (and uiSchema) in the store and use selectors to get the pieces we need for reports
const initialState = { schema: null, uiSchema: {}, hasError: false };
function dd1299Reducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_SCHEMA_SUCCESS:
      return Object.assign({}, state, {
        schema: action.schema.spec.definitions.CreateForm1299Payload,
        hasError: false,
      });
    case LOAD_SCHEMA_FAILURE:
      return Object.assign({}, state, { schema: null, hasError: true });
    case LOAD_UI_SCHEMA:
      return Object.assign({}, state, { uiSchema: action.uiSchema });
    default:
      return state;
  }
}

export default dd1299Reducer;
