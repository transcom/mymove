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
      expect(screen.getByText('Billable weight cap')).toBeInTheDocument();
      expect(screen.getByText(200, { exact: false })).toBeInTheDocument();
      expect(screen.getByText('Customer remarks')).toBeInTheDocument();
      expect(screen.getByText('Test customer remarks', { exact: false })).toBeInTheDocument();
    });
  });

  // describe('handle a payments details event', () => {

  // });
});
