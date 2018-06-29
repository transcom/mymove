import { createStore, applyMiddleware } from 'redux';
import appReducer from 'appReducer';
import createHistory from 'history/createBrowserHistory';
import { routerMiddleware } from 'react-router-redux';
import thunk from 'redux-thunk';

import { composeWithDevTools } from 'redux-devtools-extension/logOnlyInProduction';

import { isDevelopment } from 'shared/constants';
import logger from './reduxLogger';

export const history = createHistory();

const middlewares = [thunk, routerMiddleware(history)];

if (isDevelopment && !window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__) {
  middlewares.push(logger);
}

const composeEnhancers = composeWithDevTools({});

const store = composeEnhancers(applyMiddleware(...middlewares))(createStore)(
  appReducer,
);

if (isDevelopment) {
  //making state accessible to cypress
  const onStateChange = function(global) {
    global._state = this.getState();
  }.bind(store, window);
  store.subscribe(onStateChange);
  onStateChange();
}
export default store;
