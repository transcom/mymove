import { setFlashMessage, SET_FLASH_MESSAGE, clearFlashMessage, CLEAR_FLASH_MESSAGE } from './actions';

describe('flash actions', () => {
  it('setFlashMessage returns the expected action', () => {
    const expectedAction = {
      type: SET_FLASH_MESSAGE,
      key: 'GENERIC_FLASH_MESSAGE',
      messageType: 'success',
      message: 'Test flash message',
      title: 'Success!',
      slim: true,
    };

    expect(setFlashMessage('GENERIC_FLASH_MESSAGE', 'success', 'Test flash message', 'Success!', true)).toEqual(
      expectedAction,
    );
  });

  it('clearFlashMessage returns the expected action', () => {
    const expectedAction = {
      type: CLEAR_FLASH_MESSAGE,
      key: 'GENERIC_FLASH_MESSAGE',
    };

    expect(clearFlashMessage('GENERIC_FLASH_MESSAGE')).toEqual(expectedAction);
  });
});
