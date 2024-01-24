import React from 'react';
import { render, screen } from '@testing-library/react';
import { v4 } from 'uuid';

import '@testing-library/jest-dom/extend-expect';

import MultiMovesLandingPage from './MultiMovesLandingPage';

import { MockProviders } from 'testUtils';
import { MOVE_STATUSES } from 'shared/constants';

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

const defaultProps = {
  serviceMember: {
    id: v4(),
    current_location: {
      transportation_office: {
        name: 'Test Transportation Office Name',
        phone_lines: ['555-555-5555'],
      },
    },
    weight_allotment: {
      total_weight_self: 8000,
      total_weight_self_plus_dependents: 11000,
    },
  },
  showLoggedInUser: jest.fn(),
  createServiceMember: jest.fn(),
  getSignedCertification: jest.fn(),
  mtoShipments: [],
  mtoShipment: {},
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  loadMTOShipments: jest.fn(),
  updateShipmentList: jest.fn(),
  move: {
    id: v4(),
    status: MOVE_STATUSES.DRAFT,
  },
  uploadedOrderDocuments: [],
  uploadedAmendedOrderDocuments: [],
};

describe('MultiMovesLandingPage', () => {
  it('renders the component with moves', () => {
    render(
      <MockProviders>
        <MultiMovesLandingPage {...defaultProps} />
      </MockProviders>,
    );

    // Check for specific elements
    expect(screen.getByTestId('customerHeader')).toBeInTheDocument();
    expect(screen.getByTestId('welcomeHeader')).toBeInTheDocument();
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
