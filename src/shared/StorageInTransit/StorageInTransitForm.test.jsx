import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { StorageInTransitForm } from './StorageInTransitForm';

let store;
const mockStore = configureStore();
const submit = jest.fn();

const storageInTransitSchema = {
  properties: {
    location: {
      type: 'string',
      title: 'SIT Location',
    },
    estimate_start_date: {
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
  },
};

const addressSchema = {
  properties: {
    warehouse_address: {
      street_address_1: '123 Disney Rd',
      city: 'Los Angeles',
      state: 'CA',
      postal_code: '90210',
    },
  },
};

describe('StorageInTransit tests', () => {
  describe('Empty form', () => {
    let wrapper;
    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitForm
          onSubmit={submit}
          storageInTransitSchema={storageInTransitSchema}
          addressSchema={addressSchema}
        />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-request-form').length).toEqual(1);
    });
  });
});
