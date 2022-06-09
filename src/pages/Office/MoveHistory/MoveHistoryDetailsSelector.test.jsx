import React from 'react';
import { render, screen } from '@testing-library/react';

import MoveHistoryDetailsSelector from './MoveHistoryDetailsSelector';

describe('MoveHistoryDetailsSelector', () => {
  describe('for a plain text details event  (request shipment cancellation)', () => {
    const historyRecord = {
      action: 'UPDATE',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'requestShipmentCancellation',
      oldValues: { shipment_type: 'PPM' },
      tableName: '',
    };
    it('renders the plain text details appropriately', () => {
      render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);
      expect(screen.getByText('Requested cancellation for PPM shipment')).toBeInTheDocument();
    });
  });

  describe('for a labeled details event (update move task order)', () => {
    const historyRecord = {
      action: 'UPDATE',
      changedValues: {
        billable_weight_cap: '200',
        customer_remarks: 'Test customer remarks',
        counselor_remarks: '',
      },
      eventName: 'updateMoveTaskOrder',
      oldValues: { billable_: 'PPM' },
      tableName: 'moves',
    };
    it('renders the labeled details appropriately', () => {
      render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);
      expect(screen.getByText('Billable weight')).toBeInTheDocument();
      expect(screen.getByText(200, { exact: false })).toBeInTheDocument();
      expect(screen.getByText('Customer remarks')).toBeInTheDocument();
      expect(screen.getByText('Test customer remarks', { exact: false })).toBeInTheDocument();
    });
  });

  describe('handle a payments details update event', () => {
    const historyRecord = {
      action: 'UPDATE',
      changedValues: { status: 'REVIEWED', reviewed_at: '2022-04-27T20:56:24.867071' },
      eventName: 'updatePaymentRequestStatus',
      context: [
        {
          name: 'Test Service',
          price: '10123',
          status: 'APPROVED',
        },
        {
          name: 'Domestic uncrating',
          price: '5555',
          status: 'APPROVED',
        },
      ],
      tableName: 'payment_requests',
    };
    it('renders the payment details appropriately', () => {
      render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);
      expect(screen.getByText('Test Service')).toBeInTheDocument();
      expect(screen.getByText(101.23, { exact: false })).toBeInTheDocument();
      expect(screen.getByText('Domestic uncrating')).toBeInTheDocument();
      expect(screen.getByText(55.55, { exact: false })).toBeInTheDocument();
    });
  });

  describe('handle a payments details insert event', () => {
    const historyRecord = {
      action: 'INSERT',
      eventName: 'updateReweigh',
      tableName: 'payment_requests',
    };
    it('renders the payment details appropriately', () => {
      render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);
      expect(screen.getByText('Pending', { exact: false })).toBeInTheDocument();
    });
  });

  describe('handle a create payment request insert event', () => {
    const historyRecord = {
      action: 'INSERT',
      eventName: 'createPaymentRequest',
      tableName: 'payment_requests',
      context: [
        {
          name: 'Test service',
          price: '10123',
          status: 'REQUESTED',
          shipment_id: '123',
          shipment_type: 'HHG',
        },
        {
          name: 'Domestic uncrating',
          price: '5555',
          status: 'REQUESTED',
          shipment_id: '456',
          shipment_type: 'HHG_INTO_NTS_DOMESTIC',
        },
        { name: 'Move management', price: '1234', status: 'REQUESTED' },
      ],
    };
    it('renders the payment request appropriately', () => {
      render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);
      expect(screen.getByText('Move services')).toBeInTheDocument();
    });
  });
});
