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
import { reviewReducer } from 'scenes/Review/ducks';
import transportationOfficeReducer from 'shared/TransportationOffices/ducks';
import { officeReducer } from 'scenes/Office/ducks';
import { tspReducer } from 'scenes/TransportationServiceProvider/ducks';
import { documentsReducer } from 'scenes/Office/DocumentViewer/ducks';
import officePpmReducer from 'scenes/Office/Ppm/ducks';
import entititesReducer from 'shared/Entities/ducks';

export const appReducer = combineReducers({
  user: userReducer,
  loggedInUser: loggedInUserReducer,
  swagger: swaggerReducer,
  submittedIssues: issuesReducer,
  moves: moveReducer,
  moveDocuments: documentsReducer,
  ppm: ppmReducer,
  serviceMember: serviceMemberReducer,
  orders: ordersReducer,
  shipments: shipmentsReducer,
  router: routerReducer,
  form: formReducer,
  feedback: feedbackReducer,
  signedCertification: signedCertificationReducer,
  upload: documentReducer,
  review: reviewReducer,
  office: officeReducer,
  tsp: tspReducer,
  transportationOffices: transportationOfficeReducer,
  ppmIncentive: officePpmReducer,
  entities: entititesReducer,
});

export default appReducer;
