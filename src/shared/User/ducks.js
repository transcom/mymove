import * as Cookies from 'js-cookie';
import * as decode from 'jwt-decode';
import * as helpers from 'shared/ReduxHelpers';
import { GetLoggedInUser } from './api.js';

const LOAD_USER_AND_TOKEN = 'USER|LOAD_USER_AND_TOKEN';

const loggedOutUser = {
  isLoggedIn: false,
  email: null,
  jwt: null,
};

const GET_LOGGED_IN_USER = 'GET_LOGGED_IN_USER';

export const getUserTypes = helpers.generateAsyncActionTypes(
  GET_LOGGED_IN_USER,
);

async function GetLoggedInUserIfLoggedIn() {
  const userInfo = getUserInfo();
  if (userInfo.isLoggedIn) await GetLoggedInUser();
}
export const loadLoggedInUser = helpers.generateAsyncActionCreator(
  GET_LOGGED_IN_USER,
  GetLoggedInUserIfLoggedIn,
);

export const loggedInUserReducer = helpers.generateAsyncReducer(
  GET_LOGGED_IN_USER,
  u => ({ loggedInUser: u }),
);

function getUserInfo() {
  const cookie = Cookies.get('user_session');
  if (!cookie) return loggedOutUser;
  const jwt = decode(cookie);
  //if (jwt.exp <  Date.now().valueOf() / 1000) return loggedOutUser;
  return {
    jwt: cookie,
    email: jwt.email,
    userId: jwt.user_id,
    expires: Date(jwt.exp),
    isLoggedIn: true,
  };
}

export function loadUserAndToken() {
  const jwt = getUserInfo();
  return { type: LOAD_USER_AND_TOKEN, payload: jwt };
}

const userReducer = (state = getUserInfo(), action) => {
  switch (action.type) {
    case LOAD_USER_AND_TOKEN:
      return action.payload;
    default: {
      return state;
    }
  }
};

export default userReducer;
