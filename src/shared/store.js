import { createStore, applyMiddleware } from 'redux';
import { persistStore, persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import { appReducer, adminAppReducer } from 'appReducer';
import thunk from 'redux-thunk';
import createSagaMiddleware from 'redux-saga';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment, isAdminSite, isMilmoveSite } from 'shared/constants';
import logger from './reduxLogger';
import * as schema from 'shared/Entities/schema';

import rootSaga, { rootCustomerSaga } from 'sagas/index';
import { interceptorInjectionMiddleware } from 'store/interceptor/injectionMiddleware';

function appSelector() {
  if (isAdminSite) {
    return adminAppReducer();
  } else {
    return appReducer();
  }
}

export const configureStore = (initialState = {}) => {
  const sagaMiddleware = createSagaMiddleware();
  const middlewares = [thunk.withExtraArgument({ schema }), sagaMiddleware, interceptorInjectionMiddleware];

  if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
    middlewares.push(logger);
  }

  const composeEnhancers = composeWithDevTools({});

  const persistConfig = {
    key: 'root',
    storage,
    blacklist: ['auth'],
    whitelist: [],
  };

  const rootReducer = appSelector();
  const persistedReducer = persistReducer(persistConfig, rootReducer);

  const store = createStore(persistedReducer, initialState, composeEnhancers(applyMiddleware(...middlewares)));
  const persistor = persistStore(store);

  if (isMilmoveSite) {
    sagaMiddleware.run(rootCustomerSaga);
  } else {
    sagaMiddleware.run(rootSaga);
  }

  return { store, persistor };
};

export const { store, persistor } = configureStore();

export default store;
