import { createStore, applyMiddleware } from 'redux';
import { appReducer, adminAppReducer } from 'appReducer';
import { createBrowserHistory } from 'history';
import { routerMiddleware } from 'connected-react-router';
import thunk from 'redux-thunk';
import createSagaMiddleware from 'redux-saga';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment, isAdminSite } from 'shared/constants';
import logger from './reduxLogger';
import * as schema from 'shared/Entities/schema';

import rootSaga from 'sagas/index';

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
  const rootReducer = appSelector();
  const store = createStore(rootReducer, initialState, composeEnhancers(applyMiddleware(...middlewares)));

  sagaMiddleware.run(rootSaga);

  return store;
};

export const store = configureStore(history);

export default store;
