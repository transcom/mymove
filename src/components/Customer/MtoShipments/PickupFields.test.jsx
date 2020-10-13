/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';

import { PickupFields } from './PickupFields';

const defaultProps = {
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
};

function mountPickupDetails(props) {
  return mount(
    <Formik>
      <PickupFields {...defaultProps} {...props} />
    </Formik>,
  );
}

describe('PickupFields component', () => {
  it('renders expected child components', () => {
    const wrapper = mountPickupDetails();
    expect(wrapper.find('DatePickerInput').length).toBe(1);
    expect(wrapper.find('ContactInfoFields').length).toBe(1);
    expect(wrapper.find('AddressFields').length).toBe(1);
  });
});
