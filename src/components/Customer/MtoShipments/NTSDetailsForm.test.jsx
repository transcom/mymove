/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import NTSDetailsForm from './NTSDetailsForm';

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
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
};
function mountNTSDetailsForm(props = defaultProps) {
  return mount(
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <NTSDetailsForm {...props} />
      </ConnectedRouter>
    </Provider>,
  );
}
describe('NTSDetailsForm component', () => {
  it('renders expected components', () => {
    const wrapper = mountNTSDetailsForm();
    expect(wrapper.find('NTSDetailsForm').length).toBe(1);
    expect(wrapper.find('PickupFields').length).toBe(1);
    expect(wrapper.find('input[name="customerRemarks"]').length).toBe(1);
  });
});
