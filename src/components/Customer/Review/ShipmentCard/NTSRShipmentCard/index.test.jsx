/*  react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import NTSRShipmentCard from '.';

import { formatCustomerDate } from 'utils/formatters';

const defaultProps = {
  shipmentId: '#ABC123K',
  requestedDeliveryDate: new Date('03/01/2020').toISOString(),
  destinationZIP: '73523',
  receivingAgent: {
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    phone: '(999) 999-9999',
    email: 'dorothy.lagomarsino@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

function mountNTSRShipmentCard(props = defaultProps) {
  return mount(<NTSRShipmentCard {...props} />);
}
describe('NTSRShipmentCard component', () => {
  it('renders component with all fields', () => {
    const wrapper = mountNTSRShipmentCard();
    const tableHeaders = ['Requested delivery date', 'Destination', 'Receiving agent', 'Remarks'];
    const {
      firstName: receivingFirstName,
      lastName: receivingLastName,
      phone: receivingTelephone,
      email: receivingEmail,
    } = defaultProps.receivingAgent;
    const tableData = [
      formatCustomerDate(defaultProps.requestedDeliveryDate),
      defaultProps.destinationZIP,
      `${receivingFirstName} ${receivingLastName} ${receivingTelephone} ${receivingEmail}`,
    ];

    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').text()).toBe(defaultProps.remarks);
  });

  it('should render without releasing/receiving agents and remarks', () => {
    const wrapper = mountNTSRShipmentCard({ ...defaultProps, releasingAgent: null, remarks: '' });
    const tableHeaders = ['Requested delivery date', 'Destination'];
    const tableData = [formatCustomerDate(defaultProps.requestedDeliveryDate), defaultProps.destinationZIP];
    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').length).toBe(0);
  });
});
