import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import { formatWeight } from 'shared/formatters';
import WeightSummary from 'components/Office/WeightSummary/WeightSummary';

const noOverweightShipments = [
  { id: '0001', shipmentType: 'HHG', billableWeightCap: 6000, primeEstimatedWeight: 6000 },
  {
    id: '0002',
    shipmentType: 'HHG',
    billableWeightCap: 400,
    primeEstimatedWeight: 4000,
    reweigh: { id: '1234' },
  },
  { id: '0003', shipmentType: 'HHG', billableWeightCap: 3400, primeEstimatedWeight: 3400 },
];

const defaultProps = {
  maxBillableWeight: 15750,
  totalBillableWeight: 12460,
  weightRequested: 12260,
  weightAllowance: 8000,
  shipments: noOverweightShipments,
};

const overweightShipments = [
  { id: '0001', shipmentType: 'HHG', billableWeightCap: 6161, primeEstimatedWeight: 3000 },
  {
    id: '0002',
    shipmentType: 'HHG',
    billableWeightCap: 4000,
    primeEstimatedWeight: 4000,
    reweigh: { id: '1234' },
  },
  { id: '0003', shipmentType: 'HHG', billableWeightCap: 3000, primeEstimatedWeight: 3000 },
];

const maxBillableWeightExceeded = {
  maxBillableWeight: 3750,
  totalBillableWeight: 12460,
  weightRequested: 12260,
  weightAllowance: 8000,
  shipments: overweightShipments,
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
    expect(screen.getByText(formatWeight(noOverweightShipments[0].billableWeightCap))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(noOverweightShipments[1].billableWeightCap))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(noOverweightShipments[2].billableWeightCap))).toBeInTheDocument();
  });

  it('does not display flags when not appropriate', async () => {
    render(<WeightSummary {...defaultProps} />);

    await waitFor(() => {
      expect(screen.queryByTestId('shipmentIsOverweightFlag')).not.toBeInTheDocument();
    });
    expect(screen.queryByTestId('totalBillableWeightFlag')).not.toBeInTheDocument();
  });

  it('display flags when appropriate', async () => {
    render(<WeightSummary {...maxBillableWeightExceeded} />);

    await waitFor(() => {
      expect(screen.queryByTestId('shipmentIsOverweightFlag')).toBeInTheDocument();
    });
    expect(screen.queryByTestId('totalBillableWeightFlag')).toBeInTheDocument();
  });
});
