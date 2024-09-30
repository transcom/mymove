import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MobileHomeShipmentCard from './MobileHomeShipmentCard';

import { SHIPMENT_TYPES } from 'shared/constants';

const defaultProps = {
  showEditAndDeleteBtn: true,
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  requestedDeliveryDate: new Date('03/01/2020').toISOString(),
  pickupLocation: {
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '11111',
  },
  destinationLocation: {
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '73523',
  },
  releasingAgent: {
    firstName: 'Super',
    lastName: 'Mario',
    phone: '(555) 555-5555',
    email: 'superMario@gmail.com',
  },
  destinationZIP: '73523',
  receivingAgent: {
    firstName: 'Princess',
    lastName: 'Peach',
    phone: '(999) 999-9999',
    email: 'princessPeach@gmail.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '20fdbf58-879e-4692-b3a6-8a71f6dcfeaa',
    shipmentLocator: 'testMove123-01',
    shipmentType: SHIPMENT_TYPES.MOBILE_HOME,
    mobileHomeShipment: {
      year: 2020,
      make: 'Test Make',
      model: 'Test Model',
      lengthInInches: 240,
      widthInInches: 120,
      heightInInches: 72,
    },
  },
};

describe('MobileHomeShipmentCard component', () => {
  it('renders component with all fields', () => {
    render(<MobileHomeShipmentCard {...defaultProps} />);

    expect(screen.getAllByTestId('ShipmentCardNumber').length).toBe(1);
    expect(screen.getByText(/^#testMove123-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();

    const descriptionDefinitions = screen.getAllByRole('definition');

    const expectedRows = [
      ['Requested pickup date', '01 Jan 2020'],
      ['Pickup location', '17 8th St New York, NY 11111'],
      ['Releasing agent', 'Super Mario (555) 555-5555 superMario@gmail.com'],
      ['Requested delivery date', '01 Mar 2020'],
      ['Destination', '17 8th St New York, NY 73523'],
      ['Receiving agent', 'Princess Peach (999) 999-9999 princessPeach@gmail.com'],
      ['Mobile Home year', '2020'],
      ['Mobile Home make', 'Test Make'],
      ['Mobile Home model', 'Test Model'],
      ['Dimensions', `20' L x 10' W x 6' H`],
      [
        'Remarks',
        'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
      ],
    ];

    expect(descriptionDefinitions.length).toBe(expectedRows.length);

    expectedRows.forEach((expectedRow, index) => {
      // dt (definition terms) are not accessible elements that can be found with getByRole although
      // testing library claims this is fixed we need to find the node package that is out of date
      expect(descriptionDefinitions[index].previousElementSibling).toHaveTextContent(expectedRow[0]);
      expect(descriptionDefinitions[index]).toHaveTextContent(expectedRow[1]);
    });
  });

  it('omits the edit button when showEditAndDeleteBtn prop is false', () => {
    render(<MobileHomeShipmentCard {...defaultProps} showEditAndDeleteBtn={false} />);

    expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Delete' })).not.toBeInTheDocument();
  });

  it('calls onEditClick when edit button is pressed', async () => {
    render(<MobileHomeShipmentCard {...defaultProps} />);
    const editBtn = screen.getByRole('button', { name: 'Edit' });
    await userEvent.click(editBtn);
    expect(defaultProps.onEditClick).toHaveBeenCalledTimes(1);
  });

  it('calls onDeleteClick when delete button is pressed', async () => {
    render(<MobileHomeShipmentCard {...defaultProps} />);
    const deleteBtn = screen.getByRole('button', { name: 'Delete' });
    await userEvent.click(deleteBtn);
    expect(defaultProps.onDeleteClick).toHaveBeenCalledTimes(1);
  });

  it('renders incomplete shipment label and tooltip when shipment is incomplete', async () => {
    const incompleteShipmentProps = {
      ...defaultProps,
      shipment: {
        ...defaultProps.shipment,
        requestedPickupDate: '',
        mobileHomeShipment: defaultProps.shipment.mobileHomeShipment,
      },
      onIncompleteClick: jest.fn(),
    };

    render(<MobileHomeShipmentCard {...incompleteShipmentProps} />);

    expect(screen.getByText('Incomplete')).toBeInTheDocument();
    expect(screen.getByTitle('Help about incomplete shipment')).toBeInTheDocument();

    await userEvent.click(screen.getByTitle('Help about incomplete shipment'));
    expect(screen.getAllByTestId('ShipmentCardNumber').length).toBe(1);
  });
});
