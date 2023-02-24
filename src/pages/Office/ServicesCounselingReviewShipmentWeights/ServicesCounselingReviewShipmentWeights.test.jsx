import React from 'react';
import { render, screen } from '@testing-library/react';

import { missingSomeWeightQuery } from '../MoveTaskOrder/moveTaskOrderUnitTestData';

import ServicesCounselingReviewShipmentWeights from './ServicesCounselingReviewShipmentWeights';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';

jest.mock('hooks/queries', () => ({
  useReviewShipmentWeightsQuery: jest.fn(),
}));

describe('Services Counseling Review Shipment Weights', () => {
  describe('basic rendering', () => {
    it('should render the review shipment weights page', () => {
      useReviewShipmentWeightsQuery.mockReturnValue(missingSomeWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights />);

      expect(screen.getByRole('heading', { name: 'Review shipment weights', level: 1 })).toBeInTheDocument();
    });

    it('displays the weight allowance', async () => {
      render(<ServicesCounselingReviewShipmentWeights />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const weightAllowanceDisplay = weightDisplays[0];
      expect(weightAllowanceDisplay).toHaveTextContent('8,500 lbs');
    });

    it('displays the total estimated weight', async () => {
      render(<ServicesCounselingReviewShipmentWeights />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const estimatedWeightDisplay = weightDisplays[1];
      expect(estimatedWeightDisplay).toHaveTextContent('125 lbs');
    });

    it('displays the max billable weight', async () => {
      render(<ServicesCounselingReviewShipmentWeights />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const maxBillableWeightDisplay = weightDisplays[2];
      expect(maxBillableWeightDisplay).toHaveTextContent('8,000 lbs');
    });

    it('displays the total move weight', async () => {
      render(<ServicesCounselingReviewShipmentWeights />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const totalMoveWeight = weightDisplays[3];
      expect(totalMoveWeight).toHaveTextContent('125 lbs');
    });
  });
});
