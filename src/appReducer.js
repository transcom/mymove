import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import { routerReducer } from 'react-router-redux';

import userReducer from 'shared/Data/users';
import { swaggerReducerPublic, swaggerReducerInternal } from 'shared/Swagger/ducks';
import { requestsReducer } from 'shared/Swagger/requestsReducer';
import { entitiesReducer } from 'shared/Entities/reducer';
import uiReducer from 'shared/UI/ducks';

import { moveReducer } from 'scenes/Moves/ducks';
import { ppmReducer } from 'scenes/Moves/Ppm/ducks';
import { serviceMemberReducer } from 'scenes/ServiceMembers/ducks';
import { ordersReducer } from 'scenes/Orders/ducks';
import { signedCertificationReducer } from 'scenes/Legalese/ducks';
import { documentReducer } from 'shared/Uploader/ducks';
import { reviewReducer } from 'scenes/Review/ducks';
import transportationOfficeReducer from 'shared/TransportationOffices/ducks';
import { officeFlashMessagesReducer } from 'scenes/Office/ducks';
import { tspReducer } from 'scenes/TransportationServiceProvider/ducks';
import officePpmReducer from 'scenes/Office/Ppm/ducks';
import { adminReducer, i18nReducer } from 'react-admin';
import defaultMessages from 'ra-language-english';

const locale = 'en';
const i18nProvider = () => defaultMessages;

const defaultReducers = {
  form: formReducer,
  router: routerReducer,
  swaggerPublic: swaggerReducerPublic,
  requests: requestsReducer,
  ui: uiReducer,
  user: userReducer,
  entities: entitiesReducer,
};

export const appReducer = combineReducers({
  ...defaultReducers,
  swaggerInternal: swaggerReducerInternal,
  moves: moveReducer,
  ppm: ppmReducer,
  serviceMember: serviceMemberReducer,
  orders: ordersReducer,
  signedCertification: signedCertificationReducer,
  upload: documentReducer,
  review: reviewReducer,
  flashMessages: officeFlashMessagesReducer,
  transportationOffices: transportationOfficeReducer,
  ppmIncentive: officePpmReducer,
});

export const tspAppReducer = combineReducers({
  ...defaultReducers,
  tsp: tspReducer,
});

export const adminAppReducer = combineReducers({
  ...defaultReducers,
  admin: adminReducer,
  i18n: i18nReducer(locale, i18nProvider(locale)),
  form: formReducer,
  router: routerReducer,
});

export default appReducer;
