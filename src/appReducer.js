import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

import { loggedInUserReducer } from 'shared/User/ducks';
import userReducer from 'shared/User/ducks';
import swaggerReducer from 'shared/Swagger/ducks';
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

// Entities
import moveDocuments, {
  STATE_KEY as MOVEDOCUMENTS_STATE_KEY,
} from 'shared/Entities/modules/moveDocuments';
import documentModel, {
  STATE_KEY as DOCUMENTS_STATE_KEY,
} from 'shared/Entities/modules/documents';
import uploads, {
  STATE_KEY as UPLOADS_STATE_KEY,
} from 'shared/Entities/modules/uploads';
import shipments, {
  STATE_KEY as SHIPMENTS_STATE_KEY,
} from 'shared/Entities/modules/shipments';
import addresses, {
  STATE_KEY as ADDRESSES_STATE_KEY,
} from 'shared/Entities/modules/addresses';
import moves, {
  STATE_KEY as MOVES_STATE_KEY,
} from 'shared/Entities/modules/moves';
import orders, {
  STATE_KEY as ORDERS_STATE_KEY,
} from 'shared/Entities/modules/orders';

const entititesReducer = combineReducers({
  [MOVEDOCUMENTS_STATE_KEY]: moveDocuments,
  [DOCUMENTS_STATE_KEY]: documentModel,
  [UPLOADS_STATE_KEY]: uploads,
  [SHIPMENTS_STATE_KEY]: shipments,
  [ADDRESSES_STATE_KEY]: addresses,
  [MOVES_STATE_KEY]: moves,
  [ORDERS_STATE_KEY]: orders,
});

const defaultReducers = {
  loggedInUser: loggedInUserReducer,
  router: routerReducer,
  swagger: swaggerReducer,
  ui: uiReducer,
  user: userReducer,
};

export const appReducer = combineReducers(
  Object.assign({}, defaultReducers, {
    submittedIssues: issuesReducer,
    moves: moveReducer,
    ppm: ppmReducer,
    serviceMember: serviceMemberReducer,
    orders: ordersReducer,
    shipments: shipmentsReducer,
    form: formReducer,
    feedback: feedbackReducer,
    signedCertification: signedCertificationReducer,
    upload: documentReducer,
    review: reviewReducer,
    office: officeReducer,
    transportationOffices: transportationOfficeReducer,
    ppmIncentive: officePpmReducer,
    entities: entititesReducer,
  }),
);

export const tspAppReducer = combineReducers(
  Object.assign({}, defaultReducers, {
    tsp: tspReducer,
  }),
);

export default appReducer;
