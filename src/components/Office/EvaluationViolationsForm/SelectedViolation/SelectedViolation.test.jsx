/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import SelectedViolation from './SelectedViolation';

const mockViolation = {
  category: 'Category 1',
  displayOrder: 1,
  id: '9cdc8dc3-6cf4-46fb-b272-1468ef40796f',
  paragraphNumber: '1.2.3',
  requirementStatement: 'Test requirement statement for violation 1',
  requirementSummary: 'Test requirement summary for violation 1',
  subCategory: 'SubCategory 1',
  title: 'Title for violation 1',
};

const mockUnselectViolation = jest.fn();
beforeEach(() => {
  jest.clearAllMocks();
});

describe('SelectedViolation', () => {
  it('renders the violation content', async () => {
    render(<SelectedViolation violation={mockViolation} unselectViolation={mockUnselectViolation} />);

    expect(screen.getByRole('heading', { name: '1.2.3 Title for violation 1', level: 5 })).toBeInTheDocument();
    expect(screen.getByText(mockViolation.requirementSummary)).toBeInTheDocument();

    expect(screen.getByRole('button', { name: 'Remove' })).toBeInTheDocument();
    screen.debug();
  });

  it('removes the violation from selected when remove is clicked', async () => {
    render(<SelectedViolation violation={mockViolation} unselectViolation={mockUnselectViolation} />);

    expect(mockUnselectViolation).not.toHaveBeenCalled();

    userEvent.click(screen.getByRole('button', { name: 'Remove' }));

    expect(mockUnselectViolation).toHaveBeenCalledTimes(1);
    expect(mockUnselectViolation).toHaveBeenCalledWith(mockViolation.id);
  });
});
