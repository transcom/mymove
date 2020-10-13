/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import NTSrDetailsForm from './NTSrDetailsForm';

import { history, store } from 'shared/store';

const defaultProps = {
  wizardPage: {
    pageList: ['page1', 'anotherPage/:foo/:bar'],
    pageKey: 'page1',
    match: { isExact: false, path: '', url: '', params: { moveId: '' } },
    history: {
      goBack: jest.fn(),
      push: jest.fn(),
    },
  },
  showLoggedInUser: jest.fn(),
  createMTOShipment: jest.fn(),
  updateMTOShipment: jest.fn(),
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
  mtoShipment: {
    id: '',
    customerRemarks: '',
    requestedDeliveryDate: '',
    destinationAddress: {
      city: '',
      postal_code: '',
      state: '',
      street_address_1: '',
    },
  },
};

function mountNTSrDetailsForm(props = defaultProps) {
  return mount(
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <NTSrDetailsForm {...props} />
      </ConnectedRouter>
    </Provider>,
  );
}
describe('NTSrDetailsForm component', () => {
  it('renders expected child components', () => {
    const wrapper = mountNTSrDetailsForm();

    // should contain
    expect(wrapper.find('NTSrDetailsForm').length).toBe(1);
    expect(wrapper.find('DeliveryFields').length).toBe(1);
    expect(wrapper.find('input[name="customerRemarks"]').length).toBe(1);

    // should not contain
    expect(wrapper.find('PickupFields').length).toBe(0);
  });
});
