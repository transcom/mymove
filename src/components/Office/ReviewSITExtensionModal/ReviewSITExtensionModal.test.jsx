import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import moment from 'moment';
import React from 'react';

import ReviewSITExtensionModal from './ReviewSITExtensionModal';

import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';

describe('ReviewSITExtensionModal', () => {
  const sitExt = {
    requestedDays: 45,
    requestReason: 'AWAITING_COMPLETION_OF_RESIDENCE',
    contractorRemarks: 'The customer requested an extension',
    id: '123',
  };

  const sitStatus = {
    totalDaysRemaining: 30,
    totalSITDaysUsed: 15,
    calculatedTotalDaysInSIT: 15,
    currentSIT: {
      daysInSIT: 15,
      sitEntryDate: moment().subtract(15, 'days').format(swaggerDateFormat),
      sitAuthorizedEndDate: moment().add(15, 'days').format(swaggerDateFormat),
    },
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

    const daysApprovedInput = screen.getByTestId('daysApproved');
    await userEvent.clear(daysApprovedInput);
    await userEvent.type(daysApprovedInput, '90');

    const acceptExtensionField = screen.getByLabelText('Yes');
    await userEvent.click(acceptExtensionField);

    const reasonDropdown = screen.getByLabelText('Reason for edit');
    await userEvent.selectOptions(reasonDropdown, ['SERIOUS_ILLNESS_MEMBER']);

    const officeRemarksInput = screen.getByLabelText('Office remarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.type(officeRemarksInput, 'Approved!');
    await userEvent.click(submitBtn);

    const expectedEndDate = formatDateForDatePicker(moment().add(75, 'days').subtract(1, 'day'));

    await waitFor(() => {
      // expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith(sitExt.id, {
        acceptExtension: 'yes',
        convertToCustomerExpense: false,
        requestReason: 'SERIOUS_ILLNESS_MEMBER',
        officeRemarks: 'Approved!',
        daysApproved: '90',
        sitEndDate: expectedEndDate,
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
    const denyExtensionField = screen.getByLabelText('No');
    const officeRemarksInput = screen.getByLabelText('Office remarks');
    const submitBtn = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(denyExtensionField);
    await userEvent.type(officeRemarksInput, 'Denied!');
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
      expect(mockOnSubmit).toHaveBeenCalledWith(
        sitExt.id,
        expect.objectContaining({
          acceptExtension: 'no',
          convertToCustomerExpense: false,
          officeRemarks: 'Denied!',
        }),
      );
    });
  });

  it('hides Reason for edit selection when no is selected', async () => {
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
    const acceptExtensionField = screen.getByLabelText('Yes');
    await userEvent.click(acceptExtensionField);
    const denyExtensionField = screen.getByLabelText('No');
    const reasonInput = screen.getByLabelText('Reason for edit');
    await waitFor(() => {
      expect(reasonInput).toBeInTheDocument();
    });
    await userEvent.click(denyExtensionField);
    await waitFor(() => {
      expect(reasonInput).not.toBeInTheDocument();
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

    const daysApprovedInput = screen.getByTestId('daysApproved');

    const acceptExtensionField = screen.getByLabelText('Yes');
    await userEvent.click(acceptExtensionField);

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
    const sitStartAndEndTable = await screen.findByTestId('sitStartAndEndTable');
    expect(sitStartAndEndTable).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('Calculated total SIT days')).toBeInTheDocument();
    expect(within(sitStartAndEndTable).getByText('15')).toBeInTheDocument();
  });
});
