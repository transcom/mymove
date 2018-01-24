import React, { Component } from 'react';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware, compose } from 'redux';
import thunk from 'redux-thunk';

import { appReducer } from './appReducer';
import AppWrapper from 'shared/App/AppWrapper';
import './App.css';

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
const store = createStore(appReducer, composeEnhancers(applyMiddleware(thunk)));

class App extends Component {
  render() {
    return (
      <Provider store={store}>
        <AppWrapper />
      </Provider>
    );
  }
}

export default App;
