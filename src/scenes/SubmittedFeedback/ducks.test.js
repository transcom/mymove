import React from 'react';
import ReactDOM from 'react-dom';
import configureStore from 'redux-mock-store';
import { shallow } from 'enzyme';
import SubmittedFeedback from '.';
import issuesReducer, {
  createShowIssuesRequest,
  createShowIssuesSuccess,
  createShowIssuesFailure,
} from './ducks';

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
    store.dispatch(createShowIssuesSuccess());
    store.dispatch(createShowIssuesFailure());
    action = store.getActions();
    expect(action[0].type).toBe('SHOW_ISSUES');
    expect(action[1].type).toBe('SHOW_ISSUES_SUCCESS');
    expect(action[2].type).toBe('SHOW_ISSUES_FAILURE');
  });
});
