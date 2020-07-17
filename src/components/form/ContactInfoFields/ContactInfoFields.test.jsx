import React from 'react';
import { mount } from 'enzyme';

import { ContactInfoFields } from './ContactInfoFields';

describe('ContactInfoFields component', () => {
  it('renders expected number of inputs', () => {
    const wrapper = mount(<ContactInfoFields />);
    expect(wrapper.find('input[data-cy="firstName"]').length).toBe(1);
    expect(wrapper.find('input[data-cy="lastName"]').length).toBe(1);
    expect(wrapper.find('input[data-cy="phone"]').length).toBe(1);
    expect(wrapper.find('input[data-cy="email"]').length).toBe(1);
  });
});
