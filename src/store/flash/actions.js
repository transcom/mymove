export const SET_FLASH_MESSAGE = 'SET_FLASH_MESSAGE';
export const CLEAR_FLASH_MESSAGE = 'CLEAR_FLASH_MESSAGE';

export const setFlashMessage = (messageType, message, title = '', key = null) => ({
  type: SET_FLASH_MESSAGE,
  messageType,
  title,
  message,
  key,
});

export const clearFlashMessage = () => ({
  type: CLEAR_FLASH_MESSAGE,
});
