import { combineReducers } from 'redux';
import { feedbackReducer } from 'scenes/Feedback/ducks';
import { default as userReducer } from 'shared/User/ducks';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';
import { shipmentsReducer } from 'scenes/Shipments/ducks';
import dd1299Reducer from 'scenes/DD1299/ducks';
import { signedCertificationReducer } from 'scenes/Legalese/ducks';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  user: userReducer,
  submittedIssues: issuesReducer,
  shipments: shipmentsReducer,
  router: routerReducer,
  form: formReducer,
  feedback: feedbackReducer,
  signedCertification: signedCertificationReducer,
  DD1299: dd1299Reducer,
});

export default appReducer;
