import { combineReducers } from 'redux';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';
import { shipmentsReducer } from 'scenes/Shipments/ducks';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  submittedIssues: issuesReducer,
  shipments: shipmentsReducer,
  router: routerReducer,
  form: formReducer,
});

export default appReducer;
