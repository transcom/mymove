import React from 'react';
import { render, screen } from '@testing-library/react';

import CreateServiceItem from './CreateServiceItem';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createServiceItem: jest.fn(),
}));

describe('CreateServiceItem page', () => {
  describe('check loading and error component states', () => {
    const loadingReturnValue = {
      moveTaskOrder: undefined,
      isLoading: true,
      isError: false,
    };

    const errorReturnValue = {
      moveTaskOrder: undefined,
      isLoading: false,
      isError: true,
    };

    it('renders the loading placeholder when the query is still loading', () => {
      usePrimeSimulatorGetMove.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders>
          <CreateServiceItem setFlashMessage={jest.fn()} />
        </MockProviders>,
      );

      expect(screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 }));
    });

    it('renders the Something Went Wrong component when the query has an error', () => {
      usePrimeSimulatorGetMove.mockReturnValue(errorReturnValue);

      render(
        <MockProviders>
          <CreateServiceItem setFlashMessage={jest.fn()} />
        </MockProviders>,
      );

      expect(screen.getByText(/Something went wrong./));
    });
  });
});
