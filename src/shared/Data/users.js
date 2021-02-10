import * as helpers from 'shared/ReduxHelpers';

const getLoggedInUserType = 'GET_LOGGED_IN_USER';

export const GET_LOGGED_IN_USER = helpers.generateAsyncActionTypes(getLoggedInUserType);
export const getLoggedInActions = helpers.generateAsyncActions(getLoggedInUserType);

const userInfoDefault = () => ({
  email: '',
  isLoggedIn: false,
});

const currentUserReducerDefault = () => ({
  hasSucceeded: false,
  hasErrored: false,
  isLoading: true,
  userInfo: userInfoDefault(),
});

const currentUserReducer = (state = currentUserReducerDefault(), action) => {
  switch (action.type) {
    case GET_LOGGED_IN_USER.start:
      return {
        ...state,
        hasSucceeded: false,
        hasErrored: false,
        isLoading: true,
      };
    case GET_LOGGED_IN_USER.success:
      return {
        ...state,
        userInfo: {
          isLoggedIn: true,
          ...action.payload,
        },
        hasSucceeded: true,
        hasErrored: false,
        isLoading: false,
      };
    case GET_LOGGED_IN_USER.failure:
      return {
        ...state,
        isLoading: false,
        hasErrored: true,
        hasSucceeded: false,
        error: action.error,
        userInfo: userInfoDefault(),
      };
    default:
      return state;
  }
};

export default currentUserReducer;
