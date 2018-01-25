import { createStore, applyMiddleware, compose } from 'redux';
import appReducer from 'appReducer';
import createHistory from 'history/createBrowserHistory';
import { routerMiddleware } from 'react-router-redux';
import thunk from 'redux-thunk';

export const history = createHistory();
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
const store = createStore(
  appReducer,
  composeEnhancers(applyMiddleware(thunk, routerMiddleware(history))),
);

export default store;
