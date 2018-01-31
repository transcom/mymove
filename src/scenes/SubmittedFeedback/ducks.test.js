import React from 'react';
import ReactDOM from 'react-dom';
import { shallow } from 'enzyme';
import SubmittedFeedback from '.';
import issuesReducer from './ducks';

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
