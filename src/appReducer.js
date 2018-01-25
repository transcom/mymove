import { combineReducers } from 'redux';
import { showIssues } from 'reducers/index';

export const appReducer = combineReducers({
  showIssues,
});

export default appReducer;
