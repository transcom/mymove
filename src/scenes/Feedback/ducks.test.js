import configureStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import {
  feedbackReducer,
  createIssueRequest,
  createIssueSuccess,
  createIssueFailure,
  updateIssueValue,
} from './ducks';

describe('Feedback Reducer', () => {
  it('Should handle CREATE_ISSUE_SUCCESS', () => {
    const initialState = { pendingValue: '', confirmationText: '' };

    const newState = feedbackReducer(initialState, {
      type: 'CREATE_ISSUE_SUCCESS',
      item: 'Successful item!',
    });

    expect(newState).toEqual({
      pendingValue: '',
      confirmationText: 'Feedback submitted!',
    });
  });

  it('Should handle CREATE_ISSUE_FAILURE', () => {
    const initialState = { pendingValue: '', confirmationText: '' };

    const newState = feedbackReducer(initialState, {
      type: 'CREATE_ISSUE_FAILURE',
      error: 'No bueno.',
    });

    expect(newState).toEqual({
      pendingValue: '',
      confirmationText: 'Submission error.',
    });
  });

  it('Should handle CREATE_PENDING_ISSUE_VALUE', () => {
    const initialState = { pendingValue: '', confirmationText: '' };

    const newState = feedbackReducer(initialState, {
      type: 'CREATE_PENDING_ISSUE_VALUE',
      value: 'asd',
    });

    expect(newState).toEqual({ pendingValue: 'asd', confirmationText: '' });
  });
});
