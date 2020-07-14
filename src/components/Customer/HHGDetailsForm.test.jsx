import React from 'react';
import { mount } from 'enzyme';

import { HHGDetailsForm } from './HHGDetailsForm';

describe('HHGDetailsForm component', () => {
  it('renders expected form components', () => {
    const wrapper = mount(<HHGDetailsForm pageList={['page1', 'anotherPage/:foo/:bar']} pageKey="page1" />);
    expect(wrapper.find('HHGDetailsForm').length).toBe(1);
    expect(wrapper.find('DatePickerInput').length).toBe(2);
  });
});
