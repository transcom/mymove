import React from 'react';
import { render, screen } from '@testing-library/react';

import LoadingSpinner from './LoadingSpinner';

describe('LoadingSpinner Component', () => {
  test('renders the loading spinner with default message', () => {
    render(<LoadingSpinner />);

    const spinner = screen.getByTestId('loading-spinner');
    expect(spinner).toBeInTheDocument();

    expect(screen.getByText('Loading, please wait...')).toBeInTheDocument();
  });

  test('renders the loading spinner with a custom message', () => {
    const customMessage = 'Fetching data...';
    render(<LoadingSpinner message={customMessage} />);

    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();

    expect(screen.getByText(customMessage)).toBeInTheDocument();
  });
});
