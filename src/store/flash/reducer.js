import { SET_FLASH_MESSAGE, CLEAR_FLASH_MESSAGE } from './actions';

export const initialState = {
  flashMessage: null,
};

const flashReducer = (state = initialState, action) => {
  switch (action?.type) {
    case SET_FLASH_MESSAGE: {
      const { key, messageType, message, title } = action;

      return {
        ...state,
        flashMessage: {
          key,
          type: messageType,
          message,
          title,
        },
      };
    }

    case CLEAR_FLASH_MESSAGE:
      return {
        ...state,
        flashMessage: initialState.flashMessage,
      };

    default:
      return state;
  }
};

export default flashReducer;
