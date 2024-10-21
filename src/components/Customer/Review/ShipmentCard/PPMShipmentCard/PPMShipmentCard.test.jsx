import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PPMShipmentCard from './PPMShipmentCard';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import transportationOfficeFactory from 'utils/test/factories/transportationOffice';
import affiliations from 'content/serviceMemberAgencies';

const defaultProps = {
  showEditAndDeleteBtn: true,
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '20fdbf58-879e-4692-b3a6-8a71f6dcfeaa',
    shipmentLocator: 'testMove123-01',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '11111',
      },
      sitExpected: false,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
    },
  },
  marketCode: 'd',
};

const completeProps = {
  showEditAndDeleteBtn: true,
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  shipmentNumber: 1,
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '20fdbf58-879e-4692-b3a6-8a71f6dcfeaa',
    shipmentLocator: 'testMove123-01',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '10001',
      },
      secondaryPickupAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '10002',
      },
      destinationAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '11111',
      },
      secondaryDestinationAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '22222',
      },
      sitExpected: true,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
      estimatedWeight: 5999,
      proGearWeight: 1250,
      spouseProGearWeight: 375,
      estimatedIncentive: 1000099,
      hasRequestedAdvance: true,
      advanceAmountRequested: 600000,
    },
  },
  marketCode: 'd',
};

const mockedOnIncompleteClickFunction = jest.fn();
const incompleteProps = {
  showEditAndDeleteBtn: true,
  onEditClick: jest.fn(),
  onDeleteClick: jest.fn(),
  onIncompleteClick: mockedOnIncompleteClickFunction,
  shipmentNumber: 1,
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '20fdbf58-879e-4692-b3a6-8a71f6dcfeaa',
    shipmentLocator: 'testMove123-01',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '111 Test Street',
        streetAddress2: '222 Test Street',
        streetAddress3: 'Test Man',
        city: 'Test City',
        state: 'NY',
        postalCode: '11111',
      },
      sitExpected: false,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
      hasRequestedAdvance: null,
    },
  },
  marketCode: 'd',
};

describe('PPMShipmentCard component', () => {
  it('renders component with all fields', () => {
    render(<PPMShipmentCard {...completeProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('PPM 1');
    expect(screen.getByText(/^#testMove123-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();

    const descriptionDefinitions = screen.getAllByRole('definition');

    expect(descriptionDefinitions.length).toBe(11);

    const expectedRows = [
      ['Expected departure', '01 Jan 2020'],
      ['Origin address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 10001'],
      ['Second origin address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 10002'],
      ['Destination address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 11111'],
      ['Second destination address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 22222'],
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

  it('renders complete PPMShipmentCard with a heading that has a market code and shipment type', async () => {
    render(<PPMShipmentCard {...completeProps} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(`${completeProps.marketCode}PPM`);
  });

  it('renders component with incomplete fields', () => {
    render(<PPMShipmentCard {...defaultProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('PPM 1');
    expect(screen.getByText(/^#testMove123-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();

    const descriptionDefinitions = screen.getAllByRole('definition');

    expect(descriptionDefinitions.length).toBe(9);

    const expectedRows = [
      ['Expected departure', '01 Jan 2020'],
      ['Origin address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 10001'],
      ['Destination address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 11111'],
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

  it('renders PPMShipmentCard with a heading that has a market code and shipment type', async () => {
    render(<PPMShipmentCard {...defaultProps} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(`${defaultProps.marketCode}PPM`);
  });

  it('omits the edit button when showEditAndDeleteBtn prop is false', () => {
    render(<PPMShipmentCard {...completeProps} showEditAndDeleteBtn={false} />);

    expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Delete' })).not.toBeInTheDocument();
  });

  it('calls onEditClick when edit button is pressed', async () => {
    render(<PPMShipmentCard {...completeProps} />);
    const editBtn = screen.getByRole('button', { name: 'Edit' });
    await userEvent.click(editBtn);
    expect(completeProps.onEditClick).toHaveBeenCalledTimes(1);
  });

  it('calls onDeleteClick when delete button is pressed', async () => {
    render(<PPMShipmentCard {...completeProps} />);
    const deleteBtn = screen.getByRole('button', { name: 'Delete' });
    await userEvent.click(deleteBtn);
    expect(completeProps.onDeleteClick).toHaveBeenCalledTimes(1);
  });

  it('renders component with closeout office and army affiliation', () => {
    const move = { closeoutOffice: transportationOfficeFactory() };

    render(<PPMShipmentCard {...completeProps} affiliation={affiliations.ARMY} move={move} />);

    const descriptionDefinitions = screen.getAllByRole('definition');

    expect(descriptionDefinitions.length).toBe(12);

    const expectedRows = [
      ['Expected departure', '01 Jan 2020'],
      ['Origin address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 10001'],
      ['Second origin address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 10002'],
      ['Destination address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 11111'],
      ['Second destination address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 22222'],
      ['Closeout office', move.closeoutOffice.name],
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

  it('renders component with closeout office and air force affiliation', () => {
    const move = { closeoutOffice: transportationOfficeFactory() };

    render(<PPMShipmentCard {...completeProps} affiliation={affiliations.AIR_FORCE} move={move} />);

    const descriptionDefinitions = screen.getAllByRole('definition');

    expect(descriptionDefinitions.length).toBe(12);

    const expectedRows = [
      ['Expected departure', '01 Jan 2020'],
      ['Origin address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 10001'],
      ['Second origin address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 10002'],
      ['Destination address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 11111'],
      ['Second destination address', '111 Test Street, 222 Test Street, Test Man, Test City, NY 22222'],
      ['Closeout office', move.closeoutOffice.name],
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

  it('does not render incomplete label and tooltip icon for completed ppm shipment with SUBMITTED status', async () => {
    render(<PPMShipmentCard {...completeProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('PPM 1');
    expect(screen.getByText(/^#testMove123-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.queryByText('Incomplete')).toBeNull();
  });

  it('renders incomplete label and tooltip icon for incomplete ppm shipment with DRAFT status', async () => {
    render(<PPMShipmentCard {...incompleteProps} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('PPM 1');
    expect(screen.getByText(/^#testMove123-01$/, { selector: 'p' })).toBeInTheDocument();

    expect(screen.getByText(/^Incomplete$/, { selector: 'span' })).toBeInTheDocument();

    expect(screen.getByTitle('Help about incomplete shipment')).toBeInTheDocument();
    await userEvent.click(screen.getByTitle('Help about incomplete shipment'));

    // verify onclick is getting json string as parameter
    expect(mockedOnIncompleteClickFunction).toHaveBeenCalledWith('PPM 1', 'testMove123-01', 'PPM');
  });

  it('renders incomplete PPMShipmentCard with a heading that has a market code and shipment type', async () => {
    render(<PPMShipmentCard {...incompleteProps} />);
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(`${incompleteProps.marketCode}PPM`);
  });
});
