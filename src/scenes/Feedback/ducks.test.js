import configureStore from 'redux-mock-store';
// import thunk from 'redux-thunk';
import {
  feedbackReducer,
  createIssueRequest,
  createIssueSuccess,
  createIssueFailure,
  updateIssueValue,
  createIssue,
} from './ducks';

// describe('Feedback Reducer', () => {
//   it('Should handle CREATE_ISSUE_SUCCESS', () => {
//     const initialState = { pendingValue: '', confirmationText: '' };

//     const newState = feedbackReducer(initialState, {
//       type: 'CREATE_ISSUE_SUCCESS',
//       item: 'Successful item!',
//     });

//     expect(newState).toEqual({
//       pendingValue: '',
//       confirmationText: 'Feedback submitted!',
//     });
//   });

//   it('Should handle CREATE_ISSUE_FAILURE', () => {
//     const initialState = { pendingValue: '', confirmationText: '' };

//     const newState = feedbackReducer(initialState, {
//       type: 'CREATE_ISSUE_FAILURE',
//       error: 'No bueno.',
//     });

//     expect(newState).toEqual({
//       pendingValue: '',
//       confirmationText: 'Submission error.',
//     });
//   });

//   it('Should handle CREATE_PENDING_ISSUE_VALUE', () => {
//     const initialState = { pendingValue: '', confirmationText: '' };

//     const newState = feedbackReducer(initialState, {
//       type: 'CREATE_PENDING_ISSUE_VALUE',
//       value: 'asd',
//     });

//     expect(newState).toEqual({ pendingValue: 'asd', confirmationText: '' });
//   });
// });

describe('Feedback actions', () => {
  const initialState = { pendingValue: '', confirmationText: '' };
  const mockStore = configureStore();
  let store;

  beforeEach(() => {
    store = mockStore(initialState);
  });

  it('Should check action on dispatching ', () => {
    let action;
    store.dispatch(createIssue());
    store.dispatch(createIssueSuccess('Why is this site so good though')); // Did this need to be an array? Doesn't seem to affect anything right now...
    store.dispatch(createIssueFailure('Oh dear.'));
    console.log('Check, check'); // doesn't log
    action = store.getActions();
    console.log('Hey, actions', action); // doesn't log

    expect(action[0].type).toBe('CREATE_ISSUE');
    expect(action[1].type).toBe('CREATE_ISSUE_SUCCESS');
    expect(action[1].item).toEqual([
      { pendingValue: 'Why is this site so good though' },
    ]);
    expect(action[2].type).toBe('CREATE_ISSUE_FAILURE');
    expect(action[2].error).toEqual('Oh dear.');
  });
});
