export const SET_FLASH_MESSAGE = 'SET_FLASH_MESSAGE';
export const CLEAR_FLASH_MESSAGE = 'CLEAR_FLASH_MESSAGE';

export const setFlashMessage = (message, messageType = '', key = null) => ({
  type: SET_FLASH_MESSAGE,
  messageType,
  message,
  key,
});

export const clearFlashMessage = () => ({
  type: CLEAR_FLASH_MESSAGE,
});
