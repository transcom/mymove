import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';

import { SHIPMENT_OPTIONS } from '../../../../shared/constants';

import ReviewShipmentWeightsTable from './ReviewShipmentWeightsTable';
import { nonPPMReviewWeightsTableConfig, PPMReviewWeightsTableConfig } from './helpers';

import { MockProviders } from 'testUtils';

beforeEach(() => {
  jest.clearAllMocks();
});

const PPMProps = {
  tableData: [
    {
      shipmentType: 'PPM',
      ppmShipment: {
        hasReceivedAdvance: true,
        advanceAmountReceived: 60000,
        estimatedWeight: 4000,
        actualMoveDate: '2023-04-01',
        weightTickets: [
          {
            emptyWeight: 1000,
            fullWeight: 6001,
          },
        ],
      },
      actualProGearWeight: 1000,
      actualSpouseProGearWeight: 500,
    },
  ],
  tableConfig: PPMReviewWeightsTableConfig,
};
const NonPPMProps = {
  tableData: [
    {
      shipmentType: SHIPMENT_OPTIONS.HHG,
      primeEstimatedWeight: 2500,
      calculatedBillableWeight: 3000,
      primeActualWeight: 3500,
      reweigh: {
        id: 'rw01',
        weight: 3200,
      },
      actualDeliveryDate: '2023-04-23',
    },
  ],
  tableConfig: nonPPMReviewWeightsTableConfig,
};

describe('ReviewShipmentWeight component', () => {
  it('correctly renders PPM table data', async () => {
    render(
      <MockProviders>
        <ReviewShipmentWeightsTable {...PPMProps} />
      </MockProviders>,
    );
    await waitFor(async () => {
      expect(screen.getByTestId('reviewShipmentWeightsTable')).toBeInTheDocument();
      expect(screen.getByText('Weight ticket')).toBeInTheDocument();
      expect(screen.getByText('Review Documents')).toBeInTheDocument();
      expect(screen.getByText('Review Documents').tagName).toBe('A');
      expect(screen.getByText('Pro-gear (lbs)')).toBeInTheDocument();
      expect(screen.getByText('4,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Spouse pro-gear')).toBeInTheDocument();
      expect(screen.getByText('500 lbs')).toBeInTheDocument();
      expect(screen.getByText('Net weight')).toBeInTheDocument();
      expect(screen.getByText('5,001 lbs')).toBeInTheDocument();
      expect(screen.getByText('Actual Departure date')).toBeInTheDocument();
      expect(screen.getByText('Apr 01 2023')).toBeInTheDocument();
    });
  });
  it('correctly renders non-PPM table data', async () => {
    render(
      <MockProviders>
        <ReviewShipmentWeightsTable {...NonPPMProps} />
      </MockProviders>,
    );
    await waitFor(async () => {
      expect(screen.getByTestId('reviewShipmentWeightsTable')).toBeInTheDocument();
      expect(screen.getByText('Estimated weight')).toBeInTheDocument();
      expect(screen.getByText('2,500 lbs')).toBeInTheDocument();
      expect(screen.getByText('Reweigh requested')).toBeInTheDocument();
      expect(screen.getByText('Yes')).toBeInTheDocument();
      expect(screen.getByText('Billable weight')).toBeInTheDocument();
      expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Delivery date')).toBeInTheDocument();
      expect(screen.getByText('Apr 23 2023')).toBeInTheDocument();
    });
  });
});
