/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import { history, store } from '../../shared/store';

import HHGDetailsForm, { HHGDetailsFormComponent } from './HHGDetailsForm';

const defaultProps = {
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  match: { isExact: false, path: '', url: '', params: { moveId: '' } },
  showLoggedInUser: () => {},
  createMTOShipment: () => {},
  updateMTOShipment: () => {},
  push: () => {},
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
};
function mountHHGDetailsForm(props = defaultProps) {
  return mount(
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <HHGDetailsForm {...props} />
      </ConnectedRouter>
    </Provider>,
  );
}
describe('HHGDetailsForm component', () => {
  it('renders expected form components', () => {
    const wrapper = mountHHGDetailsForm();
    expect(wrapper.find('HHGDetailsForm').length).toBe(1);
    expect(wrapper.find('DatePickerInput').length).toBe(2);
    expect(wrapper.find('AddressFields').length).toBe(1);
    expect(wrapper.find('ContactInfoFields').length).toBe(2);
    expect(wrapper.find('input[name="customerRemarks"]').length).toBe(1);
  });

  it('renders second address field when has delivery address', () => {
    const wrapper = mount(<HHGDetailsFormComponent {...defaultProps} />);
    wrapper.setState({ hasDeliveryAddress: true });
    expect(wrapper.find('AddressFields').length).toBe(2);
  });
});
