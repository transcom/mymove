import { createStore, applyMiddleware } from 'redux';
import { persistStore, persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import { appReducer, adminAppReducer } from 'appReducer';
import { createBrowserHistory } from 'history';
import { routerMiddleware } from 'connected-react-router';
import thunk from 'redux-thunk';
import createSagaMiddleware from 'redux-saga';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment, isAdminSite, isMilmoveSite } from 'shared/constants';
import logger from './reduxLogger';
import * as schema from 'shared/Entities/schema';

import rootSaga, { rootCustomerSaga } from 'sagas/index';

export const history = createBrowserHistory();

function appSelector() {
  if (isAdminSite) {
    return adminAppReducer(history);
  } else {
    return appReducer(history);
  }
}

export const configureStore = (history, initialState = {}) => {
  const sagaMiddleware = createSagaMiddleware();
  const middlewares = [thunk.withExtraArgument({ schema }), routerMiddleware(history), sagaMiddleware];

  if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
    middlewares.push(logger);
  }

  const composeEnhancers = composeWithDevTools({});

  const persistConfig = {
    key: 'root',
    storage,
    whitelist: ['auth'],
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

export const { store, persistor } = configureStore(history);

export default store;
