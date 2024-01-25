import React from 'react';
import { render, screen } from '@testing-library/react';

import '@testing-library/jest-dom/extend-expect'; // For additional matchers like toBeInTheDocument

import MultiMovesLandingPage from './MultiMovesLandingPage';

import { MockProviders } from 'testUtils';

// Mock external dependencies
jest.mock('utils/featureFlags', () => ({
  detectFlags: jest.fn(() => ({ multiMove: true })),
}));

jest.mock('store/auth/actions', () => ({
  loadUser: jest.fn(),
}));

jest.mock('store/onboarding/actions', () => ({
  initOnboarding: jest.fn(),
}));

jest.mock('shared/Swagger/ducks', () => ({
  loadInternalSchema: jest.fn(),
}));

describe('MultiMovesLandingPage', () => {
  it('renders the component with moves', () => {
    render(
      <MockProviders>
        <MultiMovesLandingPage />
      </MockProviders>,
    );

    // Check for specific elements
    expect(screen.getByTestId('customerHeader')).toBeInTheDocument();
    expect(screen.getByTestId('helperText')).toBeInTheDocument();
    expect(screen.getByText('First Last')).toBeInTheDocument();
    expect(screen.getByText('Welcome to MilMove!')).toBeInTheDocument();
    expect(screen.getByText('Create a Move')).toBeInTheDocument();

    // Assuming there are two move headers and corresponding move containers
    expect(screen.getAllByText('Current Move')).toHaveLength(1);
    expect(screen.getAllByText('Previous Moves')).toHaveLength(1);
  });

  it('renders move data correctly', () => {
    render(
      <MockProviders>
        <MultiMovesLandingPage />
      </MockProviders>,
    );

    expect(screen.getByTestId('currentMoveHeader')).toBeInTheDocument();
    expect(screen.getByTestId('currentMoveContainer')).toBeInTheDocument();
    expect(screen.getByTestId('prevMovesHeader')).toBeInTheDocument();
    expect(screen.getByTestId('prevMovesContainer')).toBeInTheDocument();
  });
});
