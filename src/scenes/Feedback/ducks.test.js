import { feedbackReducer } from './ducks';

describe('Feedback Reducer', () => {
  it('Should handle CREATE_ISSUE_SUCCESS', () => {
    const initialState = null;

    const newState = feedbackReducer(initialState, {
      type: 'CREATE_ISSUE_SUCCESS',
      item: 'Successful item!',
    });

    expect(newState).toEqual({
      hasErrored: false,
      hasSucceeded: true,
      isLoading: false,
    });
  });

  it('Should handle CREATE_ISSUE_FAILURE', () => {
    const initialState = null;

    const newState = feedbackReducer(initialState, {
      type: 'CREATE_ISSUE_FAILURE',
      error: 'No bueno.',
    });

    expect(newState).toEqual({
      error: 'No bueno.',
      hasErrored: true,
      hasSucceeded: false,
      isLoading: false,
    });
  });
});
