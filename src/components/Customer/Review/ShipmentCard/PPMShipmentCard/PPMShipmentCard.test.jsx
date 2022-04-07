import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import PPMShipmentCard from './PPMShipmentCard';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const defaultProps = {
  showEditAndDeleteBtn: true,
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '20fdbf58-879e-4692-b3a6-8a71f6dcfeaa',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupPostalCode: '10001',
      destinationPostalCode: '11111',
      sitExpected: false,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
    },
  },
};

const completeProps = {
  showEditAndDeleteBtn: true,
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '20fdbf58-879e-4692-b3a6-8a71f6dcfeaa',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupPostalCode: '10001',
      secondaryPickupPostalCode: '10002',
      destinationPostalCode: '11111',
      secondaryDestinationPostalCode: '22222',
      sitExpected: true,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
      estimatedWeight: 5999,
      proGearWeight: 1250,
      spouseProGearWeight: 375,
      estimatedIncentive: 1000099,
      advanceRequested: true,
      advance: 600000,
    },
  },
};

describe('PPMShipmentCard component', () => {
  it('renders component with all fields', () => {
    render(<PPMShipmentCard {...completeProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('PPM 1');
    expect(screen.getByText(/^#20FDBF58$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();

    const descriptionDefinitions = screen.getAllByRole('definition');

    expect(descriptionDefinitions.length).toBe(11);

    const expectedRows = [
      ['Expected departure', '01 Jan 2020'],
      ['Origin ZIP', '10001'],
      ['Second origin ZIP', '10002'],
      ['Destination ZIP', '11111'],
      ['Second destination ZIP', '22222'],
      ['Storage expected? (SIT)', 'Yes'],
      ['Estimated weight', '5,999 lbs'],
      ['Pro-gear', 'Yes, 1,250 lbs'],
      ['Spouse pro-gear', 'Yes, 375 lbs'],
      ['Estimated incentive', '$10,000'],
      ['Advance', 'Yes, $6,000'],
    ];

    expectedRows.forEach((expectedRow, index) => {
      // dt (definition terms) are not accessible elements that can be found with getByRole although
      // testing library claims this is fixed we need to find the node package that is out of date
      expect(descriptionDefinitions[index].previousElementSibling).toHaveTextContent(expectedRow[0]);
      expect(descriptionDefinitions[index]).toHaveTextContent(expectedRow[1]);
    });

    expect();
  });

  it('renders component with incomplete fields', () => {
    render(<PPMShipmentCard {...defaultProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('PPM 1');
    expect(screen.getByText(/^#20FDBF58$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();

    const descriptionDefinitions = screen.getAllByRole('definition');

    expect(descriptionDefinitions.length).toBe(9);

    const expectedRows = [
      ['Expected departure', '01 Jan 2020'],
      ['Origin ZIP', '10001'],
      ['Destination ZIP', '11111'],
      ['Storage expected? (SIT)', 'No'],
      ['Estimated weight', '0 lbs'],
      ['Pro-gear', 'No'],
      ['Spouse pro-gear', 'No'],
      ['Estimated incentive', '$0'],
      ['Advance', 'No'],
    ];

    expectedRows.forEach((expectedRow, index) => {
      // dt (definition terms) are not accessible elements that can be found with getByRole although
      // testing library claims this is fixed we need to find the node package that is out of date
      expect(descriptionDefinitions[index].previousElementSibling).toHaveTextContent(expectedRow[0]);
      expect(descriptionDefinitions[index]).toHaveTextContent(expectedRow[1]);
    });
  });

  it('omits the edit button when showEditAndDeleteBtn prop is false', () => {
    render(<PPMShipmentCard {...completeProps} showEditAndDeleteBtn={false} />);

    expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Delete' })).not.toBeInTheDocument();
  });

  it('calls onEditClick when edit button is pressed', () => {
    render(<PPMShipmentCard {...completeProps} />);
    const editBtn = screen.queryByRole('button', { name: 'Edit' });
    fireEvent.click(editBtn);
    expect(completeProps.onEditClick).toHaveBeenCalledTimes(1);
  });

  it('calls onDeleteClick when delete button is pressed', () => {
    render(<PPMShipmentCard {...completeProps} />);
    const deleteBtn = screen.queryByRole('button', { name: 'Delete' });
    fireEvent.click(deleteBtn);
    expect(completeProps.onDeleteClick).toHaveBeenCalledTimes(1);
  });
});
