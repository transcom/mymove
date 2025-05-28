import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import SitCost from './SitCost';

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
  sitEndDate: '2022-12-25',
  weightStored: 2000,
  actualWeight: 3500,
  useQueries: useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue),
  setEstimatedCost: jest.fn(),
};

jest.mock('hooks/queries', () => ({
  useGetPPMSITEstimatedCostQuery: jest.fn(),
}));

describe('SitCost component', () => {
  describe('when displayed', () => {
    it('toggles ppm sit calculations when button is clicked', async () => {
      await useGetPPMSITEstimatedCostQuery.mockReturnValue(useGetPPMSITEstimatedCostQueryReturnValue);
      render(<SitCost {...defaultProps} />, {
        wrapper: MockProviders,
      });
      const toggleButtonShow = screen.getByText('Show calculations');
      expect(toggleButtonShow).toBeInTheDocument();

      await userEvent.click(toggleButtonShow);
      const toggleButtonHide = screen.getByText('Hide calculations');
      expect(toggleButtonHide).toBeInTheDocument();

      await userEvent.click(toggleButtonHide);
      expect(toggleButtonShow).toBeInTheDocument();
      await userEvent.click(toggleButtonShow);

      expect(screen.getByText('Calculations')).toBeInTheDocument();
      expect(screen.getByText('Total:')).toBeInTheDocument();
    });
  });
});
