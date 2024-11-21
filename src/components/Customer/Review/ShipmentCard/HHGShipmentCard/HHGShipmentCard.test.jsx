/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard/HHGShipmentCard';
import { formatCustomerDate } from 'utils/formatters';
import { shipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const defaultProps = {
  moveId: 'testMove123',
  editPath: '',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipmentId: '#ABC123K',
  shipmentLocator: '#ABC123K-01',
  marketCode: 'i',
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
  marketCode: 'd',
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
  marketCode: 'd',
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
      'Pickup Address',
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
    const tableHeaders = ['Requested pickup date', 'Pickup Address', 'Requested delivery date', 'Destination'];
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

  it('should not render a secondary Pickup Address if not provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} />);

    const secondPickupLocation = await screen.queryByText('Second Pickup Address');
    expect(secondPickupLocation).not.toBeInTheDocument();
  });

  it('should not render a secondary delivery address if not provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} />);

    const secondDestination = await screen.queryByText('Second Destination');
    expect(secondDestination).not.toBeInTheDocument();
  });

  it('should render a secondary Pickup Address if provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} {...secondaryPickupAddress} />);

    const secondPickupLocation = await screen.getByText('Second Pickup Address');
    expect(secondPickupLocation).toBeInTheDocument();
    const secondPickupLocationInformation = await screen.getByText(/Some Other Street Name/);
    expect(secondPickupLocationInformation).toBeInTheDocument();
  });

  it('should render a secondary delivery address if provided one', async () => {
    render(<HHGShipmentCard {...defaultProps} {...secondaryDeliveryAddress} />);

    const secondDestination = await screen.getByText('Second Destination');
    expect(secondDestination).toBeInTheDocument();
    const secondDesintationInformation = await screen.getByText(/Some Street Name/);
    expect(secondDesintationInformation).toBeInTheDocument();
  });

  it('renders HHGShipmentCard with a heading that has a market code and shipment type', async () => {
    render(<HHGShipmentCard {...defaultProps} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(`${defaultProps.marketCode}HHG 1`);
  });

  it('does not render incomplete label and tooltip icon for completed hhg shipment with SUBMITTED status', async () => {
    render(<HHGShipmentCard {...completeProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('HHG 1');
    expect(screen.getByText(/^#ABC123K-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.queryByText('Incomplete')).toBeNull();
  });

  it('renders complete HHGShipmentCard with a heading that has a market code and shipment type', async () => {
    render(<HHGShipmentCard {...completeProps} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(`${completeProps.marketCode}HHG 1`);
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

  it('renders incomplete HHGShipmentCard with a heading that has a market code and shipment type', async () => {
    render(<HHGShipmentCard {...incompleteProps} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(`${incompleteProps.marketCode}HHG 1`);
  });
});

const ubProps = {
  moveId: 'testMove123',
  editPath: '',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipmentId: '#ABC123K',
  shipmentLocator: '#ABC123K-01',
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
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

const incompleteUBProps = {
  moveId: 'testMove123',
  editPath: '',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  onIncompleteClick: mockedOnIncompleteClickFunction,
  shipmentNumber: 1,
  shipmentId: 'ABC123K',
  shipmentLocator: 'ABC123K-01',
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
  showEditAndDeleteBtn: false,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  status: shipmentStatuses.DRAFT,
};

const completeUBProps = {
  moveId: 'testMove123',
  editPath: '',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipmentId: 'ABC123K',
  shipmentLocator: 'ABC123K-01',
  shipmentType: SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE,
  showEditAndDeleteBtn: false,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  status: shipmentStatuses.SUBMITTED,
};

function mountHHGShipmentCardForUBShipment(props) {
  return mount(<HHGShipmentCard {...ubProps} {...props} />);
}

describe('HHGShipmentCard component can be reused for UB shipment card', () => {
  it('renders component with all fields', () => {
    const wrapper = mountHHGShipmentCardForUBShipment();
    const tableHeaders = [
      'Requested pickup date',
      'Pickup Address',
      'Releasing agent',
      'Requested delivery date',
      'Destination',
      'Receiving agent',
      'Remarks',
    ];
    const { streetAddress1, city, state, postalCode } = ubProps.pickupLocation;
    const {
      firstName: releasingFirstName,
      lastName: releasingLastName,
      phone: releasingTelephone,
      email: releasingEmail,
    } = ubProps.releasingAgent;
    const {
      firstName: receivingFirstName,
      lastName: receivingLastName,
      phone: receivingTelephone,
      email: receivingEmail,
    } = ubProps.receivingAgent;
    const tableData = [
      formatCustomerDate(ubProps.requestedPickupDate),
      `${streetAddress1} ${city}, ${state} ${postalCode}`,
      `${releasingFirstName} ${releasingLastName} ${releasingTelephone} ${releasingEmail}`,
      formatCustomerDate(ubProps.requestedDeliveryDate),
      ubProps.destinationZIP,
      `${receivingFirstName} ${receivingLastName} ${receivingTelephone} ${receivingEmail}`,
    ];

    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').text()).toBe(ubProps.remarks);
  });

  it('should render UB shipment card without releasing/receiving agents and remarks', () => {
    const wrapper = mountHHGShipmentCardForUBShipment({
      ...ubProps,
      releasingAgent: null,
      receivingAgent: null,
      remarks: '',
    });
    const tableHeaders = ['Requested pickup date', 'Pickup Address', 'Requested delivery date', 'Destination'];
    const { streetAddress1, city, state, postalCode } = ubProps.pickupLocation;
    const tableData = [
      formatCustomerDate(ubProps.requestedPickupDate),
      `${streetAddress1} ${city}, ${state} ${postalCode}`,
      formatCustomerDate(ubProps.requestedDeliveryDate),
      ubProps.destinationZIP,
    ];
    tableHeaders.forEach((label, index) => expect(wrapper.find('dt').at(index).text()).toBe(label));
    tableData.forEach((label, index) => expect(wrapper.find('dd').at(index).text()).toBe(label));
    expect(wrapper.find('.remarksCell').length).toBe(0);
  });

  it('should not render a secondary Pickup Address on UB shipment card if not provided one', async () => {
    render(<HHGShipmentCard {...ubProps} />);

    const secondPickupLocation = await screen.queryByText('Second Pickup Address');
    expect(secondPickupLocation).not.toBeInTheDocument();
  });

  it('should not render a secondary delivery address on UB shipment card if not provided one', async () => {
    render(<HHGShipmentCard {...ubProps} />);

    const secondDestination = await screen.queryByText('Second Destination');
    expect(secondDestination).not.toBeInTheDocument();
  });

  it('should render a UB shipment card secondary Pickup Address if provided one', async () => {
    render(<HHGShipmentCard {...ubProps} {...secondaryPickupAddress} />);

    const secondPickupLocation = await screen.getByText('Second Pickup Address');
    expect(secondPickupLocation).toBeInTheDocument();
    const secondPickupLocationInformation = await screen.getByText(/Some Other Street Name/);
    expect(secondPickupLocationInformation).toBeInTheDocument();
  });

  it('should render a UB shipment card secondary delivery address if provided one', async () => {
    render(<HHGShipmentCard {...ubProps} {...secondaryDeliveryAddress} />);

    const secondDestination = await screen.getByText('Second Destination');
    expect(secondDestination).toBeInTheDocument();
    const secondDesintationInformation = await screen.getByText(/Some Street Name/);
    expect(secondDesintationInformation).toBeInTheDocument();
  });

  it('does not render UB shipment card incomplete label and tooltip icon for completed UB shipment with SUBMITTED status', async () => {
    render(<HHGShipmentCard {...completeUBProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('UB 1');
    expect(screen.getByText(/^#ABC123K-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.queryByText('Incomplete')).toBeNull();
  });

  it('renders incomplete label and tooltip icon for incomplete UB shipment with DRAFT status', async () => {
    render(<HHGShipmentCard {...incompleteUBProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('UB 1');
    expect(screen.getByText(/^#ABC123K-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByText(/^Incomplete$/, { selector: 'span' })).toBeInTheDocument();

    expect(screen.getByTitle('Help about incomplete shipment')).toBeInTheDocument();
    await userEvent.click(screen.getByTitle('Help about incomplete shipment'));

    // verify onclick is getting json string as parameter
    expect(mockedOnIncompleteClickFunction).toHaveBeenCalledWith('UB 1', 'ABC123K-01', 'UNACCOMPANIED_BAGGAGE');
  });
});
