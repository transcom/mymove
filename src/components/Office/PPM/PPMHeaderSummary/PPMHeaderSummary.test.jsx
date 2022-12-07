import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';

import PPMHeaderSummary from './PPMHeaderSummary';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  ppmShipment: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    actualMoveDate: Date.now(),
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
    hasReceivedAdvance: true,
    advanceAmountReceived: 60000,
  },
  ppmNumber: 1,
};

describe('PPMHeaderSummary component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<PPMHeaderSummary {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'PPM 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Departure date')).toBeInTheDocument();
      expect(screen.getByText('06-Dec-2022')).toBeInTheDocument();
      expect(screen.getByText('Starting ZIP')).toBeInTheDocument();
      expect(screen.getByText('90210')).toBeInTheDocument();
      expect(screen.getByText('94611')).toBeInTheDocument();
      expect(screen.getByText('Advance received')).toBeInTheDocument();
      expect(screen.getByText('Yes, $600')).toBeInTheDocument();
    });
  });
});
