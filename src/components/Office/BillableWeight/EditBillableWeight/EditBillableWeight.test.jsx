import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import EditBillableWeight from './EditBillableWeight';

import { formatWeight } from 'shared/formatters';

describe('EditBillableWeight', () => {
  it('renders edit button intially', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
    };

    render(<EditBillableWeight {...defaultProps} />);

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
  });

  it('renders max billable weight view', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
    };

    render(<EditBillableWeight {...defaultProps} />);
    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.getByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.estimatedWeight * 1.1))).toBeInTheDocument();
    expect(screen.getByText('| weight allowance')).toBeInTheDocument();
    expect(screen.getByText('| 110% of total estimated weight')).toBeInTheDocument();
  });

  it('renders edit billable weight view', () => {
    const defaultProps = {
      title: 'Billable weight',
      originalWeight: 10000,
      estimatedWeight: 13000,
      maxBillableWeight: 6000,
      billableWeight: 7000,
      totalBillableWeight: 11000,
    };

    render(<EditBillableWeight {...defaultProps} />);
    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.getByText(formatWeight(defaultProps.originalWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.estimatedWeight * 1.1))).toBeInTheDocument();
    expect(
      screen.getByText(formatWeight(defaultProps.totalBillableWeight - defaultProps.billableWeight)),
    ).toBeInTheDocument();
    expect(screen.getByText('| original weight')).toBeInTheDocument();
    expect(screen.getByText('| 110% of total estimated weight')).toBeInTheDocument();
    expect(screen.getByText('| to fit within max billable weight')).toBeInTheDocument();
  });

  it('clicking edit button shows different view', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
    };

    render(<EditBillableWeight {...defaultProps} />);

    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.queryByText('Edit')).toBeNull();
    // weights
    expect(screen.getByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.estimatedWeight * 1.1))).toBeInTheDocument();
    // buttons
    expect(screen.getByRole('button', { name: 'Save changes' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
  });

  it('should be able to toggle between views', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
    };

    render(<EditBillableWeight {...defaultProps} />);
    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.queryByText('Edit')).toBeNull();
    expect(screen.getByRole('button', { name: 'Save changes' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(screen.queryByText('Edit')).toBeInTheDocument();
    expect(screen.queryByText('Save changes')).toBeNull();
    expect(screen.queryByText('Cancel')).toBeNull();
  });
});
