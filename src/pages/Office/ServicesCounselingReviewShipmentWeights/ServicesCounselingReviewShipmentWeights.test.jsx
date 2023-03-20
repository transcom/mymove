import React from 'react';
import { render, screen, within } from '@testing-library/react';

import {
  missingSomeWeightQuery,
  riskOfExcessWeightQuery,
  reviewWeightsQuery,
} from '../MoveTaskOrder/moveTaskOrderUnitTestData';

import ServicesCounselingReviewShipmentWeights from './ServicesCounselingReviewShipmentWeights';

import { useReviewShipmentWeightsQuery } from 'hooks/queries';

jest.mock('hooks/queries', () => ({
  useReviewShipmentWeightsQuery: jest.fn(),
}));

describe('Services Counseling Review Shipment Weights', () => {
  describe('basic rendering', () => {
    it('should render the review shipment weights page', () => {
      useReviewShipmentWeightsQuery.mockReturnValue(missingSomeWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="ABC123" />);

      expect(screen.getByRole('heading', { name: 'Review shipment weights', level: 1 })).toBeInTheDocument();
    });

    it('displays the weight allowance', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(missingSomeWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="ABC123" />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const weightAllowanceDisplay = weightDisplays[0];
      expect(weightAllowanceDisplay).toHaveTextContent('8,500 lbs');
    });

    it('displays the total estimated weight', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(missingSomeWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="ABC123" />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const estimatedWeightDisplay = weightDisplays[1];
      expect(estimatedWeightDisplay).toHaveTextContent('125 lbs');
    });

    it('displays the max billable weight', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(missingSomeWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="ABC123" />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const maxBillableWeightDisplay = weightDisplays[2];
      expect(maxBillableWeightDisplay).toHaveTextContent('8,000 lbs');
    });

    it('displays the total move weight', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(missingSomeWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="ABC123" />);

      const weightDisplays = await screen.findAllByTestId('weight-display');
      const totalMoveWeight = weightDisplays[3];
      expect(totalMoveWeight).toHaveTextContent('125 lbs');
    });

    it('displays risk of excess tag', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(riskOfExcessWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="ABC123" />);

      const riskOfExcessTag = screen.getByText(/Risk of excess/);
      expect(riskOfExcessTag).toBeInTheDocument();
    });
    it('displays PPM shipments weights list', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);
      await render(<ServicesCounselingReviewShipmentWeights moveCode="XSWT05" />);
      const container = await screen.findByTestId('ppmShipmentContainer');
      expect(container).toBeInTheDocument();
      const table = await within(container).getByRole('table');
      expect(table).toBeInTheDocument();
      expect(screen.getByText('Weight moved by customer')).toBeInTheDocument();
    });
    it('displays pro-gear weights', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);
      await render(<ServicesCounselingReviewShipmentWeights moveCode="XSWT05" />);
      const container = await screen.findByTestId('progearContainer');
      expect(container).toBeInTheDocument();
      const table = await within(container).getByRole('table');
      expect(table).toBeInTheDocument();
      expect(screen.getByText('Weight moved')).toBeInTheDocument();
    });
    it('displays non-PPM shipments weights list', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);
      await render(<ServicesCounselingReviewShipmentWeights moveCode="XSWT05" />);
      const container = await screen.findByTestId('nonPpmShipmentContainer');
      expect(container).toBeInTheDocument();
      const table = await within(container).getByRole('table');
      expect(table).toBeInTheDocument();
      expect(screen.getByText('Shipments')).toBeInTheDocument();
    });
  });
});
