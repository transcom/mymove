import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import MoveTaskOrder from './moveTaskOrder';

import { history, store } from 'shared/store';

describe('MoveTaskOrder', () => {
  const wrapper = mount(
    <Provider store={store}>
      clear
      <ConnectedRouter history={history}>
        <MoveTaskOrder />
      </ConnectedRouter>
    </Provider>,
  );

  it('should render the h1', () => {
    expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
  });
});
