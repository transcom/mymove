export const SET_FLASH_MESSAGE = 'SET_FLASH_MESSAGE';
export const CLEAR_FLASH_MESSAGE = 'CLEAR_FLASH_MESSAGE';

export const setFlashMessage = (key, messageType, message, title = '') => ({
  type: SET_FLASH_MESSAGE,
  key,
  messageType,
  message,
  title,
});

export const clearFlashMessage = () => ({
  type: CLEAR_FLASH_MESSAGE,
});
