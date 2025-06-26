import React from 'react';
import { render, waitFor, screen, act, fireEvent, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import moment from 'moment';

import SubmitSITExtensionModal from './SubmitSITExtensionModal';

import { formatDateForDatePicker, utcDateFormat } from 'shared/dates';

const defaultValues = {
  sitStatus: {
    totalDaysRemaining: 210,
    totalSITDaysUsed: 60,
    calculatedTotalDaysInSIT: 60,
    currentSIT: {
      location: 'DESTINATION',
      daysInSIT: 60,
      sitEntryDate: moment().subtract(60, 'days').format(utcDateFormat),
    },
  },
  shipment: {
    sitDaysAllowance: 270,
  },
};

describe('SubmitSITExtensionModal', () => {
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const reasonInput = screen.getByLabelText('Reason for edit *');
    const daysApprovedInput = screen.getByTestId('daysApproved');
    const officeRemarksInput = screen.getByLabelText('Office remarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await act(() => userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']));
    await act(() => userEvent.clear(daysApprovedInput));
    await act(() => userEvent.type(daysApprovedInput, '280'));
    await act(() => userEvent.type(officeRemarksInput, 'Approved!'));
    await act(() => userEvent.click(submitBtn));

    const expectedEndDate = formatDateForDatePicker(moment().add(220, 'days').subtract(1, 'day'));
    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith({
        requestReason: 'SERIOUS_ILLNESS_MEMBER',
        daysApproved: '280',
        officeRemarks: 'Approved!',
        sitEndDate: expectedEndDate,
      });
    });
  });

  it('does not allow submission of 0 approved days', async () => {
    const mockOnSubmit = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const reasonInput = screen.getByLabelText('Reason for edit *');
    const daysApprovedInput = screen.getByTestId('daysApproved');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await act(() => userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']));
    await act(() => userEvent.clear(daysApprovedInput));
    await act(() => userEvent.type(daysApprovedInput, '0'));
    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('changes the end date when the total days of SIT approved is changed', async () => {
    const mockOnSubmit = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const reasonInput = screen.getByLabelText('Reason for edit *');
    const daysApprovedInput = screen.getByTestId('daysApproved');
    const sitEndDateInput = screen.getByPlaceholderText('DD MMM YYYY');

    await act(() => userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']));
    await act(() => userEvent.clear(daysApprovedInput));
    await act(() => userEvent.type(daysApprovedInput, '280'));

    const expectedEndDate = formatDateForDatePicker(moment().add(220, 'days').subtract(1, 'day'));
    expect(sitEndDateInput.value).toBe(expectedEndDate);
  });

  it('changes the total days of SIT approved when end date is changed', async () => {
    const mockOnSubmit = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const sitEndDateInput = screen.getByPlaceholderText('DD MMM YYYY');
    await act(() => userEvent.clear(sitEndDateInput));
    const newEndDate = formatDateForDatePicker(moment().add(220, 'days').subtract(1, 'day'));
    await act(() => userEvent.type(sitEndDateInput, newEndDate));
    await fireEvent.blur(sitEndDateInput);
    const daysApprovedInput = screen.getByTestId('daysApproved');
    expect(daysApprovedInput.value).toBe('280');
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={() => {}} onClose={mockClose} {...defaultValues} />);
    const closeBtn = screen.getByRole('button', { name: 'Cancel' });

    await act(() => userEvent.click(closeBtn));

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });

  it('renders the summary SIT component and asterisks for required fields', async () => {
    await render(<SubmitSITExtensionModal onSubmit={jest.fn()} onClose={jest.fn()} {...defaultValues} />);

    await waitFor(() => {
      expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
    });
    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
    const sitStartAndEndTable = await screen.findByTestId('sitStartAndEndTable');
    expect(sitStartAndEndTable).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('Calculated total SIT days')).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('60')).toBeInTheDocument();

    const totalDaysSITApproved = screen.getByRole('columnheader', { name: 'Total days of SIT approved' });
    expect(within(totalDaysSITApproved).getByText('*')).toBeInTheDocument();

    const sitAuthEndDate = screen.getByRole('columnheader', { name: 'SIT authorized end date' });
    expect(within(sitAuthEndDate).getByText('*')).toBeInTheDocument();
  });
});
