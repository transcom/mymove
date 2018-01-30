import { combineReducers } from 'redux';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';

import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  issues: issuesReducer,
  router: routerReducer,
});

export default appReducer;
