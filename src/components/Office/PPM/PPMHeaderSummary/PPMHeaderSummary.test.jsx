import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';

import PPMHeaderSummary from './PPMHeaderSummary';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  ppmShipmentInfo: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    expectedDepartureDate: '2022-12-02',
    actualMoveDate: '2022-12-06',
    pickupAddress: {
      streetAddress1: '812 S 129th St',
      streetAddress2: '#123',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78234',
    },
    destinationAddress: {
      streetAddress1: '456 Oak Ln.',
      streetAddress2: '#123',
      city: 'Oakland',
      state: 'CA',
      postalCode: '94611',
    },
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
  },
  ppmNumber: 1,
  showAllFields: false,
};

describe('PPMHeaderSummary component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<PPMHeaderSummary {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'PPM 1' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByTestId('showRequestDetailsButton'));
      await waitFor(() => {
        expect(screen.getByText('Hide Details', { exact: false })).toBeInTheDocument();
      });
      expect(screen.getByText('Planned Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('02-Dec-2022')).toBeInTheDocument();
      expect(screen.getByText('Actual Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('06-Dec-2022')).toBeInTheDocument();
      expect(screen.getByText('Starting Address')).toBeInTheDocument();
      expect(screen.getByText('812 S 129th St, #123, San Antonio, TX 78234')).toBeInTheDocument();
      expect(screen.getByText('Ending Address')).toBeInTheDocument();
      expect(screen.getByText('456 Oak Ln., #123, Oakland, CA 94611')).toBeInTheDocument();
      expect(screen.getByText('Miles')).toBeInTheDocument();
      expect(screen.getByText('300')).toBeInTheDocument();
      expect(screen.getByText('Estimated Net Weight')).toBeInTheDocument();
      expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Actual Net Weight')).toBeInTheDocument();
      expect(screen.getByText('3,500 lbs')).toBeInTheDocument();
    });
  });
});
