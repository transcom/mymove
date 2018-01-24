import React, { Component } from 'react';
import { createStore, applyMiddleware } from 'redux';
import { Provider } from 'react-redux';
import { BrowserRouter as Router } from 'react-router-dom';
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
        <Router>
          <AppWrapper />
        </Router>
      </Provider>
    );
  }
}

export default App;
