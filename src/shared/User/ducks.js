import * as Cookies from 'js-cookie';
import * as decode from 'jwt-decode';

const LOAD_USER_AND_TOKEN = 'USER|LOAD_USER_AND_TOKEN';

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
