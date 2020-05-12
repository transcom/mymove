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
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
  backupContactName: 'Quinn Ocampo',
  backupContactPhone: '+1 999-999-9999',
  backupContactEmail: 'quinnocampo@myemail.com',
};

describe('Customer Info Table', () => {
  it('should render the data passed to its props', () => {
    const wrapper = shallow(<CustomerInfoTable customerInfo={info} />);
    expect(wrapper.find({ 'data-cy': 'name' }).text()).toMatch(info.name);
    expect(wrapper.find({ 'data-cy': 'dodId' }).text()).toMatch(info.dodId);
    expect(wrapper.find({ 'data-cy': 'phone' }).text()).toMatch(info.phone);
    expect(wrapper.find({ 'data-cy': 'email' }).text()).toMatch(info.email);
    expect(wrapper.find({ 'data-cy': 'currentAddress' }).text()).toMatch(
      `${info.currentAddress.street_address_1}, ${info.currentAddress.city}, ${info.currentAddress.state} ${info.currentAddress.postal_code}`,
    );
    expect(wrapper.find({ 'data-cy': 'destinationAddress' }).text()).toMatch(
      `${info.destinationAddress.street_address_1}, ${info.destinationAddress.city}, ${info.destinationAddress.state} ${info.destinationAddress.postal_code}`,
    );
    expect(wrapper.find({ 'data-cy': 'backupContactName' }).text()).toMatch(info.backupContactName);
    expect(wrapper.find({ 'data-cy': 'backupContactPhone' }).text()).toMatch(info.backupContactPhone);
    expect(wrapper.find({ 'data-cy': 'backupContactEmail' }).text()).toMatch(info.backupContactEmail);
  });
});
