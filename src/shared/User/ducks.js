//import * as Cookies from 'js-cookie';
const LOGOUT = 'USER|LOGOUT';

export function logOut() {
  //todo: call server endpoint, clean up cookies?
  return {
    type: LOGOUT,
  };
}
const initialState = {
  //  loggedIn: Cookies.get('user_session'),
  loggedIn: false,
};
const userReducer = (state = initialState, action) => {
  switch (action.type) {
    case LOGOUT:
      return { loggedIn: false };
    default: {
      return state;
    }
  }
};

export default userReducer;
