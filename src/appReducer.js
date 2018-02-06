import { combineReducers } from 'redux';
import { feedbackReducer } from 'scenes/Feedback/ducks';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';
import {
  availableShipmentsReducer,
  awardedShipmentsReducer,
} from 'scenes/Shipments/ducks';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

export const appReducer = combineReducers({
  submittedIssues: issuesReducer,
  availableShipments: availableShipmentsReducer,
  awardedShipments: awardedShipmentsReducer,
  router: routerReducer,
  form: formReducer,
  feedback: feedbackReducer,
});

export default appReducer;
