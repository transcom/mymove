import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { StorageInTransitOfficeDenyForm } from './StorageInTransitOfficeDenyForm';

let store;
const mockStore = configureStore();
const submit = jest.fn();

const storageInTransitSchema = {
  properties: {
    notes: {
      type: 'string',
      format: 'textarea',
      example: 'this is a note',
      title: 'Reason for denial',
    },
  },
};

describe('StorageInTransitOfficeDenyForm tests', () => {
  describe('Empty form', () => {
    let wrapper;
    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitOfficeDenyForm onSubmit={submit} storageInTransitSchema={storageInTransitSchema} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-office-deny-form').length).toEqual(1);
    });
  });
});
