import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import { history, store } from '../../shared/store';
import Accessorials from './accessorials';

describe('Accessorials', () => {
  const wrapper = mount(
    <Provider store={store}>
      clear
      <ConnectedRouter history={history}>
        <Accessorials />
      </ConnectedRouter>
    </Provider>,
  );

  it('should render the h1', () => {
    expect(wrapper.contains(<h1>This is where we will put our accessorial components!</h1>)).toBe(true);
  });
});
