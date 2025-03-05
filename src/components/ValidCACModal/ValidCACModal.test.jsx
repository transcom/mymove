import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import { ValidCACModal } from './ValidCACModal';

describe('ValidCACModal', () => {
  const onCloseMock = jest.fn();
  const onSubmitMock = jest.fn();

  beforeEach(() => {
    onCloseMock.mockClear();
    onSubmitMock.mockClear();
  });

  it('renders the modal with title, image, and description', () => {
    render(<ValidCACModal onClose={onCloseMock} onSubmit={onSubmitMock} />);

    const heading = screen.getByRole('heading', { name: /do you have a valid cac\?/i });
    expect(heading).toBeInTheDocument();

    const image = screen.getByRole('img');
    expect(image).toBeInTheDocument();

    expect(
      screen.getByText(/Common Access Card \(CAC\) authentication is required at first sign-in/i),
    ).toBeInTheDocument();
  });

  it('calls onSubmit when the "Yes" button is clicked', () => {
    render(<ValidCACModal onClose={onCloseMock} onSubmit={onSubmitMock} />);

    const yesButton = screen.getByTestId('modalSubmitButton');
    fireEvent.click(yesButton);

    expect(onSubmitMock).toHaveBeenCalledTimes(1);
  });

  it('calls onClose when the "No" button is clicked', () => {
    render(<ValidCACModal onClose={onCloseMock} onSubmit={onSubmitMock} />);

    const noButton = screen.getByTestId('modalBackButton');
    fireEvent.click(noButton);

    expect(onCloseMock).toHaveBeenCalledTimes(1);
  });
});
