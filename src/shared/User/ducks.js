import * as Cookies from 'js-cookie';
const LOGOUT = 'USER|LOGOUT';

export function logOut() {
  //todo: call server endpoint, clean up cookies?
  return {
    type: LOGOUT,
  };
}
const initialState = {
  isLoggedIn: Cookies.get('user_session') ? true : false,
};
const userReducer = (state = initialState, action) => {
  switch (action.type) {
    case LOGOUT:
      return { isLoggedIn: false };
    default: {
      return state;
    }
  }
};

export default userReducer;
