import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import UploadPaymentRequestDocuments from './UploadPaymentRequestDocuments';

// import { createUpload } from 'services/primeApi';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ paymentRequestId: 'testPaymentRequestId' }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createUpload: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('Upload Payment Request Documents Page', () => {
  it('renders the page without errors', () => {
    render(<UploadPaymentRequestDocuments />);

    expect(screen.getByText('Upload Payment Request Documents')).toBeInTheDocument();
  });

  it('navigates the user to the home page when the cancel button is clicked', async () => {
    render(<UploadPaymentRequestDocuments />);

    const cancel = screen.getByRole('button', { name: 'Cancel' });
    userEvent.click(cancel);

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/');
    });
  });
});
