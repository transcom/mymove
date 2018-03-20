import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import Transition from '.';
import store from 'shared/store';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <Transition match={{ match: 'match' }} />
    </Provider>,
    div,
  );
});
