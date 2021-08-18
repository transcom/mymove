import React from 'react';
import { render, screen } from '@testing-library/react';

import ReviewBillableWeight from './ReviewBillableWeight';

describe('ReviewBillableWeight', () => {
  it('renders the component', () => {
    render(<ReviewBillableWeight />);
    expect(screen.getByText('Review Billable Weight page')).toBeInTheDocument();
  });
});
