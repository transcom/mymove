import { GetSpec } from 'shared/api';
import { getUiSchema } from './uiSchema';

// Types
export const LOAD_SCHEMA = 'LOAD_SCHEMA';
export const LOAD_SCHEMA_SUCCESS = 'LOAD_SCHEMA_SUCCESS';
export const LOAD_SCHEMA_FAILURE = 'LOAD_SCHEMA_FAILURE';

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

// Reducer
//todo: we may want to reorganize this after we have implemented more reports
// for instance it may make sense to have the whole schema (and uiSchema) in the store and use selectors to get the pieces we need for reports
const initialState = { schema: {}, uiSchema: getUiSchema(), hasError: false };
function dd1299Reducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_SCHEMA_SUCCESS:
      return Object.assign({}, state, {
        schema: action.schema.definitions.CreateForm1299Payload,
        hasError: false,
      });
    case LOAD_SCHEMA_FAILURE:
      return Object.assign({}, state, { schema: {}, hasError: true });
    default:
      return state;
  }
}

export default dd1299Reducer;
