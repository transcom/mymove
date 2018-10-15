import configureStore from 'redux-mock-store';
import issuesReducer, { createShowIssuesRequest, createShowIssuesSuccess, createShowIssuesFailure } from './ducks';

jest.mock('./api');

describe('Issues Reducer', () => {
  it('Should handle SHOW_ISSUES', () => {
    const initialState = { issues: null, hasError: false };

    const newState = issuesReducer(initialState, { type: 'SHOW_ISSUES' });

    expect(newState).toEqual({ issues: null, hasError: false });
  });

  it('Should handle SHOW_ISSUES_SUCCESS', () => {
    const initialState = { issues: null, hasError: false };

    const newState = issuesReducer(initialState, {
      type: 'SHOW_ISSUES_SUCCESS',
      items: 'TOO FEW DOGS',
    });

    expect(newState).toEqual({ issues: 'TOO FEW DOGS', hasError: false });
  });

  it('Should handle SHOW_ISSUES_FAILURE', () => {
    const initialState = { issues: null, hasError: false };

    const newState = issuesReducer(initialState, {
      type: 'SHOW_ISSUES_FAILURE',
      error: 'Boring',
    });

    expect(newState).toEqual({ issues: null, hasError: true });
  });
});

describe('Issues Actions', () => {
  const initialState = { issues: null, hasError: false };
  const mockStore = configureStore();
  let store;

  beforeEach(() => {
    store = mockStore(initialState);
  });

  it('Should check action on dispatching ', () => {
    let action;
    store.dispatch(createShowIssuesRequest());
    store.dispatch(createShowIssuesSuccess([{ id: '11', description: 'too few dogs' }]));
    store.dispatch(createShowIssuesFailure('Tests r not fun.'));
    action = store.getActions();
    // Add expect about what the contents will be.
    expect(action[0].type).toBe('SHOW_ISSUES');
    expect(action[1].type).toBe('SHOW_ISSUES_SUCCESS');
    expect(action[1].items).toEqual([{ id: '11', description: 'too few dogs' }]);
    expect(action[2].type).toBe('SHOW_ISSUES_FAILURE');
    expect(action[2].error).toEqual('Tests r not fun.');
  });
});

// TODO: Figure out how to mock the Swagger API call
describe('given there are issues, when loadIssues is called', () => {
  // const middlewares = [thunk];
  // const initialState = { issues: null, hasError: false };
  // const mockStore = configureStore(middlewares);

  it('it creates SHOW_ISSUES_SUCCESS when submitted issues have been loaded and SHOW_ISSUES_SUCCESS payload is those issues', () => {
    // const expectedActions = [
    //   { type: SHOW_ISSUES },
    //   {
    //     type: SHOW_ISSUES_SUCCESS,
    //     items: { issues: [{ id: 11, description: 'too few dogs' }] },
    //   },
    // ];
    // const store = mockStore(initialState);
    // return store.dispatch(loadIssues()).then(() => {
    //   // return of async actions
    //   expect(store.getActions()).toEqual(expectedActions);
    // });
  });
});
