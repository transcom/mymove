import React from 'react';
import { shallow } from 'enzyme';
import CustomerRemarksAgentsDetails from '.';

describe('CustomerRemarksAgentsDetails', () => {
  it('should render empty details', () => {
    const wrapper = shallow(<CustomerRemarksAgentsDetails />);
    expect(wrapper.children()).toHaveLength(3);
  });

  it('should render with customer remarks details', () => {
    const index = 0;
    const customerRemarks = 'This is a remark.';
    const wrapper = shallow(<CustomerRemarksAgentsDetails customerRemarks={customerRemarks} />);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(customerRemarks);
  });

  it('should render with releasing agent details', () => {
    const index = 1;
    const releasingAgent = {
      firstName: 'firstname',
      lastName: 'lastname',
      phone: '(111) 111-1111',
      email: 'test@test.com',
    };
    const wrapper = shallow(<CustomerRemarksAgentsDetails releasingAgent={releasingAgent} />);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(releasingAgent.firstName);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(releasingAgent.lastName);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(releasingAgent.phone);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(releasingAgent.email);
  });

  it('should render with receiving agent details', () => {
    const index = 2;
    const receivingAgent = {
      firstName: 'firstname',
      lastName: 'lastname',
      phone: '(111) 111-1111',
      email: 'test@test.com',
    };
    const wrapper = shallow(<CustomerRemarksAgentsDetails receivingAgent={receivingAgent} />);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(receivingAgent.firstName);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(receivingAgent.lastName);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(receivingAgent.phone);
    expect(wrapper.find('div[className="container--small"]').at(index).text()).toContain(receivingAgent.email);
  });
});
