import React from 'react';
import { render, waitFor, screen, fireEvent } from '@testing-library/react';

import ReviewGunSafe from './ReviewGunSafe';

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
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
  },
  tripNumber: 1,
  ppmNumber: '1',
  showAllFields: false,
};

const gunSafeRequiredProps = {
  gunSafe: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    weight: 500,
    description: 'Kia Forte',
    hasWeightTickets: true,
  },
};

const overTheWeightLimitProps = {
  gunSafe: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    weight: 501,
    description: 'Kia Forte',
    hasWeightTickets: true,
  },
};

const missingWeightTicketProps = {
  gunSafe: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    belongsToSelf: true,
    weight: 400,
    description: 'Kia Forte',
    hasWeightTickets: false,
  },
};

const rejectedProps = {
  gunSafe: {
    ...gunSafeRequiredProps.gunSafe,
    status: ppmDocumentStatus.REJECTED,
    reason: 'Rejection reason',
  },
};

describe('ReviewGunSafe component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} />;
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Gun safe 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByText('Gun safe weight')).toBeInTheDocument();
      expect(screen.getByLabelText('Weight tickets')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Constructed weight')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Constructed gun safe weight')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('heading', { level: 3, name: 'Review gun safe 1' })).toBeInTheDocument();
      expect(screen.getByText('Add a review for this gun safe')).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Reject')).toBeInstanceOf(HTMLInputElement);
    });

    it('populates edit form with existing gun safe weight ticket values', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...gunSafeRequiredProps} />;
        </MockProviders>,
      );

      expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      expect(screen.getByLabelText(/Shipment's gun safe weight/)).toHaveDisplayValue('500');
    });

    it('populates edit form when gun safe weight ticket is missing', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...missingWeightTicketProps} />;
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Constructed weight')).toBeChecked();
        expect(screen.getByText('Constructed gun safe weight')).toBeInTheDocument();
      });
    });

    it('displays remaining character count', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...rejectedProps} />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reason')).toHaveDisplayValue('Rejection reason');
      });
      expect(screen.getByText('484 characters')).toBeInTheDocument();
    });

    it('displays over the weight limit message for gun safe', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...overTheWeightLimitProps} />
        </MockProviders>,
      );
      expect(
        screen.getByText(
          'The government authorizes the shipment of a gun safe up to 500 lbs (This is not charged against the authorized weight entitlement. Any weight over 500 lbs is charged against the weight entitlement).',
        ),
      ).toBeInTheDocument();
    });

    it('toggles the reason field when Reject is selected', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...gunSafeRequiredProps} />
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
          <ReviewGunSafe {...defaultProps} readOnly />;
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'Gun safe 1' })).toBeInTheDocument();
      });

      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByText('Gun safe weight')).toBeInTheDocument();
      expect(screen.getByLabelText('Weight tickets')).toBeDisabled();
      expect(screen.getByLabelText('Constructed weight')).toBeDisabled();
      expect(screen.getByLabelText('Constructed gun safe weight')).toBeDisabled();

      expect(screen.getByRole('heading', { level: 3, name: 'Review gun safe 1' })).toBeInTheDocument();
      expect(screen.getByText('Add a review for this gun safe')).toBeInTheDocument();
      expect(screen.getByLabelText('Accept')).toBeDisabled();
      expect(screen.getByLabelText('Reject')).toBeDisabled();
    });

    it('populates disabled edit form with existing gun safe weight ticket values', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...gunSafeRequiredProps} readOnly />;
        </MockProviders>,
      );

      expect(screen.getByText('Kia Forte')).toBeInTheDocument();
      expect(screen.getByLabelText(/Shipment's gun safe weight/)).toHaveDisplayValue('500');
      expect(screen.getByLabelText(/Shipment's gun safe weight/)).toBeDisabled();
    });

    it('populates disabled edit form when gun safe weight ticket is missing', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...missingWeightTicketProps} readOnly />;
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Constructed weight')).toBeDisabled();
      });
    });

    it('reason field is disabled', async () => {
      render(
        <MockProviders>
          <ReviewGunSafe {...defaultProps} {...rejectedProps} readOnly />
        </MockProviders>,
      );
      await waitFor(() => {
        expect(screen.getByLabelText('Reject')).toBeDisabled();
      });
      expect(screen.getByLabelText('Reason')).toBeDisabled();
    });
  });
});
