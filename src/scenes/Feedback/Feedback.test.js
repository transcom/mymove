import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { shallow } from 'enzyme';
import Feedback from './Feedback';
import store from 'shared/store';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(
    <Provider store={store}>
      <Feedback />
    </Provider>,
    div,
  );
});
