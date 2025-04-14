import React from 'react';
import { render, waitFor, act, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import SubmitMoveForm from './SubmitMoveForm';

describe('SubmitMoveForm component', () => {
  const testProps = {
    onSubmit: jest.fn(),
    onPrint: jest.fn(),
    onBack: jest.fn(),
    currentUser: 'Test User',
    initialValues: { signature: '', date: '2021-01-20' },
  };

  it('renders the signature and date inputs', () => {
    const { getByLabelText } = render(<SubmitMoveForm {...testProps} />);
    expect(getByLabelText('SIGNATURE')).toBeInTheDocument();
    expect(getByLabelText('SIGNATURE')).toBeRequired();
    expect(getByLabelText('Date')).toBeInTheDocument();
    expect(getByLabelText('Date')).toBeDisabled();
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByTestId, getByText } = render(<SubmitMoveForm {...testProps} />);
    const submitBtn = getByTestId('wizardCompleteButton');

    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getByText('Required')).toBeInTheDocument();
    });
  });

  it('submits the form when it is valid', async () => {
    await act(async () => {
      render(<SubmitMoveForm {...testProps} />);
    });

    const docContainer = screen.getByTestId('certificationTextBox');

    // Mock scroll values to simulate reaching the bottom
    Object.defineProperty(docContainer, 'scrollHeight', {
      configurable: true,
      value: 300,
    });
    Object.defineProperty(docContainer, 'clientHeight', {
      configurable: true,
      value: 100,
    });
    Object.defineProperty(docContainer, 'scrollTop', {
      configurable: true,
      writable: true,
      value: 200,
    });

    // Trigger scroll event on the correct element
    fireEvent.scroll(docContainer);

    // Wait for checkbox to become enabled
    const checkbox = await screen.findByRole('checkbox', {
      name: /i have read and understand/i,
    });
    await waitFor(() => expect(checkbox).toBeEnabled());

    // Click the checkbox
    userEvent.click(checkbox);

    // Type into the signature input (should now be enabled)
    const signatureInput = await screen.findByLabelText('SIGNATURE');
    await waitFor(() => expect(signatureInput).toBeEnabled());
    await userEvent.type(signatureInput, testProps.currentUser);

    // Click the complete button
    const submitBtn = screen.getByTestId('wizardCompleteButton');
    await userEvent.click(submitBtn);

    // Wait for form submission
    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalled();
    });
  });

  it('implements the onPrint handler', async () => {
    const { getByText } = render(<SubmitMoveForm {...testProps} />);

    const printBtn = getByText('Print');
    await userEvent.click(printBtn);

    expect(testProps.onPrint).toHaveBeenCalled();
  });

  it('implements the onBack handler', async () => {
    const { getByTestId } = render(<SubmitMoveForm {...testProps} />);

    const backBtn = getByTestId('wizardBackButton');
    await userEvent.click(backBtn);

    expect(testProps.onBack).toHaveBeenCalled();
  });
});
