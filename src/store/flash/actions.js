export const SET_FLASH_MESSAGE = 'SET_FLASH_MESSAGE';
export const CLEAR_FLASH_MESSAGE = 'CLEAR_FLASH_MESSAGE';

export const setFlashMessage = (key, messageType, message, title = '', slim = false) => ({
  type: SET_FLASH_MESSAGE,
  key,
  messageType,
  message,
  title,
  slim,
});

export const clearFlashMessage = (key = '') => ({
  type: CLEAR_FLASH_MESSAGE,
  key,
});
