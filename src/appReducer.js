import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import defaultMessages from 'ra-language-english';
import storage from 'redux-persist/lib/storage';
import { persistReducer } from 'redux-persist';

import authReducer from 'store/auth/reducer';
import onboardingReducer from 'store/onboarding/reducer';
import flashReducer from 'store/flash/reducer';
import interceptorReducer from 'store/interceptor/reducer';
import { swaggerReducerPublic, swaggerReducerInternal } from 'shared/Swagger/ducks';
import { requestsReducer } from 'shared/Swagger/requestsReducer';
import { entitiesReducer } from 'shared/Entities/reducer';
import { officeFlashMessagesReducer } from 'scenes/Office/ducks';
import officePpmReducer from 'scenes/Office/Ppm/ducks';

const locale = 'en';
const i18nProvider = () => defaultMessages;

const authPersistConfig = {
  key: 'auth',
  storage,
  whitelist: ['activeRole'],
};

const defaultReducers = {
  auth: persistReducer(authPersistConfig, authReducer),
  flash: flashReducer,
  form: formReducer,
  swaggerPublic: swaggerReducerPublic,
  requests: requestsReducer,
  entities: entitiesReducer,
};

export const appReducer = () =>
  combineReducers({
    ...defaultReducers,
    onboarding: onboardingReducer,
    swaggerInternal: swaggerReducerInternal,
    flashMessages: officeFlashMessagesReducer,
    interceptor: interceptorReducer,
    ppmIncentive: officePpmReducer,
  });

export const adminAppReducer = () =>
  combineReducers({
    ...defaultReducers,
    i18n: i18nProvider(locale),
    form: formReducer,
  });

export default appReducer;
