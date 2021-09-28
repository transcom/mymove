import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewSITExtensionModal from './ReviewSITExtensionModal';

describe('ReviewSITExtensionModal', () => {
  const sitExt = {
    requestedDays: 45,
    requestReason: 'AWAITING_COMPLETION_OF_RESIDENCE',
    contractorRemarks: 'The customer requested an extension',
    id: '123',
  };

  it('renders requested days, reason, and contractor remarks', async () => {
    const { getByText, findByText } = render(
      <ReviewSITExtensionModal sitExtension={sitExt} onSubmit={() => {}} onClose={() => {}} />,
    );

    expect(await findByText('45')).toBeInTheDocument();
    expect(getByText('Awaiting completion of residence under construction')).toBeInTheDocument();
    expect(getByText('The customer requested an extension')).toBeInTheDocument();
  });

  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    const { getByRole, getByLabelText } = render(
      <ReviewSITExtensionModal sitExtension={sitExt} onSubmit={mockOnSubmit} onClose={() => {}} />,
    );
    const daysApprovedInput = getByLabelText('Days approved');
    const officeRemarksInput = getByLabelText('Office remarks');
    const submitBtn = getByRole('button', { name: 'Save' });

    userEvent.type(daysApprovedInput, '{backspace}{backspace}20');
    userEvent.type(officeRemarksInput, 'Approved!');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith(sitExt.id, {
        acceptExtension: 'yes',
        daysApproved: '20',
        officeRemarks: 'Approved!',
      });
    });
  });

  it('calls onSubmit prop on denial with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    const { getByRole, getByLabelText } = render(
      <ReviewSITExtensionModal sitExtension={sitExt} onSubmit={mockOnSubmit} onClose={() => {}} />,
    );
    const denyExtenstionField = getByLabelText('No');
    const officeRemarksInput = getByLabelText('Office remarks');
    const submitBtn = getByRole('button', { name: 'Save' });

    userEvent.click(denyExtenstionField);
    userEvent.type(officeRemarksInput, 'Denied!');
    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith(sitExt.id, {
        acceptExtension: 'no',
        daysApproved: '',
        officeRemarks: 'Denied!',
      });
    });
  });

  it('hides days approved input when no is selected', async () => {
    const mockOnSubmit = jest.fn();
    const { getByLabelText } = render(
      <ReviewSITExtensionModal sitExtension={sitExt} onSubmit={mockOnSubmit} onClose={() => {}} />,
    );
    const daysApprovedInput = getByLabelText('Days approved');
    const denyExtenstionField = getByLabelText('No');
    await waitFor(() => {
      expect(daysApprovedInput).toBeInTheDocument();
    });
    userEvent.click(denyExtenstionField);
    await waitFor(() => {
      expect(daysApprovedInput).not.toBeInTheDocument();
    });
  });

  it('does not allow submission of more days approved than are requested', async () => {
    const mockOnSubmit = jest.fn();
    const { getByRole, getByLabelText } = render(
      <ReviewSITExtensionModal sitExtension={sitExt} onSubmit={mockOnSubmit} onClose={() => {}} />,
    );
    const daysApprovedInput = getByLabelText('Days approved');
    const submitBtn = getByRole('button', { name: 'Save' });

    userEvent.type(daysApprovedInput, '{backspace}{backspace}46');

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('does not allow submission of 0 approved days', async () => {
    const mockOnSubmit = jest.fn();
    const { getByRole, getByLabelText } = render(
      <ReviewSITExtensionModal sitExtension={sitExt} onSubmit={mockOnSubmit} onClose={() => {}} />,
    );
    const daysApprovedInput = getByLabelText('Days approved');
    const submitBtn = getByRole('button', { name: 'Save' });

    userEvent.type(daysApprovedInput, '{backspace}{backspace}0');

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    const { getByRole } = render(
      <ReviewSITExtensionModal sitExtension={sitExt} onSubmit={() => {}} onClose={mockClose} />,
    );
    const closeBtn = getByRole('button', { name: 'Cancel' });

    userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });
});
