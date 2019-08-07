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

const middlewares = [thunk.withExtraArgument({ schema }), routerMiddleware(history)];

if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
  middlewares.push(logger);
}

const composeEnhancers = composeWithDevTools({});

function appSelector(history) {
  if (isAdminSite) {
    return adminAppReducer(history);
  } else {
    return appReducer(history);
  }
}

export const store = composeEnhancers(applyMiddleware(...middlewares))(createStore)(appSelector(history));

export default store;
