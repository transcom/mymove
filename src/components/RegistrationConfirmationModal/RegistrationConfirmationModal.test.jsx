import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';

import '@testing-library/jest-dom/extend-expect';
import { RegistrationConfirmationModal } from './RegistrationConfirmationModal';

describe('RegistrationConfirmationModal', () => {
  const onSubmitMock = jest.fn();

  beforeEach(() => {
    onSubmitMock.mockClear();
  });

  it('renders the confirmation modal with expected text and button', () => {
    render(<RegistrationConfirmationModal onSubmit={onSubmitMock} />);

    expect(screen.getByText('Registration Confirmation')).toBeInTheDocument();
    expect(screen.getByText(/Your MilMove & Okta accounts have successfully been created/i)).toBeInTheDocument();
    const continueButton = screen.getByTestId('modalSubmitButton');
    expect(continueButton).toBeInTheDocument();
  });

  it('calls onSubmit when the Continue button is clicked', () => {
    render(<RegistrationConfirmationModal onSubmit={onSubmitMock} />);

    const continueButton = screen.getByTestId('modalSubmitButton');
    fireEvent.click(continueButton);

    expect(onSubmitMock).toHaveBeenCalledTimes(1);
  });
});
