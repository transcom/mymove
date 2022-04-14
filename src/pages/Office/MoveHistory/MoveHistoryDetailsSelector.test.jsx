import React from 'react';
import { render, screen } from '@testing-library/react';

import MoveHistoryDetailsSelector from './MoveHistoryDetailsSelector';

describe('MoveHistoryDetailsSelector', () => {
  describe('handle a plain text details event', () => {
    const historyRecord = {
      action: 'UPDATE',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'requestShipmentCancellation',
      oldValues: { shipment_type: 'PPM' },
      tableName: '',
    };
    render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);
    it('renders a plain text details event, request shipment cancellation', () => {
      expect(screen.getByText('Requested cancellation for PPM shipment')).toBeInTheDocument();
    });
  });

  // describe('handle a labelled details event', () => {

  // });

  // describe('handle a payments details event', () => {

  // });
});
