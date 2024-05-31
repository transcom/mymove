import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';

import HeaderSection from './HeaderSection';

beforeEach(() => {
  jest.clearAllMocks();
});

const shipmentInfoProps = {
  sectionInfo: {
    type: 'shipmentInfo',
    plannedMoveDate: '2020-03-15',
    actualMoveDate: '2022-01-12',
    actualPickupPostalCode: '42444',
    actualDestinationPostalCode: '30813',
    miles: 513,
    estimatedWeight: 4000,
    actualWeight: 4200,
  },
};

const incentivesProps = {
  sectionInfo: {
    type: 'incentives',
    isAdvanceRequested: true,
    isAdvanceReceived: true,
    advanceAmountRequested: 598700,
    advanceAmountReceived: 112244,
    grossIncentive: 7231285,
    gcc: 7231285,
    remainingIncentive: 7119041,
  },
};

const incentiveFactorsProps = {
  sectionInfo: {
    type: 'incentiveFactors',
    haulType: 'Linehaul',
    haulPrice: 6892668,
    haulFSC: -143,
    packPrice: 20000,
    unpackPrice: 10000,
    dop: 15640,
    ddp: 34640,
    sitReimbursement: 30000,
  },
};

const incentiveFactorsShorthaulProps = {
  sectionInfo: {
    type: 'incentiveFactors',
    haulType: 'Shorthaul',
    haulPrice: 6892668,
    haulFSC: -143,
    packPrice: 20000,
    unpackPrice: 10000,
    dop: 15640,
    ddp: 34640,
    sitReimbursement: 30000,
  },
};

const invalidSectionTypeProps = {
  sectionInfo: {
    type: 'someUnknownSectionType',
  },
};

const clickDetailsButton = async (buttonType) => {
  fireEvent.click(screen.getByTestId(`${buttonType}-showRequestDetailsButton`));
  await waitFor(() => {
    expect(screen.getByText('Hide Details', { exact: false })).toBeInTheDocument();
  });
};

describe('PPMHeaderSummary component', () => {
  describe('displays Shipment Info section', () => {
    it('renders Shipment Info section on load with defaults', async () => {
      render(<HeaderSection {...shipmentInfoProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 4, name: 'Shipment Info' })).toBeInTheDocument();
      });

      clickDetailsButton('shipmentInfo');

      expect(screen.getByText('Planned Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('15-Mar-2020')).toBeInTheDocument();
      expect(screen.getByText('Actual Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('12-Jan-2022')).toBeInTheDocument();
      expect(screen.getByText('Starting ZIP')).toBeInTheDocument();
      expect(screen.getByText('42444')).toBeInTheDocument();
      expect(screen.getByText('Ending ZIP')).toBeInTheDocument();
      expect(screen.getByText('30813')).toBeInTheDocument();
      expect(screen.getByText('Miles')).toBeInTheDocument();
      expect(screen.getByText('513')).toBeInTheDocument();
      expect(screen.getByText('Estimated Net Weight')).toBeInTheDocument();
      expect(screen.getByText('4,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Actual Net Weight')).toBeInTheDocument();
      expect(screen.getByText('4,200 lbs')).toBeInTheDocument();
    });
  });

  describe('displays "Incentives/Costs" section', () => {
    it('renders "Incentives/Costs" section on load with correct prop values', async () => {
      render(<HeaderSection {...incentivesProps} />);

      clickDetailsButton('incentives');

      expect(screen.getByText('Government Constructed Cost (GCC)')).toBeInTheDocument();
      expect(screen.getByTestId('gcc')).toHaveTextContent('$72,312.85');
      expect(screen.getByText('Gross Incentive')).toBeInTheDocument();
      expect(screen.getByTestId('grossIncentive')).toHaveTextContent('$72,312.85');
      expect(screen.getByText('Advance Requested')).toBeInTheDocument();
      expect(screen.getByTestId('advanceRequested')).toHaveTextContent('$5,987.00');
      expect(screen.getByText('Advance Received')).toBeInTheDocument();
      expect(screen.getByTestId('advanceReceived')).toHaveTextContent('$1,122.44');
      expect(screen.getByText('Remaining Incentive')).toBeInTheDocument();
      expect(screen.getByTestId('remainingIncentive')).toHaveTextContent('$71,190.41');
    });
  });

  describe('displays "Incentive Factors" section', () => {
    it('renders "Incentive Factors" on load with correct prop values', async () => {
      render(<HeaderSection {...incentiveFactorsProps} />);

      clickDetailsButton('incentiveFactors');

      expect(screen.getByText('Linehaul Price')).toBeInTheDocument();
      expect(screen.getByTestId('haulPrice')).toHaveTextContent('$68,926.68');
      expect(screen.getByText('Linehaul Fuel Rate Adjustment')).toBeInTheDocument();
      expect(screen.getByTestId('haulFSC')).toHaveTextContent('-$1.43');
      expect(screen.getByText('Packing Charge')).toBeInTheDocument();
      expect(screen.getByTestId('packPrice')).toHaveTextContent('$200.00');
      expect(screen.getByText('Unpacking Charge')).toBeInTheDocument();
      expect(screen.getByTestId('unpackPrice')).toHaveTextContent('$100.00');
      expect(screen.getByText('Origin Price')).toBeInTheDocument();
      expect(screen.getByTestId('originPrice')).toHaveTextContent('$156.40');
      expect(screen.getByText('Destination Price')).toBeInTheDocument();
      expect(screen.getByTestId('destinationPrice')).toHaveTextContent('$346.40');
      expect(screen.getByTestId('sitReimbursement')).toHaveTextContent('$300.00');
    });

    it('renders "Shorthaul" in place of linehaul when given a shorthaul type', async () => {
      render(<HeaderSection {...incentiveFactorsShorthaulProps} />);

      clickDetailsButton('incentiveFactors');

      expect(screen.getByText('Shorthaul Price')).toBeInTheDocument();
      expect(screen.getByTestId('haulPrice')).toHaveTextContent('$68,926.68');
      expect(screen.getByText('Shorthaul Fuel Rate Adjustment')).toBeInTheDocument();
      expect(screen.getByTestId('haulFSC')).toHaveTextContent('-$1.43');
      expect(screen.getByText('Packing Charge')).toBeInTheDocument();
      expect(screen.getByTestId('packPrice')).toHaveTextContent('$200.00');
      expect(screen.getByText('Unpacking Charge')).toBeInTheDocument();
      expect(screen.getByTestId('unpackPrice')).toHaveTextContent('$100.00');
      expect(screen.getByText('Origin Price')).toBeInTheDocument();
      expect(screen.getByTestId('originPrice')).toHaveTextContent('$156.40');
      expect(screen.getByText('Destination Price')).toBeInTheDocument();
      expect(screen.getByTestId('destinationPrice')).toHaveTextContent('$346.40');
      expect(screen.getByTestId('sitReimbursement')).toHaveTextContent('$300.00');
    });
  });

  describe('handles errors correctly', () => {
    it('renders an alert if an unknown section type was passed in', async () => {
      render(<HeaderSection {...invalidSectionTypeProps} />);

      const alert = screen.getByTestId('alert');
      expect(alert).toBeInTheDocument();
      expect(alert).toHaveTextContent('Error getting section title!');
    });

    it('renders an alert if an unknown section type was passed in and details are expanded', async () => {
      render(<HeaderSection {...invalidSectionTypeProps} />);

      clickDetailsButton(invalidSectionTypeProps.sectionInfo.type);
      expect(screen.findByText('An error occured while getting section markup!'));
    });
  });
});
