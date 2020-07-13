import React from 'react';
import { mount } from 'enzyme';

import { ContactInfoFields } from './ContactInfoFields';

const wrapper = mount(<ContactInfoFields />);

describe('ContactInfoFields component', () => {
  it('renders', () => {
    expect(wrapper.find('ContactInfoFields').length).toBe(1);
  });
});
