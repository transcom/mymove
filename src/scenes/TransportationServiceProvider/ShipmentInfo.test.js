import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { ReferrerQueueLink } from './ShipmentInfo';
import { MemoryRouter } from 'react-router';

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
          <MemoryRouter>
            <ReferrerQueueLink history={{ location: { state: { referrerPathname: '/queues/accepted' } } }} />
          </MemoryRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('Accepted Shipments Queue');
    });
    it('when no referrer is set', () => {
      wrapper = mount(
        <Provider store={store}>
          <MemoryRouter>
            <ReferrerQueueLink history={{ location: {} }} />
          </MemoryRouter>
        </Provider>,
      );
      expect(wrapper.text()).toEqual('New Shipments Queue');
    });
  });
});
