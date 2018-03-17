import configureStore from 'redux-mock-store';
import {
  feedbackReducer,
  createIssueRequest,
  createIssueSuccess,
  createIssueFailure,
  updateIssueValue,
  createIssue,
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
      hasSubmitError: false,
      hasSubmitSuccess: true,
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
      hasSubmitError: true,
      hasSubmitSuccess: false,
    });
  });
});
