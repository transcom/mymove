import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import { history, store } from '../../shared/store';
import MoveDetails from './moveDetails';

describe('MoveDetails', () => {
  const wrapper = mount(
    <Provider store={store}>
      clear
      <ConnectedRouter history={history}>
        <MoveDetails />
      </ConnectedRouter>
    </Provider>,
  );

  it('renders the h1', () => {
    expect(wrapper.find({ 'data-cy': 'too-move-details' }).exists()).toBe(true);
  });
});
