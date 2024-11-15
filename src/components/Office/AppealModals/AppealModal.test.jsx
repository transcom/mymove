import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import '@testing-library/jest-dom/extend-expect';
import { AppealModal } from './AppealModal';

describe('AppealModal', () => {
  const mockOnClose = jest.fn();
  const mockOnSubmit = jest.fn();

  const renderComponent = (props = {}) => {
    render(
      <AppealModal
        onClose={mockOnClose}
        onSubmit={mockOnSubmit}
        isSeriousIncidentAppeal={props.isSeriousIncidentAppeal}
        selectedReportViolation={props.selectedReportViolation}
      />,
    );
  };

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('renders correctly with all fields and buttons', async () => {
    renderComponent();
    expect(screen.getByText('Leave Appeal Decision')).toBeInTheDocument();
    expect(screen.getByLabelText('Remarks')).toBeInTheDocument();
    expect(screen.getByLabelText('Sustained')).toBeInTheDocument();
    expect(screen.getByLabelText('Rejected')).toBeInTheDocument();
    await waitFor(() => expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled());
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
  });

  it('displays validation messages when form is submitted with empty fields', async () => {
    renderComponent();
    await waitFor(() => expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled());

    userEvent.click(screen.getByLabelText('Remarks'));
    userEvent.click(screen.getByTestId('sustainedRadio'));

    await waitFor(() => {
      expect(screen.getByText('Remarks are required')).toBeInTheDocument();
    });

    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled();
  });

  it('the form successfully submits when all required fields are filled out', async () => {
    renderComponent();
    await userEvent.type(screen.getByLabelText('Remarks'), 'These are my remarks');
    await userEvent.click(screen.getByTestId('sustainedRadio'));

    expect(screen.getByRole('button', { name: 'Save' })).toBeEnabled();

    userEvent.click(screen.getByRole('button', { name: /Save/i }));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith(
        {
          remarks: 'These are my remarks',
          appealStatus: 'sustained',
        },
        expect.anything(),
      );
    });
  });

  it('cancel button triggers onClose callback', async () => {
    renderComponent();
    await userEvent.click(screen.getByTestId('modalCancelButton'));
    expect(mockOnClose).toHaveBeenCalled();
  });

  it('displays "Serious Incident" when isSeriousIncidentAppeal is true', () => {
    renderComponent({ isSeriousIncidentAppeal: true });
    expect(screen.getByTestId('seriousIncidentModalHint')).toBeInTheDocument();
    expect(screen.getByText('Serious Incident')).toBeInTheDocument();
  });

  it('displays violation details when isSeriousIncidentAppeal is false and selectedReportViolation is provided', () => {
    const violation = {
      paragraphNumber: '1.2.3',
      title: 'Test Violation Title',
    };
    renderComponent({ isSeriousIncidentAppeal: false, selectedReportViolation: { violation } });

    expect(screen.getByTestId('violationModalHint')).toBeInTheDocument();
    expect(screen.getByText('1.2.3 Test Violation Title')).toBeInTheDocument();
  });
});
