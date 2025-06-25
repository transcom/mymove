import React from 'react';
import { screen, waitFor, fireEvent, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import moment from 'moment';
import { generatePath } from 'react-router-dom';

import { Agreement } from './Agreement';

import { submitMoveForApproval } from 'services/internalApi';
import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import MOVE_STATUSES from 'constants/moves';
import { customerRoutes } from 'constants/routes';
import { renderWithRouterProp } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  submitMoveForApproval: jest.fn(),
}));

afterEach(jest.resetAllMocks);

describe('Agreement page', () => {
  const testProps = {
    moveId: 'testMove123',
    setFlashMessage: jest.fn(),
    updateMove: jest.fn(),
  };

  const submittedMoveSuccessResponse = {
    id: testProps.moveId,
    status: MOVE_STATUSES.SUBMITTED,
  };

  const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId: testProps.moveId });

  it('submits the move and sets the flash message before redirecting home', async () => {
    submitMoveForApproval.mockResolvedValueOnce(submittedMoveSuccessResponse);

    renderWithRouterProp(
      <Agreement {...testProps} serviceMember={{ first_name: 'Sofia', last_name: 'Clark-Nuñez' }} />,
      {
        path: customerRoutes.MOVE_REVIEW_PATH,
        params: { moveId: 'testMove123' },
      },
    );

    const scrollBox = screen.getByTestId('certificationTextBox');
    const checkbox = await screen.findByRole('checkbox', {
      name: /I have read and understand the agreement as shown above/i,
    });
    const signatureInput = screen.getByLabelText('SIGNATURE');
    const completeButton = screen.getByRole('button', { name: 'Complete' });

    // all controls should start of disabled
    expect(checkbox).toHaveAttribute('readonly');
    expect(signatureInput).toHaveAttribute('readonly');
    expect(completeButton).toBeDisabled();
    Object.defineProperty(scrollBox, 'scrollHeight', { configurable: true, value: 300 });
    Object.defineProperty(scrollBox, 'clientHeight', { configurable: true, value: 100 });
    Object.defineProperty(scrollBox, 'scrollTop', {
      configurable: true,
      writable: true,
      value: 200,
    });
    await act(async () => {
      fireEvent.scroll(scrollBox);
    });

    // scroll to bottom should enable the checkbox, but not signature (yet)
    expect(checkbox).not.toHaveAttribute('readonly');
    expect(signatureInput).toHaveAttribute('readonly');
    await userEvent.click(checkbox);
    expect(checkbox.checked).toEqual(true);

    expect(signatureInput).toBeEnabled();
    await userEvent.type(signatureInput, 'Sofia Clark-Nuñez');

    await waitFor(() => {
      expect(completeButton).toBeEnabled();
    });

    await userEvent.click(completeButton);

    await waitFor(() => {
      expect(submitMoveForApproval).toHaveBeenCalledWith(testProps.moveId, {
        certification_text: completeCertificationText,
        date: moment().format(),
        signature: 'Sofia Clark-Nuñez',
        certification_type: SIGNED_CERT_OPTIONS.SHIPMENT,
      });
    });

    expect(testProps.updateMove).toHaveBeenCalledWith(submittedMoveSuccessResponse);

    expect(testProps.setFlashMessage).toHaveBeenCalledWith(
      'MOVE_SUBMIT_SUCCESS',
      'success',
      'You’ve submitted your move request.',
    );
  });

  it('renders an error if submitting the move responds with a server error', async () => {
    submitMoveForApproval.mockRejectedValueOnce({ errors: { signature: 'Signature can not be blank.' } });

    renderWithRouterProp(
      <Agreement {...testProps} serviceMember={{ first_name: 'Sofia', last_name: 'Clark-Nuñez' }} />,
      {
        path: customerRoutes.MOVE_REVIEW_PATH,
        params: { moveId: 'testMove123' },
      },
    );

    const scrollBox = screen.getByTestId('certificationTextBox');

    Object.defineProperty(scrollBox, 'scrollHeight', { configurable: true, value: 300 });
    Object.defineProperty(scrollBox, 'clientHeight', { configurable: true, value: 100 });
    Object.defineProperty(scrollBox, 'scrollTop', {
      configurable: true,
      writable: true,
      value: 200,
    });
    await act(async () => {
      fireEvent.scroll(scrollBox);
    });
    const checkbox = await screen.findByRole('checkbox', { name: /i have read and understand/i });
    expect(checkbox).toBeEnabled();
    await userEvent.click(checkbox);

    const signatureInput = await screen.findByLabelText('SIGNATURE');
    expect(signatureInput).toBeEnabled();

    await userEvent.type(screen.getByLabelText('SIGNATURE'), 'Sofia Clark-Nuñez');
    await userEvent.click(screen.getByRole('button', { name: 'Complete' }));

    await waitFor(() => {
      expect(screen.getByTestId('alert')).toHaveTextContent('There was a problem saving your signature');
    });
  });

  it('routes back to the review page when the back button is clicked', async () => {
    renderWithRouterProp(<Agreement {...testProps} />, {
      path: customerRoutes.MOVE_REVIEW_PATH,
      params: { moveId: 'testMove123' },
    });
    await userEvent.click(screen.getByRole('button', { name: 'Back' }));

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(reviewPath);
    });
  });
});
