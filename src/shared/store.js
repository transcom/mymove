import { createStore, applyMiddleware } from 'redux';
import { appReducer, tspAppReducer, adminAppReducer } from 'appReducer';
import { createBrowserHistory } from 'history';
import { routerMiddleware } from 'react-router-redux';
import thunk from 'redux-thunk';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment, isTspSite, isAdminSite } from 'shared/constants';
import logger from './reduxLogger';
import * as schema from 'shared/Entities/schema';

export const history = createBrowserHistory();

const middlewares = [thunk.withExtraArgument({ schema }), routerMiddleware(history)];

if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
  middlewares.push(logger);
}

const composeEnhancers = composeWithDevTools({});

function appSelector() {
  if (isTspSite) {
    return tspAppReducer;
  } else if (isAdminSite) {
    return adminAppReducer;
  } else {
    return appReducer;
  }
}

export const store = composeEnhancers(applyMiddleware(...middlewares))(createStore)(appSelector());

export default store;
