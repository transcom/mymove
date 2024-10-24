import React from 'react';
import { render, screen } from '@testing-library/react';
import QueueTable from '.';

describe('QueueTable', () => {
  it('renders without crashing', () => {
    render(<QueueTable />);
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText('Confirmation')).toBeInTheDocument();
    expect(screen.getByText('Branch')).toBeInTheDocument();
    expect(screen.getByText('Original duty location')).toBeInTheDocument();
    expect(screen.getByText('Last modified by')).toBeInTheDocument();
  });
});
