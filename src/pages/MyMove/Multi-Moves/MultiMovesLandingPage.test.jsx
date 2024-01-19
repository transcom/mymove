import React from 'react';
import { render, screen } from '@testing-library/react';

import MultiMovesLandingPage from './MultiMovesLandingPage';

describe('MultiMovesLandingPage', () => {
  it('renders the MultiMovesLandingPage component with retirement moves', () => {
    render(<MultiMovesLandingPage />);

    expect(screen.getByTestId('customer-header')).toBeInTheDocument();
    expect(screen.getByText('First Last')).toBeInTheDocument();
    expect(screen.getByText('Welcome to MilMove!')).toBeInTheDocument();
    expect(screen.getByText('Create a Move')).toBeInTheDocument();

    // Assuming there are two move headers and corresponding move containers
    expect(screen.getAllByText('Current Move')).toHaveLength(1);
    expect(screen.getAllByText('Previous Moves')).toHaveLength(1);
  });
});
