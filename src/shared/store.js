import { createStore, applyMiddleware } from 'redux';
import { appReducer, tspAppReducer } from 'appReducer';
import createHistory from 'history/createBrowserHistory';
import { routerMiddleware } from 'react-router-redux';
import thunk from 'redux-thunk';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment } from 'shared/constants';
import logger from './reduxLogger';
import * as schema from 'shared/Entities/schema';

export const history = createHistory();

const middlewares = [
  thunk.withExtraArgument({ schema }),
  routerMiddleware(history),
];

if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
  middlewares.push(logger);
}

const composeEnhancers = composeWithDevTools({});

export const store = composeEnhancers(applyMiddleware(...middlewares))(createStore)(
  appReducer,
);

export const tspStore = composeEnhancers(applyMiddleware(...middlewares))(createStore)(
  tspAppReducer,
);

export default store;
