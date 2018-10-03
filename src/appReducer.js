import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

import { loggedInUserReducer } from 'shared/User/ducks';
import userReducer from 'shared/User/ducks';
import { swaggerReducer } from 'shared/Swagger/ducks';
import { requestsReducer } from 'shared/Swagger/requestsReducer';
import { entitiesReducer } from 'shared/Entities/reducer';
import uiReducer from 'shared/UI/ducks';

import { feedbackReducer } from 'scenes/Feedback/ducks';
import { moveReducer } from 'scenes/Moves/ducks';
import { ppmReducer } from 'scenes/Moves/Ppm/ducks';
import { serviceMemberReducer } from 'scenes/ServiceMembers/ducks';
import { ordersReducer } from 'scenes/Orders/ducks';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';
import { shipmentsReducer } from 'scenes/Shipments/ducks';
import { signedCertificationReducer } from 'scenes/Legalese/ducks';
import { documentReducer } from 'shared/Uploader/ducks';
import { reviewReducer } from 'scenes/Review/ducks';
import transportationOfficeReducer from 'shared/TransportationOffices/ducks';
import { officeReducer } from 'scenes/Office/ducks';
import { tspReducer } from 'scenes/TransportationServiceProvider/ducks';
import officePpmReducer from 'scenes/Office/Ppm/ducks';

const defaultReducers = {
  form: formReducer,
  loggedInUser: loggedInUserReducer,
  router: routerReducer,
  swagger: swaggerReducer,
  requests: requestsReducer,
  ui: uiReducer,
  user: userReducer,
  entities: entitiesReducer,
};

export const appReducer = combineReducers({
  ...defaultReducers,
  submittedIssues: issuesReducer,
  moves: moveReducer,
  ppm: ppmReducer,
  serviceMember: serviceMemberReducer,
  orders: ordersReducer,
  shipments: shipmentsReducer,
  feedback: feedbackReducer,
  signedCertification: signedCertificationReducer,
  upload: documentReducer,
  review: reviewReducer,
  office: officeReducer,
  transportationOffices: transportationOfficeReducer,
  ppmIncentive: officePpmReducer,
});

export const tspAppReducer = combineReducers({
  ...defaultReducers,
  tsp: tspReducer,
});

export default appReducer;
