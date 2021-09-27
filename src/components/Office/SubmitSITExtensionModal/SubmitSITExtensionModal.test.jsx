import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import SubmitSITExtensionModal from './SubmitSITExtensionModal';

describe('SubmitSITExtensionModal', () => {
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    const { getByRole, getByLabelText } = render(
      <SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} />,
    );
    const reasonInput = getByLabelText('Reason for edit');
    const daysApprovedInput = getByLabelText('Days approved');
    const officeRemarksInput = getByLabelText('Office remarks');
    const submitBtn = getByRole('button', { name: 'Save' });

    userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']);
    userEvent.type(daysApprovedInput, '20');
    userEvent.type(officeRemarksInput, 'Approved!');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith({
        requestReason: 'SERIOUS_ILLNESS_MEMBER',
        daysApproved: '20',
        officeRemarks: 'Approved!',
      });
    });
  });

  it('does not allow submission of 0 approved days', async () => {
    const mockOnSubmit = jest.fn();
    const { getByRole, getByLabelText } = render(
      <SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} />,
    );
    const reasonInput = getByLabelText('Reason for edit');
    const daysApprovedInput = getByLabelText('Days approved');
    const submitBtn = getByRole('button', { name: 'Save' });

    userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']);
    userEvent.type(daysApprovedInput, '0');

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    const { getByRole } = render(<SubmitSITExtensionModal onSubmit={() => {}} onClose={mockClose} />);
    const closeBtn = getByRole('button', { name: 'Cancel' });

    userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });
});
