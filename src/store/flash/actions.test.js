import { setFlashMessage, SET_FLASH_MESSAGE, clearFlashMessage, CLEAR_FLASH_MESSAGE } from './actions';

describe('flash actions', () => {
  it('setFlashMessage returns the expected action', () => {
    const expectedAction = {
      type: SET_FLASH_MESSAGE,
      message: 'Test flash message',
      messageType: 'success',
      key: 'GENERIC_FLASH_MESSAGE',
    };

    expect(setFlashMessage('Test flash message', 'success', 'GENERIC_FLASH_MESSAGE')).toEqual(expectedAction);
  });

  it('clearFlashMessage returns the expected action', () => {
    const expectedAction = {
      type: CLEAR_FLASH_MESSAGE,
    };

    expect(clearFlashMessage()).toEqual(expectedAction);
  });
});
