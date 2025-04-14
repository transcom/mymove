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

  it('submits the form when its valid', async () => {
    await act(async () => {
      render(<SubmitMoveForm {...testProps} />);
    });

    const signatureInput = screen.getByLabelText('SIGNATURE');
    const submitBtn = screen.getByTestId('wizardCompleteButton');

    await act(async () => {
      await userEvent.type(signatureInput, testProps.currentUser);
    });
    await act(async () => {
      await fireEvent.click(submitBtn);
    });

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
