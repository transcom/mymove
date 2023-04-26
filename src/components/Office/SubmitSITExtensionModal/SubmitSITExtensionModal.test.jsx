import React from 'react';
import { render, waitFor, screen, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import SubmitSITExtensionModal from './SubmitSITExtensionModal';

const defaultValues = {
  sitStatus: {
    daysInSIT: 30,
    location: 'DESTINATION',
    sitEntryDate: '2023-03-19T00:00:00.000Z',
    totalDaysRemaining: 210,
    totalSITDaysUsed: 60,
  },
  shipment: {
    sitDaysAllowance: 270,
  },
};

describe('SubmitSITExtensionModal', () => {
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const reasonInput = screen.getByLabelText('Reason for edit');
    const daysApprovedInput = screen.getByTestId('daysApproved');
    const officeRemarksInput = screen.getByLabelText('Office remarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await act(() => userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']));
    await act(() => userEvent.clear(daysApprovedInput));
    await act(() => userEvent.type(daysApprovedInput, '280'));
    await act(() => userEvent.type(officeRemarksInput, 'Approved!'));
    await act(() => userEvent.click(submitBtn));
    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith({
        requestReason: 'SERIOUS_ILLNESS_MEMBER',
        daysApproved: '280',
        officeRemarks: 'Approved!',
        sitEndDate: '24 Nov 2023',
      });
    });
  });

  it('does not allow submission of 0 approved days', async () => {
    const mockOnSubmit = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const reasonInput = screen.getByLabelText('Reason for edit');
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
    const reasonInput = screen.getByLabelText('Reason for edit');
    const daysApprovedInput = screen.getByTestId('daysApproved');
    const sitEndDateInput = screen.getByPlaceholderText('DD MMM YYYY');

    await act(() => userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']));
    await act(() => userEvent.clear(daysApprovedInput));
    await act(() => userEvent.type(daysApprovedInput, '280'));
    expect(sitEndDateInput.value).toBe('24 Nov 2023');
  });

  it('changes the total days of SIT approved when end date is changed', async () => {
    const mockOnSubmit = jest.fn();
    await render(<SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} {...defaultValues} />);
    const sitEndDateInput = screen.getByPlaceholderText('DD MMM YYYY');
    await act(() => userEvent.clear(sitEndDateInput));
    await act(() => userEvent.type(sitEndDateInput, '04 Nov 2023'));
    await fireEvent.blur(sitEndDateInput);
    const daysApprovedInput = screen.getByTestId('daysApproved');
    expect(daysApprovedInput.value).toBe('260');
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

  it('renders the summary SIT component', async () => {
    await render(<SubmitSITExtensionModal onSubmit={jest.fn()} onClose={jest.fn()} {...defaultValues} />);

    await waitFor(() => {
      expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
    });
  });
});
