import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CustomerSearchForm from './CustomerSearchForm';

import { searchCustomers } from 'services/ghcApi';

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  searchCustomers: jest.fn(),
}));

beforeEach(jest.resetAllMocks);

describe('CustomerSearchForm', () => {
  it('renders', () => {
    const { getByText } = render(<CustomerSearchForm onSubmit={() => {}} />);
    expect(getByText('What do you want to search for?')).toBeInTheDocument();
  });

  describe('check validation', () => {
    it('can submit DOD ID', async () => {
      const onSubmit = searchCustomers;
      const { getByLabelText, getByRole } = render(<CustomerSearchForm onSubmit={onSubmit} />);
      const submitButton = getByRole('button');

      await userEvent.click(getByLabelText('DOD ID'));
      await userEvent.type(getByLabelText('Search'), '4152341523');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('4152341523');
        expect(getByLabelText('DOD ID')).toBeChecked();
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);
      await waitFor(() => {
        expect(onSubmit).toHaveBeenCalledWith(
          {
            searchText: '4152341523',
            searchType: 'dodID',
          },
          expect.anything(),
        );
      });
    });

    it('can submit name', async () => {
      const onSubmit = searchCustomers;
      const { getByLabelText, getByRole } = render(<CustomerSearchForm onSubmit={onSubmit} />);
      const submitButton = getByRole('button');

      await userEvent.click(getByLabelText('Customer Name'));
      await userEvent.type(getByLabelText('Search'), 'Leo Spaceman');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('Leo Spaceman');
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);
      await waitFor(() => {
        expect(onSubmit).toHaveBeenCalledWith(
          {
            searchText: 'Leo Spaceman',
            searchType: 'customerName',
          },
          expect.anything(),
        );
      });
    });
  });
});
