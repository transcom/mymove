import React from 'react';
import { mount } from 'enzyme';

import { AddressFields } from './AddressFields';

const wrapper = mount(<AddressFields />);

describe('AddressFields component', () => {
  it('renders', () => {
    expect(wrapper.find('input[data-cy="mailingAddress1"]').length).toBe(1);
    expect(wrapper.find('input[data-cy="mailingAddress2"]').length).toBe(1);
    expect(wrapper.find('input[data-cy="city"]').length).toBe(1);
    expect(wrapper.find('input[data-cy="state"]').length).toBe(1);
    expect(wrapper.find('input[data-cy="zip"]').length).toBe(1);
  });
});
