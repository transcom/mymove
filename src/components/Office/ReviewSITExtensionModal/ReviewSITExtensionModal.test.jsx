import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ReviewSITExtensionModal from './ReviewSITExtensionModal';

describe('ReviewSITExtensionModal', () => {
  const sitExt = {
    requestedDays: 45,
    requestReason: 'AWAITING_COMPLETION_OF_RESIDENCE',
    contractorRemarks: 'The customer requested an extension',
    id: '123',
  };

  const sitStatus = {
    totalDaysRemaining: 300,
  };

  const shipment = {
    sitDaysAllowance: 45,
  };

  it('renders requested days, reason, and contractor remarks', async () => {
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        onSubmit={() => {}}
        onClose={() => {}}
        shipment={shipment}
        sitStatus={sitStatus}
      />,
    );

    await waitFor(() => {
      expect(screen.getByText('45')).toBeInTheDocument();
      expect(screen.getByText('Awaiting completion of residence under construction')).toBeInTheDocument();
      expect(screen.getByText('The customer requested an extension')).toBeInTheDocument();
    });
  });

  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        onSubmit={mockOnSubmit}
        onClose={() => {}}
        shipment={shipment}
        sitStatus={sitStatus}
      />,
    );
    const daysApprovedInput = screen.getByLabelText('Days approved');
    const officeRemarksInput = screen.getByLabelText('Office remarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.type(daysApprovedInput, '{backspace}{backspace}20');
    await userEvent.type(officeRemarksInput, 'Approved!');
    await userEvent.click(submitBtn);

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
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        onSubmit={mockOnSubmit}
        onClose={() => {}}
        shipment={shipment}
        sitStatus={sitStatus}
      />,
    );
    const denyExtenstionField = screen.getByLabelText('No');
    const officeRemarksInput = screen.getByLabelText('Office remarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(denyExtenstionField);
    await userEvent.type(officeRemarksInput, 'Denied!');
    await userEvent.click(submitBtn);

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
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        onSubmit={mockOnSubmit}
        onClose={() => {}}
        shipment={shipment}
        sitStatus={sitStatus}
      />,
    );
    const daysApprovedInput = screen.getByLabelText('Days approved');
    const denyExtenstionField = screen.getByLabelText('No');
    await waitFor(() => {
      expect(daysApprovedInput).toBeInTheDocument();
    });
    await userEvent.click(denyExtenstionField);
    await waitFor(() => {
      expect(daysApprovedInput).not.toBeInTheDocument();
    });
  });

  it('does not allow submission of more days approved than are requested', async () => {
    const mockOnSubmit = jest.fn();
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        onSubmit={mockOnSubmit}
        onClose={() => {}}
        shipment={shipment}
        sitStatus={sitStatus}
      />,
    );
    const daysApprovedInput = screen.getByLabelText('Days approved');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.type(daysApprovedInput, '{backspace}{backspace}46');

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('does not allow submission of 0 approved days', async () => {
    const mockOnSubmit = jest.fn();
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        shipment={shipment}
        sitStatus={sitStatus}
        onSubmit={mockOnSubmit}
        onClose={() => {}}
      />,
    );
    const daysApprovedInput = screen.getByLabelText('Days approved');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.type(daysApprovedInput, '{backspace}{backspace}0');

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        shipment={shipment}
        sitStatus={sitStatus}
        onSubmit={() => {}}
        onClose={mockClose}
      />,
    );
    const closeBtn = screen.getByRole('button', { name: 'Cancel' });

    await userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });

  it('renders the summary SIT component', async () => {
    render(
      <ReviewSITExtensionModal
        sitExtension={sitExt}
        shipment={shipment}
        sitStatus={sitStatus}
        onSubmit={jest.fn()}
        onClose={jest.fn()}
      />,
    );

    await waitFor(() => {
      expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
    });
  });
});
