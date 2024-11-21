import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import BoatShipmentCard from './BoatShipmentCard';

import { SHIPMENT_TYPES } from 'shared/constants';
import { boatShipmentTypes } from 'constants/shipments';

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
    firstName: 'Jo',
    lastName: 'Xi',
    phone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
  destinationZIP: '73523',
  receivingAgent: {
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    phone: '(999) 999-9999',
    email: 'dorothy.lagomarsino@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '20fdbf58-879e-4692-b3a6-8a71f6dcfeaa',
    shipmentLocator: 'testMove123-01',
    shipmentType: SHIPMENT_TYPES.BOAT_TOW_AWAY,
    boatShipment: {
      type: boatShipmentTypes.TOW_AWAY,
      year: 2020,
      make: 'Test Make',
      model: 'Test Model',
      lengthInInches: 240,
      widthInInches: 120,
      heightInInches: 72,
      hasTrailer: true,
      isRoadworthy: true,
    },
  },
  marketCode: 'd',
};

describe('BoatShipmentCard component', () => {
  it('renders component with all fields', () => {
    render(<BoatShipmentCard {...defaultProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(`${defaultProps.marketCode}Boat 1`);
    expect(screen.getByText(/^#testMove123-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();

    const descriptionDefinitions = screen.getAllByRole('definition');

    expect(descriptionDefinitions.length).toBe(14);

    const expectedRows = [
      ['Shipment Method', 'BTA'],
      ['Requested pickup date', '01 Jan 2020'],
      ['Pickup Address', '17 8th St New York, NY 11111'],
      ['Releasing agent', 'Jo Xi (555) 555-5555 jo.xi@email.com'],
      ['Requested delivery date', '01 Mar 2020'],
      ['Destination', '17 8th St New York, NY 73523'],
      ['Receiving agent', 'Dorothy Lagomarsino (999) 999-9999 dorothy.lagomarsino@email.com'],
      ['Boat year', '2020'],
      ['Boat make', 'Test Make'],
      ['Boat model', 'Test Model'],
      ['Dimensions', `20' L x 10' W x 6' H`],
      ['Trailer', 'Yes'],
      ['Is trailer roadworthy', 'Yes'],
      [
        'Remarks',
        'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
      ],
    ];

    expectedRows.forEach((expectedRow, index) => {
      // dt (definition terms) are not accessible elements that can be found with getByRole although
      // testing library claims this is fixed we need to find the node package that is out of date
      expect(descriptionDefinitions[index].previousElementSibling).toHaveTextContent(expectedRow[0]);
      expect(descriptionDefinitions[index]).toHaveTextContent(expectedRow[1]);
    });
  });

  it('omits the edit button when showEditAndDeleteBtn prop is false', () => {
    render(<BoatShipmentCard {...defaultProps} showEditAndDeleteBtn={false} />);

    expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Delete' })).not.toBeInTheDocument();
  });

  it('calls onEditClick when edit button is pressed', async () => {
    render(<BoatShipmentCard {...defaultProps} />);
    const editBtn = screen.getByRole('button', { name: 'Edit' });
    await userEvent.click(editBtn);
    expect(defaultProps.onEditClick).toHaveBeenCalledTimes(1);
  });

  it('calls onDeleteClick when delete button is pressed', async () => {
    render(<BoatShipmentCard {...defaultProps} />);
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
        boatShipment: defaultProps.shipment.boatShipment,
      },
      onIncompleteClick: jest.fn(),
    };

    render(<BoatShipmentCard {...incompleteShipmentProps} />);

    expect(screen.getByText('Incomplete')).toBeInTheDocument();
    expect(screen.getByTitle('Help about incomplete shipment')).toBeInTheDocument();

    await userEvent.click(screen.getByTitle('Help about incomplete shipment'));

    expect(incompleteShipmentProps.onIncompleteClick).toHaveBeenCalledWith(
      'Boat 1',
      'testMove123-01',
      SHIPMENT_TYPES.BOAT_TOW_AWAY,
    );
  });
});
