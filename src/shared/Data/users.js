import * as helpers from 'shared/ReduxHelpers';
import { GetLoggedInUser } from 'shared/User/api.js';
import { pick } from 'lodash';
import { normalize } from 'normalizr';
import { ordersArray } from 'shared/Entities/schema';
import { addEntities } from 'shared/Entities/actions';
import { getPublicShipment } from 'shared/Entities/modules/shipments';
import store from 'shared/store';

const getLoggedInUserType = 'GET_LOGGED_IN_USER';

export const GET_LOGGED_IN_USER = helpers.generateAsyncActionTypes(getLoggedInUserType);
const getLoggedInActions = helpers.generateAsyncActions(getLoggedInUserType);

export function getCurrentUserInfo() {
  return function(dispatch) {
    dispatch(getLoggedInActions.start());
    return GetLoggedInUser()
      .then(response => {
        if (response.service_member) {
          const data = normalize(response.service_member.orders, ordersArray);
          if (data.entities.shipments) {
            const shipmentIds = Object.keys(data.entities.shipments);
            shipmentIds.map(id => dispatch(getPublicShipment(id)));
          }

          // Only store addresses in a normalized way. This prevents
          // data duplication while we're using both Redux approaches.
          const filtered = pick(data.entities, ['addresses']);
          dispatch(addEntities(filtered));
        }
        return dispatch(getLoggedInActions.success(response));
      })
      .catch(error => dispatch(getLoggedInActions.error(error)));
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

export const SET_LOGGED_OUT = 'SET_LOGGED_OUT';
export const setLoggedOut = () => store.dispatch({ type: SET_LOGGED_OUT });

const currentUserReducerDefault = () => ({
  hasSucceeded: false,
  hasErrored: false,
  isLoading: false,
  userInfo: { email: '', isLoggedIn: false },
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
    case GET_LOGGED_IN_USER.error:
      return {
        ...state,
        isLoading: false,
        hasErrored: true,
        hasSucceeded: false,
        error: action.error,
      };
    case SET_LOGGED_OUT:
      return {
        ...state,
        userInfo: { email: '', isLoggedIn: false },
        hasSucceeded: false,
        hasErrored: false,
        isLoading: false,
      };
    default:
      return state;
  }
};

export default currentUserReducer;
