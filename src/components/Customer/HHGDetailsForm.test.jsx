import React from 'react';
import { mount } from 'enzyme';

import HHGDetailsForm from './HHGDetailsForm';

describe('HHGDetailsForm component', () => {
  it('renders expected form components', () => {
    const wrapper = mount(<HHGDetailsForm pageList={['page1', 'anotherPage/:foo/:bar']} pageKey="page1" />);
    expect(wrapper.find('HHGDetailsForm').length).toBe(1);
    expect(wrapper.find('DatePickerInput').length).toBe(2);
    expect(wrapper.find('AddressFields').length).toBe(1);
    expect(wrapper.find('ContactInfoFields').length).toBe(2);
    expect(wrapper.find('TextInput').length).toBe(1);
  });
  it('renders second address field when has delivery date', () => {
    const wrapper = mount(<HHGDetailsForm pageList={['page1', 'anotherPage/:foo/:bar']} pageKey="page1" />);
    wrapper.setState({ hasDeliveryAddress: true });
    expect(wrapper.find('AddressFields').length).toBe(2);
  });
});
