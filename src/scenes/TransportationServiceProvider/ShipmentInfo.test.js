import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { ReferrerQueueLink } from './ShipmentInfo';
import MockRouter from 'react-mock-router';

const mockstore = configureStore();
let wrapper;
let store;
describe('ShipmentInfo tests', () => {
  describe('Shows correct queue to return to', () => {
    beforeEach(() => {
      store = mockstore({});
    });
    it('when a referrer is set in history', () => {
      wrapper = mount(
        <Provider store={store}>
          <MockRouter push={jest.fn()}>
            <ReferrerQueueLink history={{ location: { state: { referrerPathname: '/queues/accepted' } } }} />
          </MockRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('Accepted Shipments Queue');
    });
    it('when no referrer is set', () => {
      wrapper = mount(
        <Provider store={store}>
          <MockRouter push={jest.fn()}>
            <ReferrerQueueLink history={{ location: {} }} />
          </MockRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('New Shipments Queue');
    });
  });
});
