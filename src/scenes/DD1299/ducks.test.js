import dd1299Reducer, {
  createLoadSchemaSuccess,
  createLoadSchemaFailure,
  createCreateFailure,
  createCreateSuccess,
  createCreateReset,
} from './ducks';
import { getUiSchema } from './uiSchema';

const uiSchema = getUiSchema();
describe('Reducer', () => {
  it('Should handle LOAD_SCHEMA_SUCCESS', () => {
    const newSchema = {
      definitions: { CreateForm1299Payload: 'FOO' },
    };
    const expectedState = {
      schema: newSchema.definitions.CreateForm1299Payload,
      uiSchema,
      hasSchemaError: false,
      hasCreateError: false,
      hasCreateSuccess: false,
    };
    const newState = dd1299Reducer(
      undefined,
      createLoadSchemaSuccess(newSchema),
    );
    expect(newState).toEqual(expectedState);
  });
  it('Should handle LOAD_SCHEMA_FAILURE', () => {
    const err = 'OH NO';
    const expectedState = {
      schema: {},
      uiSchema,
      hasSchemaError: true,
      hasCreateError: false,
      hasCreateSuccess: false,
    };
    const newState = dd1299Reducer(undefined, createLoadSchemaFailure(err));
    expect(newState).toEqual(expectedState);
  });
  it('Should handle CREATE_SUCCESS', () => {
    const err = 'OH NO';
    const expectedState = {
      schema: {},
      uiSchema,
      hasSchemaError: false,
      hasCreateError: false,
      hasCreateSuccess: true,
    };
    const newState = dd1299Reducer(undefined, createCreateSuccess(err));
    expect(newState).toEqual(expectedState);
  });
  it('Should handle CREATE_FAILURE', () => {
    const err = 'OH NO';
    const expectedState = {
      schema: {},
      uiSchema,
      hasSchemaError: false,
      hasCreateError: true,
      hasCreateSuccess: false,
    };
    const newState = dd1299Reducer(undefined, createCreateFailure(err));
    expect(newState).toEqual(expectedState);
  });
  it('Should handle CREATE_RESET', () => {
    const err = 'OH NO';
    const expectedState = {
      schema: {},
      uiSchema,
      hasSchemaError: false,
      hasCreateError: false,
      hasCreateSuccess: false,
    };
    const newState = dd1299Reducer(undefined, createCreateReset());
    expect(newState).toEqual(expectedState);
  });
});

// TODO: Figure out how to mock the Swagger API call
// describe('async action creators', () => {
//   const middlewares = [ thunk ]
//   const initialState = { issues: null, hasSchemaError: false };
//   const mockStore = configureStore(middlewares)

//   afterEach(() => {
//     fetchMock.reset()
//     fetchMock.restore()
//   })

//   it('creates SHOW_ISSUES_SUCCESS when submitted issues have been loaded', () => {
//     fetchMock
//       .getOnce('/submitted', { items: { issues: [{'id': 11, 'description': 'too few dogs'}] }, headers: { 'content-type': 'application/json' } })

//     const expectedActions = [
//       { type: SHOW_ISSUES },
//       { type: SHOW_ISSUES_SUCCESS, items: { issues: [{'id': 11, 'description':'too few dogs'}] } }
//     ]

//     const store = mockStore(initialState)

//     return store.dispatch(loadIssues()).then(() => {
//       // return of async actions
//       expect(store.getActions()).toEqual(expectedActions)
//     })
//   })
// })
