/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard/HHGShipmentCard';
import { formatCustomerDate } from 'utils/formatters';
import { shipmentStatuses } from 'constants/shipments';

const defaultProps = {
  moveId: 'testMove123',
  editPath: '',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipmentId: '#ABC123K',
  shipmentLocator: '#ABC123K-01',
  shipmentType: 'HHG',
  showEditAndDeleteBtn: false,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  pickupLocation: {
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '11111',
  },
  releasingAgent: {
    firstName: 'Jo',
    lastName: 'Xi',
    phone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
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

const mockedOnIncompleteClickFunction = jest.fn();
const incompleteProps = {
  moveId: 'testMove123',
  editPath: '',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  onIncompleteClick: mockedOnIncompleteClickFunction,
  shipmentNumber: 1,
  shipmentId: 'ABC123K',
  shipmentLocator: 'ABC123K-01',
  shipmentType: 'HHG',
  showEditAndDeleteBtn: false,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  status: shipmentStatuses.DRAFT,
};

const completeProps = {
  moveId: 'testMove123',
  editPath: '',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipmentId: 'ABC123K',
  shipmentLocator: 'ABC123K-01',
  shipmentType: 'HHG',
  showEditAndDeleteBtn: false,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  status: shipmentStatuses.SUBMITTED,
};

const secondaryDeliveryAddress = {
  secondaryDeliveryAddress: {
    streetAddress1: 'Some Street Name',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
};

const secondaryPickupAddress = {
  secondaryPickupAddress: {
    streetAddress1: 'Some Other Street Name',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
};

function mountHHGShipmentCard(props) {
  return mount(<HHGShipmentCard {...defaultProps} {...props} />);
}

describe('HHGShipmentCard component', () => {
  it('renders component with all fields', () => {
    const wrapper = mountHHGShipmentCard();
    const tableHeaders = [
      'Requested pickup date',
      'Pickup location',
      'Releasing agent',
      'Requested delivery date',
      'Destination',
      'Receiving agent',
      'Remarks',
    ];
    const { streetAddress1, city, state, postalCode } = defaultProps.pickupLocation;
    const {
      firstName: releasingFirstName,
      lastName: releasingLastName,
      phone: releasingTelephone,
      email: releasingEmail,
    } = defaultProps.releasingAgent;
    const {
      firstName: receivingFirstName,
      lastName: receivingLastName,
      phone: receivingTelephone,
      email: receivingEmail,
    } = defaultProps.receivingAgent;
    const tableData = [
      formatCustomerDate(defaultProps.requestedPickupDate),
      `${streetAddress1} ${city}, ${state} ${postalCode}`,
      `${releasingFirstName} ${releasingLastName} ${releasingTelephone} ${releasingEmail}`,
      formatCustomerDate(defaultProps.requestedDeliveryDate),
      defaultProps.destinationZIP,
      `${receivingFirstName} ${receivingLastName} ${receivingTelephone} ${receivingEmail}`,
    ];

    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').text()).toBe(defaultProps.remarks);
  });

  it('should render without releasing/receiving agents and remarks', () => {
    const wrapper = mountHHGShipmentCard({ ...defaultProps, releasingAgent: null, receivingAgent: null, remarks: '' });
    const tableHeaders = ['Requested pickup date', 'Pickup location', 'Requested delivery date', 'Destination'];
    const { streetAddress1, city, state, postalCode } = defaultProps.pickupLocation;
    const tableData = [
      formatCustomerDate(defaultProps.requestedPickupDate),
      `${streetAddress1} ${city}, ${state} ${postalCode}`,
      formatCustomerDate(defaultProps.requestedDeliveryDate),
      defaultProps.destinationZIP,
    ];
    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').length).toBe(0);
  });

  it('should not render a secondary pickup location if not provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} />);

    const secondPickupLocation = await screen.queryByText('Second pickup location');
    expect(secondPickupLocation).not.toBeInTheDocument();
  });

  it('should not render a secondary destination location if not provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} />);

    const secondDestination = await screen.queryByText('Second Destination');
    expect(secondDestination).not.toBeInTheDocument();
  });

  it('should render a secondary pickup location if provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} {...secondaryPickupAddress} />);

    const secondPickupLocation = await screen.getByText('Second pickup location');
    expect(secondPickupLocation).toBeInTheDocument();
    const secondPickupLocationInformation = await screen.getByText(/Some Other Street Name/);
    expect(secondPickupLocationInformation).toBeInTheDocument();
  });

  it('should render a secondary destination location if provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} {...secondaryDeliveryAddress} />);

    const secondDestination = await screen.getByText('Second Destination');
    expect(secondDestination).toBeInTheDocument();
    const secondDesintationInformation = await screen.getByText(/Some Street Name/);
    expect(secondDesintationInformation).toBeInTheDocument();
  });

  it('does not render incomplete label and tooltip icon for completed hhg shipment with SUBMITTED status', async () => {
    render(<HHGShipmentCard {...completeProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('HHG 1');
    expect(screen.getByText(/^#ABC123K-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.queryByText('Incomplete')).toBeNull();
  });

  it('renders incomplete label and tooltip icon for incomplete HHG shipment with DRAFT status', async () => {
    render(<HHGShipmentCard {...incompleteProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('HHG 1');
    expect(screen.getByText(/^#ABC123K-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByText(/^Incomplete$/, { selector: 'span' })).toBeInTheDocument();

    expect(screen.getByTitle('Help about incomplete shipment')).toBeInTheDocument();
    await userEvent.click(screen.getByTitle('Help about incomplete shipment'));

    // verify onclick is getting json string as parameter
    expect(mockedOnIncompleteClickFunction).toHaveBeenCalledWith('HHG 1', 'ABC123K-01', 'HHG');
  });
});
