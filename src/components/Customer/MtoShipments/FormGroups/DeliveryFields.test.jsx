/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';

import { DeliveryFields } from './DeliveryFields';

const defaultProps = {
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
};

function mountDeliveryDetails(props) {
  return mount(
    <Formik>
      <DeliveryFields {...defaultProps} {...props} />
    </Formik>,
  );
}

describe('DeliveryFields component', () => {
  it('renders expected child components', () => {
    const wrapper = mountDeliveryDetails();
    expect(wrapper.find('DatePickerInput').length).toBe(1);
    expect(wrapper.find('ContactInfoFields').length).toBe(1);
    expect(wrapper.find('AddressFields').length).toBe(0);
  });

  it('renders second address field when has delivery address', () => {
    const wrapper = mountDeliveryDetails({ hasDeliveryAddress: true });
    expect(wrapper.find('AddressFields').length).toBe(1);
  });
});
