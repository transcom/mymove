import { SET_FLASH_MESSAGE, CLEAR_FLASH_MESSAGE } from './actions';

export const initialState = {
  flashMessage: null,
};

const flashReducer = (state = initialState, action) => {
  switch (action?.type) {
    case SET_FLASH_MESSAGE: {
      const { key, messageType, message, title, slim } = action;

      return {
        ...state,
        flashMessage: {
          key,
          type: messageType,
          message,
          title,
          slim,
        },
      };
    }

    case CLEAR_FLASH_MESSAGE: {
      const { key } = action;

      if (key && state.flashMessage?.key === key) {
        return {
          ...state,
          flashMessage: initialState.flashMessage,
        };
      }

      return state;
    }

    default:
      return state;
  }
};

export default flashReducer;
