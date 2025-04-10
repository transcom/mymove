import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import moment from 'moment';

import EditSitEntryDateModal from './EditSitEntryDateModal';

import { formatDateForDatePicker, utcDateFormat, swaggerDateFormat } from 'shared/dates';

const defaultValues = {
  serviceItem: {
    id: 'fakeID',
    currentSIT: {
      sitEntryDate: moment().subtract(60, 'days').format(utcDateFormat),
    },
  },
};

describe('EditSitEntryDateModal', () => {
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    await render(<EditSitEntryDateModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const officeRemarksInput = screen.getByTestId('officeRemarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    const datePickers = screen.getAllByPlaceholderText('DD MMM YYYY');
    const sitEntryDate = datePickers[1];
    const newEndDate = formatDateForDatePicker(moment().add(220, 'DD MM YYYY'));
    formatDateForDatePicker(moment(newEndDate, swaggerDateFormat));
    await waitFor(() => userEvent.type(sitEntryDate, newEndDate));
    await fireEvent.blur(sitEntryDate);
    await waitFor(() => userEvent.type(officeRemarksInput, 'Approved!'));
    await waitFor(() => userEvent.click(submitBtn));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith('fakeID', new Date(newEndDate));
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    await render(<EditSitEntryDateModal onSubmit={() => {}} onClose={mockClose} {...defaultValues} />);
    const closeBtn = screen.getByRole('button', { name: 'Cancel' });

    await waitFor(() => userEvent.click(closeBtn));

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });

  it('renders the summary SIT component', async () => {
    await render(<EditSitEntryDateModal onSubmit={jest.fn()} onClose={jest.fn()} {...defaultValues} />);

    await waitFor(() => {
      expect(screen.getByText('Edit SIT Entry Date')).toBeInTheDocument();
      expect(screen.getByText('Original SIT entry date')).toBeInTheDocument();
      expect(screen.getByText('New SIT entry date')).toBeInTheDocument();
      expect(screen.getByText('Office remarks')).toBeInTheDocument();
    });
  });
});
