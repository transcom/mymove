import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { StorageInTransitOfficeEditForm } from './StorageInTransitOfficeEditForm';

let store;
const mockStore = configureStore();
const submit = jest.fn();

const storageInTransitSchema = {
  properties: {
    location: {
      type: 'string',
      title: 'SIT Location',
    },
    estimated_start_date: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Estimated start date',
    },
    warehouse_id: {
      type: 'string',
      example: '000383',
      title: 'Warehouse Name',
    },
    warehouse_name: {
      type: 'string',
      example: 'ABC Warehouse, Inc.',
      title: 'Warehouse Name',
    },
    notes: {
      type: 'string',
      example: 'Good to go!',
      format: 'textarea',
      title: 'Note',
    },
  },
};

const storageInTransit = {
  estimated_start_date: '2019-02-12',
  id: '5cd370a1-ac3d-4fb3-86a3-c4f23e289687',
  shipment_id: 'dd67cec5-334a-4209-a9d9-a14485414052',
  status: 'APPROVED',
  warehouse_address: {
    city: 'Beverly Hills',
    postal_code: '90210',
    state: 'CA',
    street_address_1: '123 Any Street',
  },
  warehouse_id: '76567867',
  warehouse_name: 'haus',
};

describe('StorageInTransitOfficeEditForm tests', () => {
  describe('Empty form', () => {
    let wrapper;
    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitOfficeEditForm
          onSubmit={submit}
          storageInTransitSchema={storageInTransitSchema}
          initialValues={storageInTransit}
        />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-form').length).toEqual(1);
    });
  });
});
