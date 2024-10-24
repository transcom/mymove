import React from 'react';
import { render, screen, } from '@testing-library/react';
import QueueTable from '.';

describe('QueueTable', () => {
  it('renders without crashing', () => {
    render(<QueueTable />);
    expect(screen.getByText('Payment requested')).toBeInTheDocument();
  });
});
