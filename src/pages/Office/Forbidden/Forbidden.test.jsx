import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Forbidden from './Forbidden';

import { MockProviders } from 'testUtils';

const testMoveCode = '123ABC';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({
    moveCode: testMoveCode,
  }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

describe('Forbidden', () => {
  it('component renders', async () => {
    render(
      <MockProviders
        initialEntries={[`/moves/${testMoveCode}/evaluation-reports/11111111-1111-1111-1111-111111111111`]}
      >
        <Forbidden />
      </MockProviders>,
    );
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 1, name: "Sorry, you can't access this page" })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Go to move details' })).toBeInTheDocument();
    });
  });

  it('has a button to go to Move Details page', async () => {
    render(
      <MockProviders
        initialEntries={[`/moves/${testMoveCode}/evaluation-reports/11111111-1111-1111-1111-111111111111`]}
      >
        <Forbidden />
      </MockProviders>,
    );
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 1, name: "Sorry, you can't access this page" })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Go to move details' })).toBeInTheDocument();
    });
    await userEvent.click(screen.getByRole('button', { name: 'Go to move details' }));
    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith(`/moves/${testMoveCode}/details`);
    });
  });
});
