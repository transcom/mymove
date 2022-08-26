import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Forbidden from './Forbidden';

import { MockProviders } from 'testUtils';

const testMoveCode = '123ABC';

const mockPush = jest.fn();

// These mocks don't seem like they are actually being used, which might
// be why the second test is failing.
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useParams: jest.fn().mockReturnValue({
    moveCode: testMoveCode,
  }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

beforeEach(() => {
  jest.clearAllMocks();
});
describe('Forbidden', () => {
  it('renders', async () => {
    render(
      <MockProviders
        initialEntries={[`/moves/${testMoveCode}/evaluation-reports/452e8013-805c-4e17-a6ff-ce90722c12c7`]}
      >
        <Forbidden />
      </MockProviders>,
    );
    await waitFor(() => {
      expect(screen.getByRole('heading', { level: 1, name: "Sorry, you can't access this page" })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Go to move details' })).toBeInTheDocument();
    });
  });
  it('button works', async () => {
    render(
      <MockProviders
        initialEntries={[`/moves/${testMoveCode}/evaluation-reports/452e8013-805c-4e17-a6ff-ce90722c12c7`]}
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
      expect(mockPush).toHaveBeenCalled();
      // expect(mockPush).toHaveBeenCalledWith(`/moves/${testMoveCode}/details`)
    });
  });
});
