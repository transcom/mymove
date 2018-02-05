import { GetSpec } from 'shared/api';

// Types
export const LOAD_SCHEMA = 'LOAD_SCHEMA';
export const LOAD_SCHEMA_SUCCESS = 'LOAD_SCHEMA_SUCCESS';
export const LOAD_SCHEMA_FAILURE = 'LOAD_SCHEMA_FAILURE';

// Actions
export const createLoadSchemaRequest = () => ({
  type: LOAD_SCHEMA,
});

export const createLoadSchemaSuccess = spec => ({
  type: LOAD_SCHEMA_SUCCESS,
  spec,
});

export const createLoadSchemaFailure = error => ({
  type: LOAD_SCHEMA_FAILURE,
  error,
});

// Action Creator
export function loadSchema() {
  // Interpreted by the thunk middleware:
  return function(dispatch, getState) {
    dispatch(createLoadSchemaRequest());
    GetSpec()
      .then(spec => dispatch(createLoadSchemaSuccess(spec)))
      .catch(error => dispatch(createLoadSchemaFailure(error)));
  };
}

// Reducer
function dd1299Reducer(state = { schema: {}, hasError: false }, action) {
  switch (action.type) {
    case LOAD_SCHEMA_SUCCESS:
      return {
        schema: action.spec.definitions.CreateForm1299Payload,
        hasError: false,
      };
    case LOAD_SCHEMA_FAILURE:
      return { schema: {}, hasError: true };
    default:
      return state;
  }
}

export default dd1299Reducer;
