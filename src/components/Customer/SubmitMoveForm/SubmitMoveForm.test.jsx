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
    expect(getByLabelText('SIGNATURE *')).toBeInTheDocument();
    expect(getByLabelText('SIGNATURE *')).toBeRequired();
    expect(getByLabelText('Date')).toBeInTheDocument();
    expect(getByLabelText('Date')).toHaveAttribute('readonly');
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
    const signatureInput = await screen.findByLabelText('SIGNATURE *');
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
  it('disables the signature input until the agreement checkbox is checked', async () => {
    render(<SubmitMoveForm {...testProps} />);

    const signatureInput = screen.getByLabelText('SIGNATURE *');
    const checkbox = screen.getByRole('checkbox', {
      name: /i have read and understand/i,
    });

    expect(signatureInput).toHaveAttribute('readonly');

    // Simulate scroll-to-bottom to enable checkbox
    const docContainer = screen.getByTestId('certificationTextBox');
    Object.defineProperty(docContainer, 'scrollHeight', { configurable: true, value: 300 });
    Object.defineProperty(docContainer, 'clientHeight', { configurable: true, value: 100 });
    Object.defineProperty(docContainer, 'scrollTop', { configurable: true, writable: true, value: 200 });
    fireEvent.scroll(docContainer);

    await waitFor(() => expect(checkbox).toBeEnabled());
    userEvent.click(checkbox);

    await waitFor(() => expect(signatureInput).toBeEnabled());
  });
  it('shows validation error if signature does not match currentUser', async () => {
    const mockSubmit = jest.fn();
    render(<SubmitMoveForm {...testProps} onSubmit={mockSubmit} />);

    // Simulate scroll-to-bottom to enable checkbox
    const docContainer = screen.getByTestId('certificationTextBox');
    Object.defineProperty(docContainer, 'scrollHeight', { configurable: true, value: 300 });
    Object.defineProperty(docContainer, 'clientHeight', { configurable: true, value: 100 });
    Object.defineProperty(docContainer, 'scrollTop', { configurable: true, writable: true, value: 200 });
    fireEvent.scroll(docContainer);

    // Wait for checkbox to be enabled
    const checkbox = await screen.findByRole('checkbox');
    await waitFor(() => expect(checkbox).toBeEnabled());
    await userEvent.click(checkbox);

    // Wait for signature input to become enabled
    const signatureInput = screen.getByLabelText('SIGNATURE *');
    await waitFor(() => expect(signatureInput).toBeEnabled());

    // Type mismatched signature
    await userEvent.clear(signatureInput);
    await userEvent.type(signatureInput, 'Wrong Name');

    const submitBtn = screen.getByTestId('wizardCompleteButton');
    await userEvent.click(submitBtn);

    // Expect error message
    const validationError = await screen.findByText((text) =>
      text.includes('Typed signature must match your exact user name'),
    );

    expect(validationError).toBeInTheDocument();
    expect(mockSubmit).not.toHaveBeenCalled();
  });

  it('does not show validation error if signature has extra spaces', async () => {
    const mockSubmit = jest.fn();
    render(<SubmitMoveForm {...testProps} onSubmit={mockSubmit} />);

    // Simulate scroll-to-bottom to enable checkbox
    const docContainer = screen.getByTestId('certificationTextBox');
    Object.defineProperty(docContainer, 'scrollHeight', { configurable: true, value: 300 });
    Object.defineProperty(docContainer, 'clientHeight', { configurable: true, value: 100 });
    Object.defineProperty(docContainer, 'scrollTop', { configurable: true, writable: true, value: 200 });
    fireEvent.scroll(docContainer);

    // Wait for checkbox to be enabled
    const checkbox = await screen.findByRole('checkbox');
    await waitFor(() => expect(checkbox).toBeEnabled());
    await userEvent.click(checkbox);

    // Wait for signature input to become enabled
    const signatureInput = screen.getByLabelText('SIGNATURE *');
    await waitFor(() => expect(signatureInput).toBeEnabled());

    // Type mismatched signature
    await userEvent.clear(signatureInput);
    await userEvent.type(signatureInput, 'Test  User');

    const submitBtn = screen.getByTestId('wizardCompleteButton');
    await userEvent.click(submitBtn);

    // Expect error message
    const validationError = await screen.queryByText((text) =>
      text.includes('Typed signature must match your exact user name'),
    );

    expect(validationError).not.toBeInTheDocument();
    expect(mockSubmit).toHaveBeenCalled();

    // Type mismatched signature
    await userEvent.clear(signatureInput);
    await userEvent.type(signatureInput, ' Test User ');

    await userEvent.click(submitBtn);

    // Expect error message
    const validationError2 = await screen.queryByText((text) =>
      text.includes('Typed signature must match your exact user name'),
    );

    expect(validationError2).not.toBeInTheDocument();
    expect(mockSubmit).toHaveBeenCalled();
  });

  it('does not render certification text if certificationText is null', () => {
    render(<SubmitMoveForm {...testProps} certificationText={null} />);

    // It still renders the wrapper, but there should be no markdown inside
    const certTextBox = screen.getByTestId('certificationTextBox');
    expect(certTextBox).toBeEmptyDOMElement();
  });
});
