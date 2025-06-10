import { combineReducers } from 'redux';
import storage from 'redux-persist/lib/storage';
import { persistReducer } from 'redux-persist';

import authReducer from 'store/auth/reducer';
import onboardingReducer from 'store/onboarding/reducer';
import flashReducer from 'store/flash/reducer';
import interceptorReducer from 'store/interceptor/reducer';
import generalStateReducer from 'store/general/reducer';
import { swaggerReducerPublic, swaggerReducerInternal } from 'shared/Swagger/ducks';
import { requestsReducer } from 'shared/Swagger/requestsReducer';
import { adminEntitiesReducer, entitiesReducer } from 'shared/Entities/reducer';

const authPersistConfig = {
  key: 'auth',
  storage,
  whitelist: ['activeRole', 'activeOffice'],
};

const defaultReducers = {
  auth: persistReducer(authPersistConfig, authReducer),
  flash: flashReducer,
  swaggerPublic: swaggerReducerPublic,
  requests: requestsReducer,
  entities: entitiesReducer,
};

export const appReducer = () =>
  combineReducers({
    ...defaultReducers,
    onboarding: onboardingReducer,
    swaggerInternal: swaggerReducerInternal,
    interceptor: interceptorReducer,
    generalState: generalStateReducer,
  });

export const adminAppReducer = () =>
  combineReducers({
    ...defaultReducers,
    entities: adminEntitiesReducer,
  });

export default appReducer;
