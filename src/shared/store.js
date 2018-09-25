import { createStore, applyMiddleware } from 'redux';
import { appReducer, tspAppReducer } from 'appReducer';
import createHistory from 'history/createBrowserHistory';
import { routerMiddleware } from 'react-router-redux';
import thunk from 'redux-thunk';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment } from 'shared/constants';
import logger from './reduxLogger';
import * as schema from 'shared/Entities/schema';
import { getClient } from 'shared/Swagger/api';

export const history = createHistory();

const middlewares = [
  thunk.withExtraArgument({ schema, client: getClient() }),
  routerMiddleware(history),
];

if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
  middlewares.push(logger);
}

const composeEnhancers = composeWithDevTools({ name: 'sm-office' });

export const store = composeEnhancers(applyMiddleware(...middlewares))(
  createStore,
)(appReducer);

const tspComposeEnhancers = composeWithDevTools({ name: 'tsp' });

export const tspStore = tspComposeEnhancers(applyMiddleware(...middlewares))(
  createStore,
)(tspAppReducer);

export default store;
