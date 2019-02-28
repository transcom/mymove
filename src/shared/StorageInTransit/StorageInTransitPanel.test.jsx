import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';

import { StorageInTransitPanel } from './StorageInTransitPanel';

import * as CONSTANTS from 'shared/constants.js';

const mockStore = configureStore();
let store;

describe('StorageInTransit tests', () => {
  describe('When no items exist', () => {
    let wrapper;
    const sitRequests = [];

    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitPanel sitRequests={sitRequests} shipmentId="" sitEntitlement={90} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-panel').length).toEqual(1);
    });
  });
  describe('When no items exists and Request SIT appears on TSP app', () => {
    CONSTANTS.isTspSite = true;
    let wrapper;
    const sitRequests = [];

    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitPanel sitRequests={sitRequests} shipmentId="" sitEntitlement={90} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-panel').length).toEqual(1);
      expect(wrapper.find('.add-request').length).toEqual(1);
    });
  });
  describe('When no items exists and Request SIT does not appears on Office app', () => {
    CONSTANTS.isTspSite = false;
    let wrapper;
    const sitRequests = [];

    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitPanel sitRequests={sitRequests} shipmentId="" sitEntitlement={90} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-panel').length).toEqual(1);
      expect(wrapper.find('.add-request').length).toEqual(0);
    });
  });
});
