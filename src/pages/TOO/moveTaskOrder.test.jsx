import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import { history, store } from '../../shared/store';

import MoveTaskOrder from './moveTaskOrder';

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
    expect(wrapper.find({ 'data-cy': 'too-shipment-container' }).exists()).toBe(true);
  });
});
