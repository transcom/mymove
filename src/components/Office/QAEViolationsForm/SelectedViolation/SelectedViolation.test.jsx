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
    render(
      <SelectedViolation violation={mockViolation} unselectViolation={mockUnselectViolation} isReadOnly={false} />,
    );

    expect(screen.getByRole('heading', { name: '1.2.3 Title for violation 1', level: 5 })).toBeInTheDocument();
    expect(screen.getByText(mockViolation.requirementSummary)).toBeInTheDocument();

    // remove button should display when read only is false
    expect(screen.getByRole('button', { name: 'Remove' })).toBeInTheDocument();
  });

  it('removes the violation from selected when remove is clicked', async () => {
    render(<SelectedViolation violation={mockViolation} unselectViolation={mockUnselectViolation} />);

    expect(mockUnselectViolation).not.toHaveBeenCalled();

    await userEvent.click(screen.getByRole('button', { name: 'Remove' }));

    expect(mockUnselectViolation).toHaveBeenCalledTimes(1);
    expect(mockUnselectViolation).toHaveBeenCalledWith(mockViolation.id);
  });

  it('does not render anything when there is no violation', async () => {
    render(<SelectedViolation violation={null} unselectViolation={mockUnselectViolation} />);

    expect(screen.queryByRole('heading', { name: '1.2.3 Title for violation 1', level: 5 })).not.toBeInTheDocument();
    expect(screen.queryByText(mockViolation.requirementSummary)).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Remove' })).not.toBeInTheDocument();
  });
});
