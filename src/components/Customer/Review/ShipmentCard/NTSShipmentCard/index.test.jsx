/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import NTSShipmentCard from './index';

import { formatCustomerDate } from 'utils/formatters';

const defaultProps = {
  moveId: 'testMove123',
  onEditClick: jest.fn(),
  shipmentId: '#ABC123K',
  shipmentType: 'HHG_INTO_NTS_DOMESTIC',
  showEditBtn: false,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  pickupLocation: {
    street_address_1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postal_code: '11111',
  },
  releasingAgent: {
    firstName: 'Jo',
    lastName: 'Xi',
    phone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

function mountNTSShipmentCard(props) {
  return mount(<NTSShipmentCard {...defaultProps} {...props} />);
}

const secondaryPickupAddress = {
  secondaryPickupAddress: {
    street_address_1: 'Some Other Street Name',
    city: 'New York',
    state: 'NY',
    postal_code: '111111',
  },
};

describe('NTSShipmentCard component', () => {
  it('renders component with all fields', () => {
    const wrapper = mountNTSShipmentCard();
    const tableHeaders = ['Requested pickup date', 'Pickup location', 'Releasing agent', 'Remarks'];
    const { street_address_1: streetAddress1, city, state, postal_code: postalCode } = defaultProps.pickupLocation;
    const {
      firstName: releasingFirstName,
      lastName: releasingLastName,
      phone: releasingTelephone,
      email: releasingEmail,
    } = defaultProps.releasingAgent;
    const tableData = [
      formatCustomerDate(defaultProps.requestedPickupDate),
      `${streetAddress1} ${city}, ${state} ${postalCode}`,
      `${releasingFirstName} ${releasingLastName} ${releasingTelephone} ${releasingEmail}`,
    ];

    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').text()).toBe(defaultProps.remarks);
  });

  it('should render without releasing/receiving agents and remarks', () => {
    const wrapper = mountNTSShipmentCard({ ...defaultProps, releasingAgent: null, remarks: '' });
    const tableHeaders = ['Requested pickup date', 'Pickup location'];
    const { street_address_1: streetAddress1, city, state, postal_code: postalCode } = defaultProps.pickupLocation;
    const tableData = [
      formatCustomerDate(defaultProps.requestedPickupDate),
      `${streetAddress1} ${city}, ${state} ${postalCode}`,
    ];
    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').length).toBe(0);
  });

  it('should not render a secondary pickup location if not provided one', async () => {
    render(<NTSShipmentCard {...defaultProps} />);

    const secondPickupLocation = await screen.queryByText('Second pickup location');
    expect(secondPickupLocation).not.toBeInTheDocument();
  });

  it('should render a secondary pickup location if provided one', async () => {
    render(<NTSShipmentCard {...defaultProps} {...secondaryPickupAddress} />);

    const secondPickupLocation = await screen.getByText('Second pickup location');
    expect(secondPickupLocation).toBeInTheDocument();
    const secondPickupLocationInformation = await screen.getByText(/Some Other Street Name/);
    expect(secondPickupLocationInformation).toBeInTheDocument();
  });
});
