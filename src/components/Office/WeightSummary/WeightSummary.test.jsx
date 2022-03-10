import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import { formatWeight } from 'utils/formatters';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';

const noOverweightShipments = [
  { id: '0001', shipmentType: 'HHG', calculatedBillableWeight: 6000, primeEstimatedWeight: 6000 },
  {
    id: '0002',
    shipmentType: 'HHG',
    calculatedBillableWeight: 400,
    primeEstimatedWeight: 4000,
    reweigh: { id: '1234' },
  },
  { id: '0003', shipmentType: 'HHG', calculatedBillableWeight: 3400, primeEstimatedWeight: 3400 },
];

const defaultProps = {
  maxBillableWeight: 15750,
  totalBillableWeight: 12460,
  weightRequested: 12260,
  weightAllowance: 8000,
  shipments: noOverweightShipments,
};

const overweightShipments = [
  { id: '0001', shipmentType: 'HHG', calculatedBillableWeight: 6161, primeEstimatedWeight: 3000 },
  {
    id: '0002',
    shipmentType: 'HHG',
    calculatedBillableWeight: 4000,
    primeEstimatedWeight: 4000,
    reweigh: { id: '1234' },
  },
  { id: '0003', shipmentType: 'HHG', calculatedBillableWeight: 3000, primeEstimatedWeight: 3000 },
];

const missingEstimatedWeightShipments = [{ id: '0001', shipmentType: 'HHG', calculatedBillableWeight: 6161 }];
const missingReweighWeightShipments = [
  {
    id: '0001',
    shipmentType: 'HHG',
    calculatedBillableWeight: 6161,
    reweigh: {
      dateReweighRequested: '2021-09-01',
    },
  },
];

const maxBillableWeightExceeded = {
  maxBillableWeight: 3750,
  totalBillableWeight: 12460,
  weightRequested: 12260,
  weightAllowance: 8000,
  shipments: overweightShipments,
};

const missingEstimatedWeight = {
  maxBillableWeight: 13750,
  totalBillableWeight: 12460,
  weightRequested: 12260,
  weightAllowance: 8000,
  shipments: missingEstimatedWeightShipments,
};

const missingReweighWeight = {
  maxBillableWeight: 13750,
  totalBillableWeight: 12460,
  weightRequested: 12260,
  weightAllowance: 8000,
  shipments: missingReweighWeightShipments,
};

describe('WeightSummary', () => {
  it('renders without crashing', () => {
    render(<WeightSummary {...defaultProps} />);
    // labels
    expect(screen.getByText('Max billable weight')).toBeInTheDocument();
    expect(screen.getByText('Weight requested')).toBeInTheDocument();
    expect(screen.getByText('Weight allowance')).toBeInTheDocument();
    expect(screen.getByText('Total billable weight')).toBeInTheDocument();

    // weights
    expect(screen.getByText(formatWeight(defaultProps.maxBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.totalBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightRequested))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();

    // shipment weights
    expect(screen.getByText(formatWeight(noOverweightShipments[0].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(noOverweightShipments[1].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(noOverweightShipments[2].calculatedBillableWeight))).toBeInTheDocument();
  });

  it('does not display flags when not appropriate', async () => {
    render(<WeightSummary {...defaultProps} />);

    await waitFor(() => {
      expect(screen.queryByTestId('shipmentHasFlag')).not.toBeInTheDocument();
    });
    expect(screen.queryByTestId('totalBillableWeightFlag')).not.toBeInTheDocument();
  });

  it('display max billable weight flag when appropriate', async () => {
    render(<WeightSummary {...maxBillableWeightExceeded} />);

    await waitFor(() => {
      expect(screen.queryByTestId('shipmentHasFlag')).toBeInTheDocument();
    });
    expect(screen.queryByTestId('totalBillableWeightFlag')).toBeInTheDocument();
  });

  it('display missing estimated weight flag when appropriate', async () => {
    render(<WeightSummary {...missingEstimatedWeight} />);

    await waitFor(() => {
      expect(screen.queryByTestId('shipmentHasFlag')).toBeInTheDocument();
    });
  });

  it('display missing reweigh weight flag when appropriate', async () => {
    render(<WeightSummary {...missingReweighWeight} />);

    await waitFor(() => {
      expect(screen.queryByTestId('shipmentHasFlag')).toBeInTheDocument();
    });
  });
});
