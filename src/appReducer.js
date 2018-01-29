import { combineReducers } from 'redux';
import issues from 'scenes/SubmittedFeedback/ducks';

import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  issues,
  router: routerReducer,
});

export default appReducer;
