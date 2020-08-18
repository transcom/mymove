/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import EditShipment from './EditShipment';

function mountEditShipment(props) {
  return mount(<EditShipment {...props} />);
}
describe('EditShipment component', () => {
  it('renders expected form components', () => {
    const wrapper = mountEditShipment();
    expect(wrapper.find('EditShipment').length).toBe(1);
    expect(wrapper.find('DatePickerInput').length).toBe(2);
    expect(wrapper.find('AddressFields').length).toBe(1);
    expect(wrapper.find('ContactInfoFields').length).toBe(2);
  });

  it('renders second address field when has delivery address', () => {
    const wrapper = mountEditShipment();
    wrapper.setState({ hasDeliveryAddress: true });
    expect(wrapper.find('AddressFields').length).toBe(2);
  });
});
