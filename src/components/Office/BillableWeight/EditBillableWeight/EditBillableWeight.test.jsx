import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import EditBillableWeight from './EditBillableWeight';

import { formatWeight } from 'shared/formatters';

describe('EditBillableWeight', () => {
  it('renders edit button intially', () => {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
    };

    render(<EditBillableWeight {...defaultProps} />);

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
  });

  it('clicking edit button shows different view', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
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
      weightAllowance: 8000,
      estimatedWeight: 13750,
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
