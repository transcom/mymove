import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentCard from './ShipmentCard';

import { formatAddressShort, formatDateFromIso } from 'shared/formatters';
import { formatWeight } from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const tomorrow = new Date();
tomorrow.setDate(tomorrow.getDate() + 1);

const defaultShipmentCardProps = {
  billableWeight: 2000,
  maxBillableWeight: 0,
  dateReweighRequested: new Date().toISOString(),
  departedDate: tomorrow.toISOString(),
  pickupAddress: {
    city: 'Rancho Santa Margarita',
    state: 'CA',
    postalCode: '92688',
  },
  destinationAddress: {
    city: 'West Springfield Town',
    state: 'MA',
    postalCode: '01089',
  },
  estimatedWeight: 5000,
  originalWeight: 4999,
  reweighWeight: 4999,
  adjustedWeight: null,
  reweighRemarks: 'Unable to perform reweigh because shipment was already unloaded',
  editEntity: () => {},
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

describe('ShipmentCard', () => {
  it('renders address and weight information', () => {
    const defaultProps = {
      billableWeight: 4014,
      maxBillableWeight: 0,
      dateReweighRequested: new Date().toISOString(),
      departedDate: tomorrow.toISOString(),
      pickupAddress: {
        city: 'Rancho Santa Margarita',
        state: 'CA',
        postalCode: '92688',
      },
      destinationAddress: {
        city: 'West Springfield Town',
        state: 'MA',
        postalCode: '01089',
      },
      estimatedWeight: 5000,
      originalWeight: 4300,
      reweighRemarks: 'Unable to perform reweigh because shipment was already unloaded',
      editEntity: () => {},
      shipmentType: SHIPMENT_OPTIONS.HHG,
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

  describe('warning indicator', () => {
    it('renders no yellow highlight for original and reweigh weights if weights do not exceeds 110% estimated weight', () => {
      const defaultProps = {
        ...defaultShipmentCardProps,
        originalWeight: 4999,
        reweighWeight: 4999,
        adjustedWeight: null,
      };

      render(<ShipmentCard {...defaultProps} />);

      expect(screen.queryByTestId('originalWeightContainer')).not.toHaveClass('warning');
      expect(screen.queryByTestId('reweighWeightContainer')).not.toHaveClass('warning');
    });

    it('renders no yellow highlight for original and reweigh weights if adjusted weight is set', () => {
      const defaultProps = {
        ...defaultShipmentCardProps,
        estimatedWeight: 5000,
        originalWeight: 3000,
        reweighWeight: 2000,
        adjustedWeight: 1000,
      };

      render(<ShipmentCard {...defaultProps} />);

      expect(screen.queryByTestId('originalWeightContainer')).not.toHaveClass('warning');
      expect(screen.queryByTestId('reweighWeightContainer')).not.toHaveClass('warning');
    });

    it('renders yellow highlight for original weight that exceeds 110% estimated weight', () => {
      const defaultProps = {
        ...defaultShipmentCardProps,
        estimatedWeight: 5000,
        originalWeight: 5510,
        reweighWeight: 6000,
        adjustedWeight: null,
      };

      render(<ShipmentCard {...defaultProps} />);

      expect(screen.getByTestId('originalWeightContainer')).toHaveClass('warning');
      expect(screen.queryByTestId('reweighWeightContainer')).not.toHaveClass('warning');
    });

    it('renders yellow highlight for original weight that exceeds 110% estimated weight and reweigh weight missing', () => {
      const defaultProps = {
        ...defaultShipmentCardProps,
        estimatedWeight: 5000,
        originalWeight: 6000,
        reweighWeight: null,
        adjustedWeight: null,
      };

      render(<ShipmentCard {...defaultProps} />);

      expect(screen.getByTestId('originalWeightContainer')).toHaveClass('warning');
      expect(screen.getByTestId('reweighWeightContainer')).toHaveClass('warning');
    });

    it('renders yellow highlight for reweigh weight that exceeds 110% estimated weight', () => {
      const defaultProps = {
        ...defaultShipmentCardProps,
        estimatedWeight: 5000,
        originalWeight: 6000,
        reweighWeight: 5510,
        adjustedWeight: null,
      };

      render(<ShipmentCard {...defaultProps} />);

      expect(screen.queryByTestId('originalWeightContainer')).not.toHaveClass('warning');
      expect(screen.getByTestId('reweighWeightContainer')).toHaveClass('warning');
    });

    it('renders yellow highlight for missing estimated weight', () => {
      const defaultProps = {
        ...defaultShipmentCardProps,
        estimatedWeight: null,
        originalWeight: 6000,
        reweighWeight: 5510,
        adjustedWeight: null,
      };

      render(<ShipmentCard {...defaultProps} />);

      expect(screen.getByTestId('estimatedWeightContainer')).toHaveClass('warning');
      expect(screen.queryByTestId('originalWeightContainer')).not.toHaveClass('warning');
      expect(screen.queryByTestId('reweighWeightContainer')).not.toHaveClass('warning');
    });
  });

  it('does not render the reweigh remarks if there are no reweigh remarks', () => {
    const defaultProps = {
      billableWeight: 4014,
      maxBillableWeight: 0,
      dateReweighRequested: new Date().toISOString(),
      departedDate: tomorrow.toISOString(),
      pickupAddress: {
        city: 'Rancho Santa Margarita',
        state: 'CA',
        postalCode: '92688',
      },
      destinationAddress: {
        city: 'West Springfield Town',
        state: 'MA',
        postalCode: '01089',
      },
      estimatedWeight: 5000,
      originalWeight: 4300,
      editEntity: () => {},
      shipmentType: SHIPMENT_OPTIONS.HHG,
    };

    render(<ShipmentCard {...defaultProps} />);
    // labels
    expect(screen.queryByText('Reweigh remarks')).not.toBeInTheDocument();
  });
});
