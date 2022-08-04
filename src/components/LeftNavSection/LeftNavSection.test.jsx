import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import LeftNavSection from './LeftNavSection';

describe('Left Nav Section Component', () => {
  it('renders without errors', () => {
    render(<LeftNavSection sectionName="testSection">Test Section</LeftNavSection>);
    expect(screen.getByText('Test Section')).toBeInTheDocument();
  });

  it('uses the onClickHandler when provided', async () => {
    const mockOnClick = jest.fn();
    render(
      <LeftNavSection sectionName="testSection" onClickHandler={mockOnClick}>
        Test Section
      </LeftNavSection>,
    );

    const sectionLabel = screen.getByText('Test Section');

    await userEvent.click(sectionLabel);

    expect(mockOnClick).toHaveBeenCalledTimes(1);
  });

  it('uses the active property when it the link is active', () => {
    render(
      <LeftNavSection sectionName="testSection" isActive>
        Test Section
      </LeftNavSection>,
    );

    const sectionLabel = screen.getByText('Test Section');

    expect(sectionLabel).toHaveClass('active');
  });
});
