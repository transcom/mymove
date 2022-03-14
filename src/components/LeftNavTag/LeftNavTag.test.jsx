import React from 'react';
import { render, screen } from '@testing-library/react';

import LeftNavTag from './LeftNavTag';

describe('Left Nav Tag Component', () => {
  it('renders without errors', () => {
    render(<LeftNavTag showTag>Test Tag</LeftNavTag>);
    expect(screen.getByText('Test Tag')).toBeInTheDocument();
  });

  it('does not render if showTag is false', () => {
    render(<LeftNavTag showTag={false}>Test Tag</LeftNavTag>);
    expect(screen.queryByText('Test Tag')).not.toBeInTheDocument();
  });
});
