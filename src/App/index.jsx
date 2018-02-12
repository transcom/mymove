import React, { Component } from 'react';
import { Provider } from 'react-redux';

import AppWrapper from 'shared/AppWrapper';
import store from 'shared/store';
import './index.css';

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
