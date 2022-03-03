import React from 'react';
import { render, screen } from '@testing-library/react';

import EstimatedIncentiveDetails from './EstimatedIncentiveDetails';

const defaultProps = {
  shipment: {
    id: '1234',
    ppmShipment: {
      pickupPostalCode: '10001',
      destinationPostalCode: '10002',
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 3456,
      proGearWeight: 1333,
      proGearWeightSpouse: 425,
      estimatedIncentive: 876543,
    },
  },
};

const optionalSecondaryProps = {
  shipment: {
    id: '1234',
    ppmShipment: {
      pickupPostalCode: '10001',
      destinationPostalCode: '10002',
      secondaryPickupPostalCode: '10003',
      secondaryDestinationPostalCode: '10004',
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 3456,
      proGearWeight: 1333,
      proGearWeightSpouse: 425,
      estimatedIncentive: 876543,
    },
  },
};

describe('EstimatedIncentiveDetails component', () => {
  it('renders the details with required fields', async () => {
    render(<EstimatedIncentiveDetails {...defaultProps} />);
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('$8,765 is your estimated incentive');

    expect(
      screen.getByText("That's about how much you could earn for moving your PPM, based on what you've entered:"),
    ).toBeInTheDocument();

    const incentiveListItems = screen.getAllByRole('listitem');
    expect(incentiveListItems).toHaveLength(4);
    expect(incentiveListItems[0]).toHaveTextContent('3,456 lbs estimated weight');
    expect(incentiveListItems[1]).toHaveTextContent('Starting from 10001');
    expect(incentiveListItems[2]).toHaveTextContent('Ending in 10002');
    expect(incentiveListItems[3]).toHaveTextContent('Starting your PPM on 04 Jul 2022');

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Your actual incentive amount will vary');

    expect(
      screen.getByText(
        'Finance will determine your final incentive based on the total weight you move and the actual date you start moving your PPM.',
      ),
    ).toBeInTheDocument();

    expect(
      screen.getByText(
        'You must get certified weight tickets to document the weight you move. You are responsible for uploading them to MilMove.',
      ),
    );
  });

  it('conditionally renders secondary postal codes', () => {
    render(<EstimatedIncentiveDetails {...optionalSecondaryProps} />);

    const incentiveListItems = screen.getAllByRole('listitem');
    expect(incentiveListItems).toHaveLength(6);
    expect(incentiveListItems[0]).toHaveTextContent('3,456 lbs estimated weight');
    expect(incentiveListItems[1]).toHaveTextContent('Starting from 10001');
    expect(incentiveListItems[2]).toHaveTextContent('Picking up things in 10003');
    expect(incentiveListItems[3]).toHaveTextContent('Dropping off things in 10004');
    expect(incentiveListItems[4]).toHaveTextContent('Ending in 10002');
    expect(incentiveListItems[5]).toHaveTextContent('Starting your PPM on 04 Jul 2022');
  });
});
