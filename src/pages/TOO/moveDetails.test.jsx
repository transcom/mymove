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
  it('renders the Orders Table', () => {
    expect(wrapper.find({ 'data-cy': 'currentDutyStation' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'newDutyStation' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'issuedDate' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'reportByDate' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'departmentIndicator' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'ordersNumber' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'ordersType' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'ordersTypeDetail' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'tacMDC' }).exists()).toBe(true);
    expect(wrapper.find({ 'data-cy': 'sacSDN' }).exists()).toBe(true);
  });
});
