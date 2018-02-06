import dd1299Reducer, {
  createLoadSchemaSuccess,
  createLoadSchemaFailure,
} from './ducks';

describe('Reducer', () => {
  it('Should handle LOAD_SCHEMA_SUCCESS', () => {
    const newSchema = {
      spec: { definitions: { CreateForm1299Payload: 'FOO' } },
    };
    const expectedState = {
      schema: newSchema.spec.definitions.CreateForm1299Payload,
      uiSchema: {},
      hasError: false,
    };
    const newState = dd1299Reducer(
      undefined,
      createLoadSchemaSuccess(newSchema),
    );
    expect(newState).toEqual(expectedState);
  });
  it('Should handle LOAD_SCHEMA_FAILURE', () => {
    const expectedState = {
      schema: null,
      uiSchema: {},
      hasError: true,
    };
    const newState = dd1299Reducer(undefined, createLoadSchemaFailure('OH NO'));
    expect(newState).toEqual(expectedState);
  });
});

// TODO: Figure out how to mock the Swagger API call
// describe('async action creators', () => {
//   const middlewares = [ thunk ]
//   const initialState = { issues: null, hasError: false };
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
