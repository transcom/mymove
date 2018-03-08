import * as Cookies from 'js-cookie';
import * as decode from 'jwt-decode';

const LOGOUT = 'USER|LOGOUT';

const LOAD_USER_AND_TOKEN = 'USER|ZLOAD_USER_AND_TOKEN';

const loggedOutUser = {
  isLoggedIn: false,
  email: null,
  jwt: null,
};
function getUserInfo() {
  const cookie = Cookies.get('user_session');
  if (!cookie) return loggedOutUser;
  const jwt = decode(cookie);
  //if (jwt.exp <  Date.now().valueOf() / 1000) return loggedOutUser;
  return {
    jwt: cookie,
    email: jwt.email,
    isLoggedIn: true,
  };
}

export function loadUserAndToken() {
  const jwt = getUserInfo();
  return { type: LOAD_USER_AND_TOKEN, payload: jwt };
}

export function logOut() {
  //todo: call server endpoint, clean up cookies?
  return {
    type: LOGOUT,
  };
}

const userReducer = (state = loggedOutUser, action) => {
  switch (action.type) {
    case LOAD_USER_AND_TOKEN:
      return action.payload;
    case LOGOUT:
      return loggedOutUser;
    default: {
      return state;
    }
  }
};

export default userReducer;
