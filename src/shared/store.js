import { createStore, applyMiddleware } from 'redux';
import { appReducer, adminAppReducer } from 'appReducer';
import { createBrowserHistory } from 'history';
import { routerMiddleware } from 'connected-react-router';
import thunk from 'redux-thunk';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment, isAdminSite } from 'shared/constants';
import logger from './reduxLogger';
import * as schema from 'shared/Entities/schema';

export const history = createBrowserHistory();

function appSelector() {
  if (isAdminSite) {
    return adminAppReducer(history);
  } else {
    return appReducer(history);
  }
}

export const configureStore = (history, initialState = {}) => {
  const middlewares = [thunk.withExtraArgument({ schema }), routerMiddleware(history)];

  if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
    middlewares.push(logger);
  }

  const composeEnhancers = composeWithDevTools({});
  const rootReducer = appSelector();
  const store = createStore(rootReducer, initialState, composeEnhancers(applyMiddleware(...middlewares)));
  return store;
};

export const store = configureStore(history);

export default store;
