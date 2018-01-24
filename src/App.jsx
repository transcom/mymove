import React, { Component } from 'react';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';

import AppWrapper from 'shared/App/AppWrapper';

import './App.css';

function issues(state = [], action) {
  return state;
}

const store = createStore(issues, {}, applyMiddleware(thunk));

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
