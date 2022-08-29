import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ViolationsAccordion from './ViolationsAccordion';

const mockPwsViolations = [
  {
    category: 'Category 1',
    displayOrder: 1,
    id: '9cdc8dc3-6cf4-46fb-b272-1468ef40796f',
    paragraphNumber: '1.2.3',
    requirementStatement: 'Test requirement statement for violation 1',
    requirementSummary: 'Test requirement summary for violation 1',
    subCategory: 'SubCategory 1',
    title: 'Title for violation 1',
  },
  {
    category: 'Category 1',
    displayOrder: 2,
    id: 'c359ebc3-a506-4f41-8f91-409d59c97b22',
    paragraphNumber: '1.2.5.1',
    requirementStatement: 'Test requirement statement for violation 2',
    requirementSummary: 'Test requirement summary for violation 2',
    subCategory: 'SubCategory 2',
    title: 'Title for violation 2',
  },
  {
    category: 'Category 1',
    displayOrder: 3,
    id: 'd359ebc3-a506-4f41-8f91-409d59c97b23',
    paragraphNumber: '1.2.5.1',
    requirementStatement: 'Test requirement statement for violation 3',
    requirementSummary: 'Test requirement summary for violation 3',
    subCategory: 'SubCategory 3',
    title: 'Title for violation 3',
  },
];
const mockViolation = mockPwsViolations[0];

describe('ViolationsAccordion', () => {
  it('can render successfully', () => {
    render(<ViolationsAccordion violations={[mockViolation]} />);

    // Check that intitially displayed content is rendered
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(mockViolation.category);
    expect(screen.getByRole('button', { name: mockViolation.subCategory })).toBeInTheDocument();
    expect(screen.getByText(`${mockViolation.paragraphNumber} ${mockViolation.title}`)).toBeInTheDocument();
    expect(screen.getByText(mockViolation.requirementSummary)).toBeInTheDocument();
  });

  it('can expand and collapse the accordion to view requirement statement', () => {
    render(<ViolationsAccordion violations={[mockViolation]} />);

    // Requiremnt statement not shown initially, requires user to expand to view
    expect(screen.queryByText(mockViolation.requirementStatement)).not.toBeInTheDocument();

    // Should have expand but not collapse button
    expect(screen.getByTestId('expand-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('collapse-icon')).not.toBeInTheDocument();

    // Click expand
    userEvent.click(screen.getByTestId('expand-icon'));

    // Should now be showning the requirement statement
    expect(screen.getByText(mockViolation.requirementStatement)).toBeInTheDocument();

    // Expand icon should have become collapse icon
    expect(screen.getByTestId('collapse-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('expand-icon')).not.toBeInTheDocument();

    // Collapse section
    userEvent.click(screen.getByTestId('collapse-icon'));

    // Requirments statement should no longer be shown and buttons reverted back
    expect(screen.queryByText(mockViolation.requirementStatement)).not.toBeInTheDocument();
    expect(screen.getByTestId('expand-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('collapse-icon')).not.toBeInTheDocument();
  });

  it('groups violations by subcategory', () => {
    render(<ViolationsAccordion violations={mockPwsViolations} />);

    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('Category 1');

    // Should have all 3 subcategories
    expect(screen.getByRole('button', { name: 'SubCategory 1' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'SubCategory 2' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'SubCategory 3' })).toBeInTheDocument();
  });
});
