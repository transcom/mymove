import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import QAECSRMoveSearch from './QAECSRMoveSearch';

import { MockProviders } from 'testUtils';
import { useMoveSearchQueries } from 'hooks/queries';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('hooks/queries'),
  useMoveSearchQueries: jest.fn(),
}));

const mockSearchResults = {
  searchResult: {
    data: [],
    page: 1,
    perPage: 20,
    totalCount: 0,
  },
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const mockSearchResultsWithMove = {
  searchResult: {
    data: [
      {
        branch: 'ARMY',
        destinationDutyLocationPostalCode: '30813',
        destinationGBLOC: 'CNNQ',
        dodID: '5805438095',
        firstName: 'Spaceman',
        id: '94b7e1a4-6ca6-439c-af86-24a135bacb7b',
        lastName: 'Mover',
        locator: 'MOVE12',
        orderType: 'PERMANENT_CHANGE_OF_STATION',
        originDutyLocationPostalCode: '50309',
        originGBLOC: 'KKFA',
        requestedPickupDate: '2020-03-16',
        shipmentsCount: 1,
        status: 'SUBMITTED',
      },
    ],
    page: 1,
    perPage: 20,
    totalCount: 1,
  },
  isLoading: false,
  isError: false,
  isSuccess: true,
};

describe('QAECSRMoveSearch page', () => {
  it('page loads', async () => {
    render(
      <MockProviders>
        <QAECSRMoveSearch />
      </MockProviders>,
    );

    const h1 = await screen.getByRole('heading', { name: 'Search for a move', level: 1 });
    expect(h1).toBeInTheDocument();

    const results = screen.queryByText(/Results/);
    expect(results).not.toBeInTheDocument();
  });

  it('can submit a search by DOD ID', async () => {
    useMoveSearchQueries.mockReturnValue(mockSearchResults);
    render(
      <MockProviders>
        <QAECSRMoveSearch />
      </MockProviders>,
    );
    await act(async () => {
      const submitButton = screen.getByTestId('searchTextSubmit');
      await screen.getByLabelText('DOD ID').click();
      await userEvent.type(screen.getByLabelText('Search'), '1234567890');
      await waitFor(() => {
        expect(screen.getByLabelText('Search')).toHaveValue('1234567890');
        expect(screen.getByLabelText('DOD ID')).toBeChecked();
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);

      const results = await screen.queryByText(/Results/);
      expect(results).toBeInTheDocument();
    });
  });

  it('can submit a search by Move Code', async () => {
    useMoveSearchQueries.mockReturnValue(mockSearchResults);
    render(
      <MockProviders>
        <QAECSRMoveSearch />
      </MockProviders>,
    );

    await act(async () => {
      const submitButton = screen.getByTestId('searchTextSubmit');
      await screen.getByLabelText('Move Code').click();
      await userEvent.type(screen.getByLabelText('Search'), 'MOVE12');
      await waitFor(() => {
        expect(screen.getByLabelText('Search')).toHaveValue('MOVE12');
        expect(screen.getByLabelText('Move Code')).toBeChecked();
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);

      const results = await screen.queryByText(/Results/);
      expect(results).toBeInTheDocument();
      const noResults = await screen.queryByText(/No results found/);
      expect(noResults).toBeInTheDocument();
    });
  });

  it('can submit a search by Customer Name', async () => {
    useMoveSearchQueries.mockReturnValue(mockSearchResults);
    render(
      <MockProviders>
        <QAECSRMoveSearch />
      </MockProviders>,
    );

    const submitButton = screen.getByTestId('searchTextSubmit');
    await act(async () => {
      await screen.getByLabelText('Customer Name').click();
      await userEvent.type(screen.getByLabelText('Search'), 'Leo Spaceman');
      await waitFor(() => {
        expect(screen.getByLabelText('Search')).toHaveValue('Leo Spaceman');
        expect(screen.getByLabelText('Customer Name')).toBeChecked();
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);

      const results = await screen.queryByText(/Results/);
      expect(results).toBeInTheDocument();
      const noResults = await screen.queryByText(/No results found/);
      expect(noResults).toBeInTheDocument();
    });
  });

  it('can submit a search by Payment Request Number', async () => {
    useMoveSearchQueries.mockReturnValue(mockSearchResults);
    render(
      <MockProviders>
        <QAECSRMoveSearch />
      </MockProviders>,
    );
    await act(async () => {
      const submitButton = screen.getByTestId('searchTextSubmit');
      await screen.getByLabelText('Payment Request Number').click();
      await userEvent.type(screen.getByLabelText('Search'), '1234-5678-9');
      await waitFor(() => {
        expect(screen.getByLabelText('Search')).toHaveValue('1234-5678-9');
        expect(screen.getByLabelText('Payment Request Number')).toBeChecked();
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);

      const results = await screen.queryByText(/Results/);
      expect(results).toBeInTheDocument();
      const noResults = await screen.queryByText(/No results found/);
      expect(noResults).toBeInTheDocument();
    });
  });

  it('can navigate to move afer search', async () => {
    useMoveSearchQueries.mockReturnValue(mockSearchResultsWithMove);
    render(
      <MockProviders>
        <QAECSRMoveSearch />
      </MockProviders>,
    );
    await act(async () => {
      const submitButton = screen.getByTestId('searchTextSubmit');
      await screen.getByLabelText('Move Code').click();
      await userEvent.type(screen.getByLabelText('Search'), 'MOVE12');
      await waitFor(() => {
        expect(screen.getByLabelText('Search')).toHaveValue('MOVE12');
        expect(screen.getByLabelText('Move Code')).toBeChecked();
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);

      const noResults = await screen.queryByText('Results (1)');
      expect(noResults).toBeInTheDocument();

      await screen.getByTestId('locator-0').click();
      expect(mockNavigate).toHaveBeenCalledWith('/moves/MOVE12/details');
    });
  });
});
