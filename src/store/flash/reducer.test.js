import flashReducer, { initialState } from './reducer';
import { setFlashMessage, clearFlashMessage } from './actions';

describe('flashReducer', () => {
  it('returns the initial state by default', () => {
    expect(flashReducer(undefined, undefined)).toEqual(initialState);
  });

  it('handles the setFlashMessage action', () => {
    expect(flashReducer(initialState, setFlashMessage('success', 'test message', 'Success!'))).toEqual({
      ...initialState,
      flashMessage: {
        title: 'Success!',
        message: 'test message',
        type: 'success',
        key: null,
      },
    });
  });

  it('handles the clearFlashMessage action', () => {
    expect(flashReducer(initialState, clearFlashMessage())).toEqual(initialState);
  });
});
