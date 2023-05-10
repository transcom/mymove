import React from 'react';
import { screen, waitFor } from '@testing-library/react';

import InvalidPermissions from './InvalidPermissions';

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

describe('InvalidPermissions', () => {
  it('component renders', async () => {
    renderWithRouter(<InvalidPermissions />, mockRoutingOptions);

    await waitFor(() => {
      expect(
        screen.getByRole('heading', { level: 1, name: 'You do not have permission to access this site.' }),
      ).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Sign Out' })).toBeInTheDocument();
    });
  });
});
