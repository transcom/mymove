import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import SMName from './SMName.jsx';
import store from 'shared/store';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <SMName />
    </Provider>,
    div,
  );
});
