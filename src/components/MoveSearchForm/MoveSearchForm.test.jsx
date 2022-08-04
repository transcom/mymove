import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MoveSearchForm from './MoveSearchForm';

describe('MoveSearchForm', () => {
  it('renders', () => {
    const { getByText } = render(<MoveSearchForm onSubmit={() => {}} />);
    expect(getByText('What do you want to search for?')).toBeInTheDocument();
  });

  describe('check validation', () => {
    it('can submit move code', async () => {
      const onSubmit = jest.fn();
      const { getByLabelText, getByRole } = render(<MoveSearchForm onSubmit={onSubmit} />);
      const submitButton = getByRole('button');

      await userEvent.type(getByLabelText('Search'), '123456');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('123456');
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);
      await waitFor(() => {
        expect(onSubmit).toHaveBeenCalledWith(
          {
            searchText: '123456',
            searchType: 'moveCode',
          },
          expect.anything(),
        );
      });
    });

    it('can submit DOD ID', async () => {
      const onSubmit = jest.fn();
      const { getByLabelText, getByRole } = render(<MoveSearchForm onSubmit={onSubmit} />);
      const submitButton = getByRole('button');

      await userEvent.click(getByLabelText('DOD ID'));
      await userEvent.type(getByLabelText('Search'), '1111111111');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('1111111111');
      });
      expect(submitButton).toBeEnabled();
      await userEvent.click(submitButton);
      await waitFor(() => {
        expect(onSubmit).toHaveBeenCalledWith(
          {
            searchText: '1111111111',
            searchType: 'dodID',
          },
          expect.anything(),
        );
      });
    });

    it('can submit name', async () => {
      const onSubmit = jest.fn();
      const { getByLabelText, getByRole } = render(<MoveSearchForm onSubmit={onSubmit} />);
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

    it('submits move code when it is 6 characters', async () => {
      const { getByLabelText, getByRole } = render(<MoveSearchForm onSubmit={jest.fn()} />);
      await userEvent.click(getByLabelText('Move Code'));
      await userEvent.type(getByLabelText('Search'), '123456');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('123456');
      });
      expect(getByRole('button')).toBeEnabled();
    });

    it('disables submit button when dod id is not 10 characters', async () => {
      const { getByLabelText, getByRole } = render(<MoveSearchForm onSubmit={jest.fn()} />);
      await userEvent.click(getByLabelText('DOD ID'));
      await userEvent.type(getByLabelText('Search'), '12345');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('12345');
      });
      expect(getByRole('button')).toBeDisabled();
      await userEvent.type(getByLabelText('Search'), '67890');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('1234567890');
      });
      expect(getByRole('button')).toBeEnabled();
      await userEvent.type(getByLabelText('Search'), '1');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('12345678901');
      });
      expect(getByRole('button')).toBeDisabled();
    });

    it('disables submit button when move code is not 6 characters', async () => {
      const { getByLabelText, getByRole } = render(<MoveSearchForm onSubmit={jest.fn()} />);
      await userEvent.click(getByLabelText('Move Code'));
      await userEvent.type(getByLabelText('Search'), '12345');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('12345');
      });
      expect(getByRole('button')).toBeDisabled();
      await userEvent.type(getByLabelText('Search'), '6');
      await waitFor(() => {
        expect(getByLabelText('Search')).toHaveValue('123456');
      });
      expect(getByRole('button')).toBeEnabled();
    });
  });
});
