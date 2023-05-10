import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Forbidden from './Forbidden';

import { renderWithRouter } from 'testUtils';
import { qaeCSRRoutes } from 'constants/routes';

const testMoveCode = '123ABC';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { moveCode: testMoveCode, reportId: '11111111-1111-1111-1111-111111111111' };
const mockRoutingOptions = { path: qaeCSRRoutes.BASE_EVALUATION_REPORT_PATH, params: routingParams };

describe('Forbidden', () => {
  it('component renders', async () => {
    renderWithRouter(<Forbidden />, mockRoutingOptions);

    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 1, name: "Sorry, you can't access this page." })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Go to move details' })).toBeInTheDocument();
    });
  });

  it('has a button to go to Move Details page', async () => {
    renderWithRouter(<Forbidden />, mockRoutingOptions);
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 1, name: "Sorry, you can't access this page." })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Go to move details' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Go to move details' }));
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(`/moves/${testMoveCode}/details`);
    });
  });
});
