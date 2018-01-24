import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  routing: routerReducer,
});

export default appReducer;
