import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ShipmentSITDisplay from '../ShipmentSITDisplay/ShipmentSITDisplay';

import SubmitSITExtensionModal from './SubmitSITExtensionModal';

describe('SubmitSITExtensionModal', () => {
  const summarySITExtension = (
    <ShipmentSITDisplay {...{ sitExtensions: [], sitStatus: {}, shipment: {}, hideSITExtensionAction: true }} />
  );
  it('calls onSubmit prop on approval with form values when validations pass', async () => {
    const mockOnSubmit = jest.fn();
    render(
      <SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} summarySITComponent={summarySITExtension} />,
    );
    const reasonInput = screen.getByLabelText('Reason for edit');
    const daysApprovedInput = screen.getByLabelText('Days approved');
    const officeRemarksInput = screen.getByLabelText('Office remarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']);
    await userEvent.type(daysApprovedInput, '20');
    await userEvent.type(officeRemarksInput, 'Approved!');
    await userEvent.click(submitBtn);

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
    render(
      <SubmitSITExtensionModal onSubmit={mockOnSubmit} onClose={() => {}} summarySITComponent={summarySITExtension} />,
    );
    const reasonInput = screen.getByLabelText('Reason for edit');
    const daysApprovedInput = screen.getByLabelText('Days approved');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.selectOptions(reasonInput, ['SERIOUS_ILLNESS_MEMBER']);
    await userEvent.type(daysApprovedInput, '0');

    await waitFor(() => {
      expect(submitBtn).toBeDisabled();
    });
  });

  it('calls onclose prop on modal close', async () => {
    const mockClose = jest.fn();
    render(
      <SubmitSITExtensionModal onSubmit={() => {}} onClose={mockClose} summarySITComponent={summarySITExtension} />,
    );
    const closeBtn = screen.getByRole('button', { name: 'Cancel' });

    await userEvent.click(closeBtn);

    await waitFor(() => {
      expect(mockClose).toHaveBeenCalled();
    });
  });

  it('renders the summary SIT component', async () => {
    render(
      <SubmitSITExtensionModal onSubmit={jest.fn()} onClose={jest.fn()} summarySITComponent={summarySITExtension} />,
    );

    await waitFor(() => {
      expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
    });
  });
});
