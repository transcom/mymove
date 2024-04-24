import { combineReducers } from 'redux';
import { reducer as formReducer } from 'redux-form';
import storage from 'redux-persist/lib/storage';
import { persistReducer } from 'redux-persist';

import authReducer from 'store/auth/reducer';
import onboardingReducer from 'store/onboarding/reducer';
import flashReducer from 'store/flash/reducer';
import interceptorReducer from 'store/interceptor/reducer';
import generalStateReducer from 'store/general/reducer';
import { swaggerReducerPublic, swaggerReducerInternal } from 'shared/Swagger/ducks';
import { requestsReducer } from 'shared/Swagger/requestsReducer';
import { entitiesReducer } from 'shared/Entities/reducer';

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
    interceptor: interceptorReducer,
    generalState: generalStateReducer,
  });

export default appReducer;
