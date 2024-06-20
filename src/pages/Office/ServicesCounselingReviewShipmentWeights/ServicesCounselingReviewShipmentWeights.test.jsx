import React from 'react';
import { render, screen, within } from '@testing-library/react';

import {
  missingSomeWeightQuery,
  riskOfExcessWeightQuery,
  reviewWeightsQuery,
  reviewWeightsNoProGearQuery,
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
      render(<ServicesCounselingReviewShipmentWeights moveCode="XSWT05" />);
      const container = await screen.findByTestId('ppmShipmentContainer');
      expect(container).toBeInTheDocument();
      const table = await within(container).getByRole('table');
      expect(table).toBeInTheDocument();
      expect(screen.getByText('Weight moved by customer')).toBeInTheDocument();
    });

    it('displays non-PPM shipments weights list', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="XSWT05" />);
      const container = await screen.findByTestId('nonPpmShipmentContainer');
      expect(container).toBeInTheDocument();
      const table = await within(container).getByRole('table');
      expect(table).toBeInTheDocument();
      expect(screen.getByText('Shipments')).toBeInTheDocument();
    });

    it('displays non-PPM shipments weights section when there are shipments without pro-gear', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsNoProGearQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="XSWT05" />);
      const container = await screen.findByTestId('nonPpmShipmentContainer');
      expect(container).toBeInTheDocument();
      const table = await within(container).getByRole('table');
      expect(table).toBeInTheDocument();
      const progear = await screen.queryByTestId('progearContainer');
      expect(progear).not.toBeInTheDocument();
      expect(screen.getByText('Shipments')).toBeInTheDocument();
    });

    it('displays excess weight warning when move has excess weight', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(reviewWeightsQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="XSWT01" />);

      const excessWeightWarning = await screen.findByTestId('alert');
      expect(excessWeightWarning).toBeInTheDocument();
      expect(excessWeightWarning).toHaveTextContent(
        'This move has excess weight. Review PPM weight ticket documents to resolve.',
      );
    });

    it('does NOT display excess weight warning when move does NOT have excess weight', async () => {
      useReviewShipmentWeightsQuery.mockReturnValue(missingSomeWeightQuery);
      render(<ServicesCounselingReviewShipmentWeights moveCode="CLOSE0" />);

      const excessWeightWarning = await screen.queryByTestId('alert');
      expect(excessWeightWarning).not.toBeInTheDocument();
    });
  });
});
