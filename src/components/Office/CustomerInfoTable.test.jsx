import React from 'react';
import { shallow } from 'enzyme';

import CustomerInfoTable from './CustomerInfoTable';

const info = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  backupContact: {
    name: 'Quinn Ocampo',
    email: 'quinnocampo@myemail.com',
    phone: '+1 999-999-9999',
  },
};

describe('Customer Info Table', () => {
  it('should render the data passed to its props', () => {
    const wrapper = shallow(<CustomerInfoTable customerInfo={info} />);
    expect(wrapper.find({ 'data-testid': 'name' }).text()).toMatch(info.name);
    expect(wrapper.find({ 'data-testid': 'dodId' }).text()).toMatch(info.dodId);
    expect(wrapper.find({ 'data-testid': 'phone' }).text()).toMatch(info.phone);
    expect(wrapper.find({ 'data-testid': 'email' }).text()).toMatch(info.email);
    expect(wrapper.find({ 'data-testid': 'currentAddress' }).text()).toMatch(
      `${info.currentAddress.street_address_1}, ${info.currentAddress.city}, ${info.currentAddress.state} ${info.currentAddress.postal_code}`,
    );
    expect(wrapper.find({ 'data-testid': 'backupContactName' }).text()).toMatch(info.backupContact.name);
    expect(wrapper.find({ 'data-testid': 'backupContactPhone' }).text()).toMatch(info.backupContact.phone);
    expect(wrapper.find({ 'data-testid': 'backupContactEmail' }).text()).toMatch(info.backupContact.email);
  });
});
