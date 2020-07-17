import * as helpers from 'shared/ReduxHelpers';
import { GetLoggedInUser, GetIsLoggedIn } from 'shared/User/api';
import { pick } from 'lodash';
import { normalize } from 'normalizr';
import { ordersArray } from 'shared/Entities/schema';
import { addEntities } from 'shared/Entities/actions';

export const setIsLoggedInType = 'SET_IS_LOGGED_IN';
const getLoggedInUserType = 'GET_LOGGED_IN_USER';

export const GET_LOGGED_IN_USER = helpers.generateAsyncActionTypes(getLoggedInUserType);
export const getLoggedInActions = helpers.generateAsyncActions(getLoggedInUserType);

export function getCurrentUserInfo() {
  return function (dispatch) {
    dispatch(getLoggedInActions.start());
    return GetIsLoggedIn()
      .then((response) => {
        if (response.isLoggedIn === true) {
          return GetLoggedInUser()
            .then((response) => {
              if (response.service_member) {
                const data = normalize(response.service_member.orders, ordersArray);

                // Only store addresses in a normalized way. This prevents
                // data duplication while we're using both Redux approaches.
                const filtered = pick(data.entities, ['addresses']);
                dispatch(addEntities(filtered));
              }
              return dispatch(getLoggedInActions.success(response));
            })
            .catch((error) => dispatch(getLoggedInActions.error(error)));
        }
        return dispatch(getLoggedInActions.error('User is not logged in'));
      })
      .catch((error) => dispatch(getLoggedInActions.error(error)));
  };
}

export function setUserIsLoggedIn(isLoggedIn) {
  return function (dispatch) {
    return dispatch({ type: setIsLoggedInType, isLoggedIn });
  };
}

export function selectCurrentUser(state) {
  return state.user.userInfo || {};
}

export function selectGetCurrentUserIsLoading(state) {
  return state.user.isLoading;
}

export function selectGetCurrentUserIsSuccess(state) {
  return state.user.hasSucceeded;
}

export function selectGetCurrentUserIsError(state) {
  return state.user.hasErrored;
}

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
    case setIsLoggedInType:
      return {
        ...state,
        userInfo: {
          ...userInfoDefault(),
          isLoggedIn: action.isLoggedIn,
        },
      };
    default:
      return state;
  }
};

export default currentUserReducer;
