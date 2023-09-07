import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import MilMoveErrorBoundary from './index';

jest.mock('utils/milmoveLog', () => ({
  milmoveLogger: {
    error: jest.fn(),
  },
}));

jest.mock('utils/retryPageLoading', () => ({
  retryPageLoading: jest.fn(),
}));

describe('MoveMoveErrorBoundary', () => {
  // to quiet the test, mock console.error
  beforeEach(() => {
    jest.spyOn(console, 'error').mockImplementation(() => null);
  });

  const Fallback = () => {
    return (
      <div>
        <span>My Fallback</span>
      </div>
    );
  };
  it('catches errors', async () => {
    const Thrower = () => {
      throw new Error('MyError');
    };
    render(
      <MilMoveErrorBoundary fallback={<Fallback />}>
        <Thrower />
      </MilMoveErrorBoundary>,
    );
    await waitFor(() => {
      expect(screen.getByText('My Fallback')).toBeVisible();
    });
  });

  it('shows children', () => {
    render(
      <MilMoveErrorBoundary fallback={<Fallback />}>
        <div>
          <span>All Good</span>
        </div>
      </MilMoveErrorBoundary>,
    );
    expect(screen.getByText('All Good')).toBeVisible();
  });
});
