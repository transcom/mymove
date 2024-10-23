import React from 'react';
import { render, screen } from '@testing-library/react';

import EstimatedIncentiveDetails from 'components/Customer/PPM/Booking/EstimatedIncentiveDetails/EstimatedIncentiveDetails';

const defaultProps = {
  shipment: {
    id: '1234',
    ppmShipment: {
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '813 S 129th St',
        streetAddress2: '#124',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 3456,
      proGearWeight: 1333,
      proGearWeightSpouse: 425,
      estimatedIncentive: 876543,
    },
  },
};

const zeroIncentiveProps = {
  shipment: {
    id: '1234',
    ppmShipment: {
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '813 S 129th St',
        streetAddress2: '#124',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
      expectedDepartureDate: '2022-07-04',
      estimatedWeight: 3456,
      proGearWeight: 1333,
      proGearWeightSpouse: 425,
      estimatedIncentive: 0,
    },
  },
};

const optionalSecondaryProps = {
  shipment: {
    id: '1234',
    ppmShipment: {
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10001',
      },
      destinationAddress: {
        streetAddress1: '813 S 129th St',
        streetAddress2: '#124',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10002',
      },
      secondaryPickupAddress: {
        streetAddress1: '813 S 129th St',
        streetAddress2: '#125',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10003',
      },
      secondaryDestinationAddress: {
        streetAddress1: '814 S 129th St',
        streetAddress2: '#126',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '10004',
      },
      hasSecondaryPickupAddress: true,
      hasSecondaryDestinationAddress: true,
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
      screen.getByText(
        'This is an estimate of how much you could earn by moving your PPM, based on what you have entered:',
      ),
    ).toBeInTheDocument();

    const incentiveListItems = screen.getAllByRole('listitem');
    expect(incentiveListItems).toHaveLength(4);
    expect(incentiveListItems[0]).toHaveTextContent('3,456 lbs estimated weight');
    expect(incentiveListItems[1]).toHaveTextContent('Starting from 812 S 129th St, #123, San Antonio, TX 10001');
    expect(incentiveListItems[2]).toHaveTextContent('Ending at 813 S 129th St, #124, San Antonio, TX 10002');
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

  it('renders the DTOD unavailable when estimated incentive is', async () => {
    render(<EstimatedIncentiveDetails {...zeroIncentiveProps} />);

    expect(
      screen.getByText('The Defense Table of Distances (DTOD) was unavailable during your PPM creation, so we are currently unable to provide your estimated incentive. Your estimated incentive information will be updated and provided to you during your counseling session.'),
    ).toBeInTheDocument();
  });

  it('conditionally renders secondary postal codes', () => {
    render(<EstimatedIncentiveDetails {...optionalSecondaryProps} />);

    const incentiveListItems = screen.getAllByRole('listitem');
    expect(incentiveListItems).toHaveLength(6);
    expect(incentiveListItems[0]).toHaveTextContent('3,456 lbs estimated weight');
    expect(incentiveListItems[1]).toHaveTextContent('Starting from 812 S 129th St, #123, San Antonio, TX 10001');
    expect(incentiveListItems[2]).toHaveTextContent('Picking up things at 813 S 129th St, #125, San Antonio, TX 10003');
    expect(incentiveListItems[3]).toHaveTextContent(
      'Dropping off things at 814 S 129th St, #126, San Antonio, TX 10004',
    );
    expect(incentiveListItems[4]).toHaveTextContent('Ending at 813 S 129th St, #124, San Antonio, TX 10002');
    expect(incentiveListItems[5]).toHaveTextContent('Starting your PPM on 04 Jul 2022');
  });
});
