import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import Forbidden from './Forbidden';

import { MockProviders } from 'testUtils';

const testMoveCode = '123ABC';

const mockPush = jest.fn();
const mockGoBack = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useHistory: () => ({
    push: mockPush,
    goBack: mockGoBack,
  }),
}));

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
    userEvent.click(screen.getByRole('button', { name: 'Go to move details' }));
    // the stuff below is not working, we're sticking on the page
    // i tried this with and without mocking push.
    await waitFor(() => {
      // expect(mockPush).toHaveBeenCalled();
      expect(screen.getByRole('heading', { level: 1, name: 'Move details' })).toBeInTheDocument();
    });
  });
});
