import React from 'react';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import LeftNav from './LeftNav';

import LeftNavTag from 'components/LeftNavTag/LeftNavTag';

describe('Left Nav Component', () => {
  const sections = ['orders', 'allowances', 'requested-shipments'];

  it('renders without errors', () => {
    render(<LeftNav sections={sections} />);

    expect(screen.getByText('Orders')).toBeInTheDocument();
    expect(screen.getByText('Allowances')).toBeInTheDocument();
    expect(screen.getByText('Requested shipments')).toBeInTheDocument();
  });

  it('activates a section when the section is clicked', () => {
    render(<LeftNav sections={sections} />);

    const sectionLabel = screen.getByText('Orders');

    userEvent.click(sectionLabel);

    expect(sectionLabel).toHaveClass('active');
  });

  it('correctly pairs left nav tags to their associated section', () => {
    render(
      <LeftNav sections={sections}>
        <LeftNavTag associatedSectionName="orders" showTag>
          Test Tag
        </LeftNavTag>
      </LeftNav>,
    );

    const section = screen.getByText('Orders');
    expect(within(section).getByText('Test Tag')).toBeInTheDocument();
  });
});
