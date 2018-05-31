import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

import { loggedInUserReducer } from 'shared/User/ducks';
import userReducer from 'shared/User/ducks';
import swaggerReducer from 'shared/Swagger/ducks';

import { feedbackReducer } from 'scenes/Feedback/ducks';
import { moveReducer } from 'scenes/Moves/ducks';
import { ppmReducer } from 'scenes/Moves/Ppm/ducks';
import { serviceMemberReducer } from 'scenes/ServiceMembers/ducks';
import { ordersReducer } from 'scenes/Orders/ducks';
import issuesReducer from 'scenes/SubmittedFeedback/ducks';
import { shipmentsReducer } from 'scenes/Shipments/ducks';
import { signedCertificationReducer } from 'scenes/Legalese/ducks';
import { documentReducer } from 'shared/Uploader/ducks';
import transportationOfficeReducer from 'shared/TransportationOffices/ducks';
import { officeReducer } from 'scenes/Office/ducks';

export const appReducer = combineReducers({
  user: userReducer,
  loggedInUser: loggedInUserReducer,
  swagger: swaggerReducer,
  submittedIssues: issuesReducer,
  moves: moveReducer,
  ppm: ppmReducer,
  serviceMember: serviceMemberReducer,
  orders: ordersReducer,
  shipments: shipmentsReducer,
  router: routerReducer,
  form: formReducer,
  feedback: feedbackReducer,
  signedCertification: signedCertificationReducer,
  upload: documentReducer,
  office: officeReducer,
  transportationOffices: transportationOfficeReducer,
});

export default appReducer;
