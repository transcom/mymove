import React from 'react';
import { render, screen } from '@testing-library/react';

import SitCostBreakdown from './SitCostBreakdown';

import { MockProviders } from 'testUtils';
import { useGetPPMSITEstimatedCostQuery } from 'hooks/queries';
import { LOCATION_TYPES } from 'types/sitStatusShape';

const useGetPPMSITEstimatedCostQueryReturnValue = {
  estimatedCost: {
    sitCost: 5000,
    paramsAdditionalDaySIT: {
      contractYearName: 'Award Term 1',
      escalationCompounded: '1.10701',
      isPeak: 'true',
      numberDaysSIT: '6',
      priceRateOrFactor: '0.89',
      serviceAreaDestination: '456',
      serviceAreaOrigin: '',
    },
    paramsFirstDaySIT: {
      contractYearName: 'Award Term 1',
      escalationCompounded: '1.10701',
      isPeak: 'true',
      priceRateOrFactor: '20.35',
      serviceAreaDestination: '456',
      serviceAreaOrigi: '',
    },
    priceAdditionalDaySIT: 897,
    priceFirstDaySIT: 3402,
  },
  isError: false,
  isLoading: false,
  isSuccess: true,
};

const defaultProps = {
  ppmShipmentInfo: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    expectedDepartureDate: '2022-12-02',
    actualMoveDate: '2022-12-06',
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
    estimatedCost: 3000,
  },
  ppmSITLocation: LOCATION_TYPES.DESTINATION,
  sitStartDate: '2022-12-15',
  sitAdditionalStartDate: '2022-12-16',
  sitEndDate: '2022-12-25',
  weightStored: 2000,
  actualWeight: 3500,
  useQueries: useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue),
  setEstimatedCost: jest.fn(),
};

jest.mock('hooks/queries', () => ({
  useGetPPMSITEstimatedCostQuery: jest.fn(),
}));

describe('SitCostBreakdown component', () => {
  describe('when displayed', () => {
    it('shows sit cost components and their information', async () => {
      await useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);
      render(<SitCostBreakdown {...defaultProps} />, {
        wrapper: MockProviders,
      });
      expect(screen.getByText('Calculations')).toBeInTheDocument();
      expect(screen.getByText('SIT Information:')).toBeInTheDocument();
      expect(screen.getByText('Destination service area: 456')).toBeInTheDocument();
      expect(screen.getByText('Actual move date: 06-Dec-22')).toBeInTheDocument();
      expect(screen.getByText('Domestic peak')).toBeInTheDocument();
      expect(screen.getByText('Billable weight:')).toBeInTheDocument();
      expect(screen.getByText('20 cwt')).toBeInTheDocument();
      expect(screen.getByText('Adjusted weight: 2,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Estimated SIT weight: 0 lbs')).toBeInTheDocument();
      expect(screen.getByText('Actual PPM weight: 3,500 lbs')).toBeInTheDocument();
      expect(screen.getByText('Estimated PPM weight: 3,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('First day SIT price:')).toBeInTheDocument();
      expect(screen.getByText('$34.02')).toBeInTheDocument();
      expect(screen.getByText('SIT start date: 15-Dec-22')).toBeInTheDocument();
      expect(screen.getByText('Base price: $20.35/cwt')).toBeInTheDocument();
      expect(screen.getByText('Additional Day SIT price:')).toBeInTheDocument();
      expect(screen.getByText('$8.97')).toBeInTheDocument();
      expect(screen.getByText("SIT add'l day start: 16-Dec-22")).toBeInTheDocument();
      expect(screen.getByText('SIT end date: 25-Dec-22')).toBeInTheDocument();
      expect(screen.getByText('Additional days used: 6')).toBeInTheDocument();
      expect(screen.getByText('Price per day: $0.89/cwt')).toBeInTheDocument();
      expect(screen.getByText('Price escalation factor:')).toBeInTheDocument();
      expect(screen.getByText('1.10701')).toBeInTheDocument();
      expect(screen.getByText('Base year: Award Term 1')).toBeInTheDocument();
      expect(screen.getByText('Total:')).toBeInTheDocument();
      expect(screen.getByText('$50.00')).toBeInTheDocument();
    });
  });
});
