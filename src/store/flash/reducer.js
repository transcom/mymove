import { SET_FLASH_MESSAGE, CLEAR_FLASH_MESSAGE } from './actions';

export const initialState = {
  flashMessage: {
    type: null,
    message: null,
    key: null,
  },
};

const flashReducer = (state = initialState, action) => {
  switch (action?.type) {
    case SET_FLASH_MESSAGE: {
      const { message, messageType, key } = action;

      return {
        ...state,
        flashMessage: {
          type: messageType,
          message,
          key,
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
