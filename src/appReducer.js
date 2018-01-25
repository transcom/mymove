import { combineReducers } from 'redux';

import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({ router: routerReducer });

export default appReducer;
