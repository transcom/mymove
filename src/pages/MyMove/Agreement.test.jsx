import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import moment from 'moment';
import { generatePath } from 'react-router';

import { Agreement } from './Agreement';

import { submitMoveForApproval } from 'services/internalApi';
import { completeCertificationText } from 'scenes/Legalese/legaleseText';
import { SIGNED_CERT_OPTIONS } from 'shared/constants';
import MOVE_STATUSES from 'constants/moves';
import { customerRoutes } from 'constants/routes';

jest.mock('services/internalApi', () => ({
  submitMoveForApproval: jest.fn(),
}));

afterEach(jest.resetAllMocks);

describe('Agreement page', () => {
  const testProps = {
    moveId: 'testMove123',
    setFlashMessage: jest.fn(),
    push: jest.fn(),
    updateMove: jest.fn(),
  };

  const submittedMoveSuccessResponse = {
    id: testProps.moveId,
    status: MOVE_STATUSES.SUBMITTED,
  };

  const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId: testProps.moveId });

  it('submits the move and sets the flash message before redirecting home', async () => {
    submitMoveForApproval.mockResolvedValueOnce(submittedMoveSuccessResponse);

    render(<Agreement {...testProps} />);
    await userEvent.type(screen.getByLabelText('Signature'), 'Sofia Clark-Nuñez');
    await userEvent.click(screen.getByRole('button', { name: 'Complete' }));

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

    render(<Agreement {...testProps} />);
    await userEvent.type(screen.getByLabelText('Signature'), 'Sofia Clark-Nuñez');
    await userEvent.click(screen.getByRole('button', { name: 'Complete' }));

    await waitFor(() => {
      expect(screen.getByTestId('alert')).toHaveTextContent('There was a problem saving your signature');
    });
  });

  it('routes back to the review page when the back button is clicked', async () => {
    render(<Agreement {...testProps} />);
    await userEvent.click(screen.getByRole('button', { name: 'Back' }));

    await waitFor(() => {
      expect(testProps.push).toHaveBeenCalledWith(reviewPath);
    });
  });
});
