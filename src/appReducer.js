import { combineReducers } from 'redux';
import { showIssues } from 'reducers/index';

import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  showIssues,
  router: routerReducer,
});

export default appReducer;
