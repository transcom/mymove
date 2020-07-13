import React from 'react';
import { mount } from 'enzyme';

import { HHGDetailsForm } from './HHGDetailsForm';

const wrapper = mount(<HHGDetailsForm pageKey="" pages={[]} />);

describe('HHGDetailsForm component', () => {
  it('renders', () => {
    expect(wrapper.find('HHGDetailsForm').length).toBe(1);
  });
});
