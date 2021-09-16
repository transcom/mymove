import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentCard from './ShipmentCard';

import { formatWeight, formatAddressShort, formatDateFromIso } from 'shared/formatters';

describe('ShipmentCard', () => {
  it('renders address and weight information', () => {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const defaultProps = {
      billableWeight: 4014,
      maxBillableWeight: 0,
      dateReweighRequested: new Date().toISOString(),
      departedDate: tomorrow.toISOString(),
      pickupAddress: {
        city: 'Rancho Santa Margarita',
        state: 'CA',
        postal_code: '92688',
      },
      destinationAddress: {
        city: 'West Springfield Town',
        state: 'MA',
        postal_code: '01089',
      },
      estimatedWeight: 5000,
      originalWeight: 4300,
      reweighRemarks: 'Unable to perform reweigh because shipment was already unloaded',
    };

    render(<ShipmentCard {...defaultProps} />);
    // labels
    expect(screen.getByText('Departed')).toBeInTheDocument();
    expect(screen.getByText('From')).toBeInTheDocument();
    expect(screen.getByText('To')).toBeInTheDocument();
    expect(screen.getByText('Estimated weight')).toBeInTheDocument();
    expect(screen.getByText('Original weight')).toBeInTheDocument();
    expect(screen.getByText('Reweigh weight')).toBeInTheDocument();
    expect(screen.getByText('Date reweigh requested')).toBeInTheDocument();
    expect(screen.getByText('Reweigh remarks')).toBeInTheDocument();
    expect(screen.getByText('Billable weight')).toBeInTheDocument();

    //  weights
    expect(screen.getByText(formatWeight(defaultProps.billableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.estimatedWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.originalWeight))).toBeInTheDocument();

    // dates
    expect(screen.getByText(formatDateFromIso(defaultProps.dateReweighRequested, 'DD MMM YYYY'))).toBeInTheDocument();
    expect(screen.getByText(formatDateFromIso(defaultProps.departedDate, 'DD MMM YYYY'))).toBeInTheDocument();

    // addresses
    expect(screen.getByText(formatAddressShort(defaultProps.pickupAddress))).toBeInTheDocument();
    expect(screen.getByText(formatAddressShort(defaultProps.destinationAddress))).toBeInTheDocument();
  });
});
