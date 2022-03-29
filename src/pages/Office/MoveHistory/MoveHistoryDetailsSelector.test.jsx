import React from 'react';
import { render, screen } from '@testing-library/react';

import MoveHistoryDetailsSelector from './MoveHistoryDetailsSelector';

describe('MoveHistoryDetailsSelector', () => {
  it.each([
    ['counselingUpdateOrder', 'Labeled'],
    ['updateOrder', 'Labeled'],
    ['updateAllowance', 'Labeled'],
    ['counselingUpdateAllowance', 'Labeled'],
    ['updateMoveTaskOrder', 'Labeled'],
    ['updateMTOShipment', 'Labeled'],
    ['approveShipment', 'Approved shipment'],
    ['requestShipmentDiversion', 'Requested diversion'],
    ['updateMTOServiceItem', 'Service Items'],
    ['updateMTOServiceItemStatus', 'Service item status'],
    ['requestShipmentCancellation', 'Shipment cancelled'],
    ['createOrders', '-'],
    ['updateOrders', 'Labeled'],
    ['uploadAmendedOrders', '-'],
    ['submitMoveForApproval', '-'],
    ['submitAmendedOrders', 'Labeled'],
    ['createMTOShipment', '-'],
    ['updateMTOShipmentAddress', 'Labeled'],
    ['createMTOServiceItem', 'Service Items'],
    ['default', '-'],
  ])('for event name %s it renders %s', (eventName, text) => {
    render(<MoveHistoryDetailsSelector eventName={eventName} />);

    expect(screen.getByText(text, { exact: false })).toBeInTheDocument();
  });

  it.each([
    [
      'setFinancialReviewFlag',
      'Move flagged for financial review',
      [{ columnName: 'financial_review_flag', columnValue: 'true' }],
    ],
    [
      'setFinancialReviewFlag',
      'Move unflagged for financial review',
      [{ columnName: 'financial_review_flag', columnValue: 'false' }],
    ],
    ['updateMoveTaskOrderStatus', 'Created Move Task Order (MTO)', [{ columnName: 'status', columnValue: 'APPROVED' }]],
    [
      'updateMoveTaskOrderStatus',
      'Rejected Move Task Order (MTO)',
      [{ columnName: 'status', columnValue: 'Rejected' }],
    ],
  ])('for event name %s it renders %s', (eventName, text, changedValues) => {
    render(<MoveHistoryDetailsSelector eventName={eventName} changedValues={changedValues} />);

    expect(screen.getByText(text, { exact: false })).toBeInTheDocument();
  });
});
