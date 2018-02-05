import { combineReducers } from 'redux';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  submittedIssues: issuesReducer,
  router: routerReducer,
  form: formReducer,
});

export default appReducer;
