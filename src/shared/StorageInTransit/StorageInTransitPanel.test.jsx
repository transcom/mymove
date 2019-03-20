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

  describe('When some items exist', () => {
    const sitRequests = [
      {
        estimated_start_date: '2019-02-12',
        id: '5cd370a1-ac3d-4fb3-86a3-c4f23e289687',
        location: 'ORIGIN',
        shipment_id: 'dd67cec5-334a-4209-a9d9-a14485414052',
        status: 'REQUESTED',
        warehouse_address: {
          city: 'Beverly Hills',
          postal_code: '90210',
          state: 'CA',
          street_address_1: '123 Any Street',
        },
        warehouse_id: '76567867',
        warehouse_name: 'haus',
      },
      {
        estimated_start_date: '2019-02-12',
        id: '5cd370a1-ac3d-4fb3-86a3-c4f23e289689',
        location: 'DESTINATION',
        notes: 'notes',
        shipment_id: 'dd67cec5-334a-4209-a9d9-a14485414052',
        status: 'REQUESTED',
        warehouse_address: {
          city: 'Beverly Hills',
          postal_code: '90210',
          state: 'CA',
          street_address_1: '123 Any Street',
        },
        warehouse_id: '76567869',
        warehouse_name: 'hausen',
      },
    ];

    let store = mockStore({});
    let wrapper = mount(
      <Provider store={store}>
        <StorageInTransitPanel
          storageInTransits={sitRequests}
          shipmentId="dd67cec5-334a-4209-a9d9-a14485414052"
          sitEntitlement={90}
        />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-panel').length).toEqual(1);
    });

    it('renders the first mocked Storage In Transit request', () => {
      expect(wrapper.find('.column-head').get(1).props.children).toContain('Origin');
      expect(wrapper.find('.column-head').get(1).props.children[3].props.children[1]).toEqual('Requested');
      expect(wrapper.find('.column-subhead').get(0).props.children).toEqual('Dates');
      expect(wrapper.find('.column-subhead').get(1).props.children).toEqual('Warehouse');
    });

    it('renders the second mocked Storage In Transit request', () => {
      expect(wrapper.find('.column-head').get(2).props.children).toContain('Destination');
      expect(wrapper.find('.column-head').get(2).props.children[3].props.children[1]).toEqual('Requested');

      expect(wrapper.find('.column-subhead').get(3).props.children).toEqual('Note');
      expect(wrapper.find('.column-subhead').get(0).props.children).toEqual('Dates');
      expect(wrapper.find('.column-subhead').get(1).props.children).toEqual('Warehouse');
    });
  });
});
