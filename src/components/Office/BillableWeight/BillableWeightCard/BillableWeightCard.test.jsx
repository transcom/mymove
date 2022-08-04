import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import BillableWeightCard from './BillableWeightCard';

import { formatWeight } from 'utils/formatters';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

describe('BillableWeightCard', () => {
  const defaultProps = {
    maxBillableWeight: 13750,
    totalBillableWeight: 12460,
    weightRequested: 12260,
    weightAllowance: 8000,
    onReviewWeights: jest.fn(),
    secondaryReviewWeightsBtn: false,
  };

  const renderWithPermission = (props) => {
    render(
      <MockProviders permissions={[permissionTypes.updateMaxBillableWeight]}>
        <BillableWeightCard {...defaultProps} {...props} />
      </MockProviders>,
    );
  };

  it('renders maximum billable weight, total billable weight, weight requested and weight allowance', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 2161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234', weight: 40 },
      },
      {
        id: '0002',
        shipmentType: 'HHG',
        calculatedBillableWeight: 3200,
        estimatedWeight: 5000,
        primeEstimatedWeight: 1000,
        reweigh: { id: '1234', weight: 300 },
      },
      {
        id: '0003',
        shipmentType: 'HHG',
        calculatedBillableWeight: 3400,
        estimatedWeight: 5000,
        primeEstimatedWeight: 200,
        reweigh: { id: '1234', weight: 500 },
      },
    ];

    renderWithPermission({ shipments });

    // labels
    expect(screen.getByText('Maximum billable weight')).toBeInTheDocument();
    expect(screen.getByText('Weight requested')).toBeInTheDocument();
    expect(screen.getByText('Weight allowance')).toBeInTheDocument();
    expect(screen.getByText('Total billable weight')).toBeInTheDocument();

    // weights
    expect(screen.getByText(formatWeight(defaultProps.maxBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.totalBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightRequested))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();

    // shipment weights
    expect(screen.getByText(formatWeight(shipments[0].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[1].calculatedBillableWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(shipments[2].calculatedBillableWeight))).toBeInTheDocument();
  });

  it('implements the review weights handler when the review weights button is clicked', async () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234', weight: 40 },
      },
    ];
    renderWithPermission({ shipments });

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });

    await userEvent.click(reviewWeights);

    await waitFor(() => {
      expect(defaultProps.onReviewWeights).toHaveBeenCalled();
    });
  });

  it('displays secondary styling button when flag is set', async () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234', weight: 40 },
      },
    ];
    renderWithPermission({ shipments, secondaryReviewWeightsBtn: true });

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });
    expect(reviewWeights).toHaveClass('usa-button--secondary');
  });

  it('displays primary styling button when shipment has missing estimated weight', async () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        estimatedWeight: 5600,
        reweigh: { id: '1234', weight: 40 },
      },
    ];
    renderWithPermission({ shipments });

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });
    expect(reviewWeights).not.toHaveClass('usa-button--secondary');
  });

  it('displays primary styling button when shipment has missing reweigh weight', async () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        primeEstimatedWeight: 5800,
        reweigh: { id: '1234' },
      },
    ];
    render(
      <MockProviders permissions={[permissionTypes.updateMaxBillableWeight]}>
        <BillableWeightCard {...defaultProps} shipments={shipments} />
      </MockProviders>,
    );

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });
    expect(reviewWeights).not.toHaveClass('usa-button--secondary');
  });

  it('displays primary styling button when shipment has an overweight weight', async () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 60161,
        primeEstimatedWeight: 5800,
        reweigh: { id: '1234', weight: 2344 },
      },
    ];
    renderWithPermission({ shipments });

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });
    expect(reviewWeights).not.toHaveClass('usa-button--secondary');
  });

  it('displays primary styling button when the moves total weight exceeds the max billable weight', async () => {
    const props = {
      maxBillableWeight: 3750,
      totalBillableWeight: 12460,
      weightRequested: 12260,
      weightAllowance: 8000,
      onReviewWeights: jest.fn(),
      secondaryReviewWeightsBtn: false,
    };

    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        primeEstimatedWeight: 5800,
        reweigh: { id: '1234', weight: 2344 },
      },
    ];
    renderWithPermission({ ...props, shipments });

    const reviewWeights = screen.getByRole('button', { name: 'Review weights' });
    expect(reviewWeights).not.toHaveClass('usa-button--secondary');
  });

  it('renders external vendor weight summary with one NTSR external vendor shipment', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
        calculatedBillableWeight: 1000,
        estimatedWeight: 5600,
        ntsRecordedWeight: 1234,
        usesExternalVendor: true,
      },
      {
        id: '0002',
        shipmentType: 'HHG',
        calculatedBillableWeight: 3200,
        estimatedWeight: 5000,
        primeEstimatedWeight: 1000,
        reweigh: { id: '1234', weight: 300 },
      },
      {
        id: '0003',
        shipmentType: 'HHG',
        calculatedBillableWeight: 3400,
        estimatedWeight: 5000,
        primeEstimatedWeight: 200,
        reweigh: { id: '1234', weight: 500 },
      },
    ];
    renderWithPermission({ shipments });

    expect(screen.getByText('1 other shipment:')).toBeInTheDocument();
    expect(screen.getByText('1,234 lbs')).toBeInTheDocument();
  });

  it('renders external vendor weight summary with multiple external vendor NTSR shipments', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
        calculatedBillableWeight: 1000,
        estimatedWeight: 5600,
        ntsRecordedWeight: 4000,
        usesExternalVendor: true,
      },
      {
        id: '0002',
        shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
        calculatedBillableWeight: 3200,
        estimatedWeight: 5000,
        ntsRecordedWeight: 500,
        usesExternalVendor: true,
      },
      {
        id: '0003',
        shipmentType: 'HHG',
        calculatedBillableWeight: 3400,
        estimatedWeight: 5000,
        primeEstimatedWeight: 200,
        reweigh: { id: '1234', weight: 500 },
      },
    ];
    renderWithPermission({ shipments });

    expect(screen.getByText('2 other shipments:')).toBeInTheDocument();
    expect(screen.getByText('4,500 lbs')).toBeInTheDocument();
  });

  it('does not show the review weights button with no permission', () => {
    const shipments = [
      {
        id: '0001',
        shipmentType: 'HHG',
        calculatedBillableWeight: 6161,
        estimatedWeight: 5600,
        primeEstimatedWeight: 100,
        reweigh: { id: '1234', weight: 40 },
      },
    ];
    render(<BillableWeightCard {...defaultProps} shipments={shipments} />);

    const reviewWeights = screen.queryByRole('button', { name: 'Review weights' });
    expect(reviewWeights).not.toBeInTheDocument();
  });
});
