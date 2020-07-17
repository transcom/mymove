import React from 'react';
import { mount } from 'enzyme';

import { AddressFields } from './AddressFields';

const wrapper = mount(<AddressFields />);

describe('AddressFields component', () => {
  it('renders', () => {
    expect(wrapper.find('input[data-testid="mailingAddress1"]').length).toBe(1);
    expect(wrapper.find('input[data-testid="mailingAddress2"]').length).toBe(1);
    expect(wrapper.find('input[data-testid="city"]').length).toBe(1);
    expect(wrapper.find('input[data-testid="state"]').length).toBe(1);
    expect(wrapper.find('input[data-testid="zip"]').length).toBe(1);
  });
});
