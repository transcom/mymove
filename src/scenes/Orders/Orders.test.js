import store from 'shared/store';
import Orders from '../Orders/Orders';
import React from 'react';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router-dom';
import { mount } from 'enzyme';

describe('Orders page', () => {
  it('renders', () => {});
  const wrapper = mount(
    <Provider store={store}>
      <Router push={jest.fn()}>
        <Orders pages={[]} pageKey="" updateOrders={jest.fn()} />
      </Router>
    </Provider>,
  );
  expect(wrapper.length).toEqual(1);
});
