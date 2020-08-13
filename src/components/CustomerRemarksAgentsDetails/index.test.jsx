import React from 'react';
import { shallow } from 'enzyme';

import CustomerRemarksAgentsDetails from './index';

import DataPoint from 'components/DataPoint';

describe('CustomerRemarksAgentsDetails', () => {
  it('renders empty details', () => {
    const wrapper = shallow(<CustomerRemarksAgentsDetails />);
    expect(wrapper.children()).toHaveLength(3);
  });

  it('renders with customer remarks details', () => {
    const index = 0;
    const customerRemarks = 'This is a remark.';
    const wrapper = shallow(<CustomerRemarksAgentsDetails customerRemarks={customerRemarks} />);
    expect(wrapper.find(DataPoint).at(index).dive().text()).toContain(customerRemarks);
  });

  it('renders with releasing agent details', () => {
    const index = 1;
    const releasingAgent = {
      firstName: 'firstname',
      lastName: 'lastname',
      phone: '(111) 111-1111',
      email: 'test@test.com',
    };
    const wrapper = shallow(<CustomerRemarksAgentsDetails releasingAgent={releasingAgent} />);
    const releasingSection = wrapper.find(DataPoint).at(index).dive().text();
    expect(releasingSection).toContain(releasingAgent.firstName);
    expect(releasingSection).toContain(releasingAgent.lastName);
    expect(releasingSection).toContain(releasingAgent.phone);
    expect(releasingSection).toContain(releasingAgent.email);
  });

  it('renders with receiving agent details', () => {
    const index = 2;
    const receivingAgent = {
      firstName: 'firstname',
      lastName: 'lastname',
      phone: '(111) 111-1111',
      email: 'test@test.com',
    };
    const wrapper = shallow(<CustomerRemarksAgentsDetails receivingAgent={receivingAgent} />);
    const receivingSection = wrapper.find(DataPoint).at(index).dive().text();
    expect(receivingSection).toContain(receivingAgent.firstName);
    expect(receivingSection).toContain(receivingAgent.lastName);
    expect(receivingSection).toContain(receivingAgent.phone);
    expect(receivingSection).toContain(receivingAgent.email);
  });
});
