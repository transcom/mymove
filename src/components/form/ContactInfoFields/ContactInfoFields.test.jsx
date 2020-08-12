import React from 'react';
import { mount } from 'enzyme';
import { Formik } from 'formik';

import { ContactInfoFields } from './ContactInfoFields';

describe('ContactInfoFields component', () => {
  it('renders expected number of inputs', () => {
    const wrapper = mount(
      <Formik>
        <ContactInfoFields handleChange={jest.fn()} name="fieldsetName" />
      </Formik>,
    );
    expect(wrapper.find('input[data-testid="firstName"]').length).toBe(1);
    expect(wrapper.find('input[data-testid="lastName"]').length).toBe(1);
    expect(wrapper.find('input[data-testid="phone"]').length).toBe(1);
    expect(wrapper.find('input[data-testid="email"]').length).toBe(1);
  });
});
