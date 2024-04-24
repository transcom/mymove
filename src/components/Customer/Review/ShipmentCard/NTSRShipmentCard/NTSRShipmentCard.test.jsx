/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import NTSRShipmentCard from 'components/Customer/Review/ShipmentCard/NTSRShipmentCard/NTSRShipmentCard';
import { formatCustomerDate } from 'utils/formatters';
import { shipmentStatuses } from 'constants/shipments';

const defaultProps = {
  moveId: 'testMove123',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  onIncompleteClick: jest.fn(),
  showEditAndDeleteBtn: false,
  shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
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

const completeProps = {
  moveId: 'testMove123',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  onIncompleteClick: jest.fn(),
  showEditAndDeleteBtn: false,
  shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
  shipmentId: 'ABC123K',
  shipmentLocator: 'ABC123K-01',
  status: shipmentStatuses.SUBMITTED,
};

const mockedOnIncompleteClickFunction = jest.fn();
const incompleteProps = {
  moveId: 'testMove123',
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  onIncompleteClick: mockedOnIncompleteClickFunction,
  showEditAndDeleteBtn: false,
  shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
  shipmentId: 'ABC123K',
  shipmentLocator: 'ABC123K-01',
  status: shipmentStatuses.DRAFT,
};

const secondaryDeliveryAddress = {
  secondaryDeliveryAddress: {
    streetAddress1: 'Some Street Name',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
};

function mountNTSRShipmentCard(props) {
  return mount(<NTSRShipmentCard {...defaultProps} {...props} />);
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
    expect(wrapper.find('.remarksCell').at(0).text()).toBe('—');
  });

  it('should not render a secondary destination location if not provided one', async () => {
    render(<NTSRShipmentCard {...defaultProps} />);

    const secondDestination = await screen.queryByText('Second Destination');
    expect(secondDestination).not.toBeInTheDocument();
  });

  it('should render a secondary destination location if provided one', async () => {
    render(<NTSRShipmentCard {...defaultProps} {...secondaryDeliveryAddress} />);

    const secondDestination = await screen.getByText('Second Destination');
    expect(secondDestination).toBeInTheDocument();
    const secondDesintationInformation = await screen.getByText(/Some Street Name/);
    expect(secondDesintationInformation).toBeInTheDocument();
  });

  it('does not render incomplete label and tooltip icon for completed shipment with SUBMITTED status', async () => {
    render(<NTSRShipmentCard {...completeProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('NTS-release');
    expect(screen.getByText(/^#ABC123K-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.queryByText('Incomplete')).toBeNull();
  });

  it('renders incomplete label and tooltip icon for incomplete HHG shipment with DRAFT status', async () => {
    render(<NTSRShipmentCard {...incompleteProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('NTS-release');
    expect(screen.getByText(/^#ABC123K-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByText(/^Incomplete$/, { selector: 'span' })).toBeInTheDocument();

    expect(screen.getByTitle('Help about incomplete shipment')).toBeInTheDocument();
    await userEvent.click(screen.getByTitle('Help about incomplete shipment'));

    // verify onclick is getting json string as parameter
    expect(mockedOnIncompleteClickFunction).toHaveBeenCalledWith('NTS-release', 'ABC123K-01', 'NTS-release');
  });
});
