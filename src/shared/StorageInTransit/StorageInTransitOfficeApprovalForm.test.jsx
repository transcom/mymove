import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import { StorageInTransitOfficeApprovalForm } from './StorageInTransitOfficeApprovalForm';

let store;
const mockStore = configureStore();
const submit = jest.fn();

const storageInTransitSchema = {
  properties: {
    estimate_start_date: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Estimated start date',
    },
  },
};

describe('StorageInTransitOfficeApprovalForm tests', () => {
  describe('Empty form', () => {
    let wrapper;
    store = mockStore({});
    wrapper = mount(
      <Provider store={store}>
        <StorageInTransitOfficeApprovalForm onSubmit={submit} storageInTransitSchema={storageInTransitSchema} />
      </Provider>,
    );

    it('renders without crashing', () => {
      expect(wrapper.find('.storage-in-transit-office-approval-form').length).toEqual(1);
    });
  });
});
