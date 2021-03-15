import flashReducer, { initialState } from './reducer';
import { setFlashMessage, clearFlashMessage } from './actions';

describe('flashReducer', () => {
  it('returns the initial state by default', () => {
    expect(flashReducer(undefined, undefined)).toEqual(initialState);
  });

  it('handles the setFlashMessage action', () => {
    expect(flashReducer(initialState, setFlashMessage('TEST_SUCCESS', 'success', 'test message', 'Success!'))).toEqual({
      ...initialState,
      flashMessage: {
        key: 'TEST_SUCCESS',
        type: 'success',
        message: 'test message',
        title: 'Success!',
        slim: false,
      },
    });
  });

  it('handles the clearFlashMessage action', () => {
    expect(flashReducer(initialState, clearFlashMessage('TEST_SUCCESS'))).toEqual(initialState);
  });
});
