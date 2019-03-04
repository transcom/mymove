// SINGLE RESOURCE ACTION TYPES
const REMOVE_BANNER = 'REMOVE_BANNER';
const SHOW_BANNER = 'SHOW_BANNER';

// SINGLE-RESOURCE ACTION CREATORS

export const removeBanner = () => {
  return {
    type: REMOVE_BANNER,
  };
};

export const showBanner = ({ messageLines }) => {
  return {
    type: SHOW_BANNER,
    payload: { messageLines },
  };
};

// Reducer
const initialState = {
  display: false,
  messageLines: [],
};

export function officeFlashMessagesReducer(state = initialState, action) {
  switch (action.type) {
    // SINGLE-RESOURCE ACTION TYPES
    case SHOW_BANNER:
      return Object.assign({}, state, {
        display: true,
        messageLines: action.payload.messageLines,
      });
    case REMOVE_BANNER:
      return Object.assign({}, state, {
        display: false,
        messageLines: [],
      });

    default:
      return state;
  }
}
