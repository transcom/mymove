import React from 'react';
import { mount } from 'enzyme';

import { AddressFields } from './AddressFields';

const wrapper = mount(<AddressFields />);

describe('AddressFields component', () => {
  it('renders', () => {
    expect(wrapper.find('AddressFields').length).toBe(1);
  });
});
