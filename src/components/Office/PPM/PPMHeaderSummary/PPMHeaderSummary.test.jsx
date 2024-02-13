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
    actualPickupPostalCode: '90210',
    actualDestinationPostalCode: '94611',
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
      expect(screen.getByText('Starting ZIP')).toBeInTheDocument();
      expect(screen.getByText('90210')).toBeInTheDocument();
      expect(screen.getByText('Ending ZIP')).toBeInTheDocument();
      expect(screen.getByText('94611')).toBeInTheDocument();
      expect(screen.getByText('Miles')).toBeInTheDocument();
      expect(screen.getByText('300')).toBeInTheDocument();
      expect(screen.getByText('Estimated Net Weight')).toBeInTheDocument();
      expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Actual Net Weight')).toBeInTheDocument();
      expect(screen.getByText('3,500 lbs')).toBeInTheDocument();
    });
  });
});
