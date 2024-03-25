import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';

import HHGWeightSummary from './HHGWeightSummary';

beforeEach(() => {
  jest.clearAllMocks();
});

const mtoShipments = [
  {
    shipmentType: 'HHG',
    primeEstimatedWeight: 110,
    primeActualWeight: 100,
  },
  {
    shipmentType: 'HHG',
    primeEstimatedWeight: 201,
    primeActualWeight: 200,
  },
];

describe('HHGWeightSummary component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(<HHGWeightSummary mtoShipments={mtoShipments} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'HHG 1' })).toBeInTheDocument();
        expect(screen.getByRole('heading', { level: 3, name: 'HHG 2' })).toBeInTheDocument();
      });

      expect(screen.getAllByText('Estimated Weight')[0]).toBeInTheDocument();
      expect(screen.getAllByText('Estimated Weight')[1]).toBeInTheDocument();
      expect(screen.getAllByText('Actual Weight')[0]).toBeInTheDocument();
      expect(screen.getAllByText('Actual Weight')[1]).toBeInTheDocument();
      expect(screen.getByText('110 lbs')).toBeInTheDocument();
      expect(screen.getByText('100 lbs')).toBeInTheDocument();
      expect(screen.getByText('201 lbs')).toBeInTheDocument();
      expect(screen.getByText('200 lbs')).toBeInTheDocument();
    });
  });
});
