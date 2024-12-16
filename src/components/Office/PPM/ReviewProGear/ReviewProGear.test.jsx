import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';

import ReviewProGear from './ReviewProGear';

import ppmDocumentStatus from 'constants/ppms';
import { MockProviders } from 'testUtils';

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
  tripNumber: 1,
  ppmNumber: '1',
  showAllFields: false,
};

const proGearRequiredProps = {
  proGear: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    belongsToSelf: true,
    weight: 400,
    description: 'Kia Forte',
    hasWeightTickets: true,
  },
};

const missingWeightTicketProps = {
  proGear: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    belongsToSelf: true,
    weight: 400,
    description: 'Kia Forte',
    hasWeightTickets: false,
  },
};

const rejectedProps = {
  proGear: {
    ...proGearRequiredProps.proGear,
    status: ppmDocumentStatus.REJECTED,
    reason: 'Rejection reason',
  },
};

describe('ReviewProGear component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} />;
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Pro-gear 1' })).toBeInTheDocument();
      });
      expect(screen.getByText('Belongs to')).toBeInTheDocument();
      expect(screen.getByLabelText('Customer')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Spouse')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByText('Description')).toBeInTheDocument();

      expect(screen.getByText('Pro-gear weight')).toBeInTheDocument();
      expect(screen.getByLabelText('Weight tickets')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Constructed weight')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByLabelText('Constructed pro-gear weight')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: 'Review pro-gear 1' })).toBeInTheDocument();
      expect(screen.getByText('Add a review for this pro-gear')).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing pro-gear weight ticket values', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} {...proGearRequiredProps} />;
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByLabelText('Customer')).toBeChecked();
      });
      expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      expect(screen.getByLabelText(/Shipment's pro-gear weight/)).toHaveDisplayValue('400');
    });

    it('populates edit form when pro-gear weight ticket is missing', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} {...missingWeightTicketProps} />;
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Constructed weight')).toBeChecked();
        expect(screen.getByText('Constructed pro-gear weight')).toBeInTheDocument();
      });
    });

    it('displays remaining character count', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} {...rejectedProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reason')).toHaveDisplayValue('Rejection reason');
      });
      expect(screen.getByText('484 characters')).toBeInTheDocument();
    });

    it('toggles the reason field when Reject is selected', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} {...proGearRequiredProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
      });
      await fireEvent.click(screen.getByLabelText('Reject'));
      expect(screen.getByLabelText('Reason')).toBeInstanceOf(HTMLTextAreaElement);
      await fireEvent.click(screen.getByLabelText('Accept'));
      expect(screen.queryByLabelText('Reason')).not.toBeInTheDocument();
    });
  });

  describe('displays disabled read only form', () => {
    it('renders disabled blank form on load with defaults', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} readOnly />;
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Pro-gear 1' })).toBeInTheDocument();
      });
      expect(screen.getByText('Belongs to')).toBeInTheDocument();
      expect(screen.getByLabelText('Customer')).toBeDisabled();
      expect(screen.getByLabelText('Spouse')).toBeDisabled();

      expect(screen.getByText('Description')).toBeInTheDocument();

      expect(screen.getByText('Pro-gear weight')).toBeInTheDocument();
      expect(screen.getByLabelText('Weight tickets')).toBeDisabled();
      expect(screen.getByLabelText('Constructed weight')).toBeDisabled();

      expect(screen.getByLabelText('Constructed pro-gear weight')).toBeDisabled();

      expect(screen.getByRole('heading', { level: 3, name: 'Review pro-gear 1' })).toBeInTheDocument();
      expect(screen.getByText('Add a review for this pro-gear')).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toBeDisabled();
      expect(screen.getByLabelText('Reject')).toBeDisabled();
    });

    it('populates disabled edit form with existing pro-gear weight ticket values', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} {...proGearRequiredProps} readOnly />;
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByLabelText('Customer')).toBeChecked();
      });
      expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      expect(screen.getByLabelText(/Shipment's pro-gear weight/)).toHaveDisplayValue('400');
      expect(screen.getByLabelText(/Shipment's pro-gear weight/)).toBeDisabled();
    });

    it('populates disabled edit form when pro-gear weight ticket is missing', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} {...missingWeightTicketProps} readOnly />;
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Constructed weight')).toBeDisabled();
      });
    });

    it('reason field is disabled', async () => {
      render(
        <MockProviders>
          <ReviewProGear {...defaultProps} {...rejectedProps} readOnly />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeDisabled();
      });
      expect(screen.getByLabelText('Reason')).toBeDisabled();
    });
  });
});
