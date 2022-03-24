import React from 'react';
import { render, screen } from '@testing-library/react';

import MoveHistoryDetailsSelector from './MoveHistoryDetailsSelector';

describe('MoveHistoryDetailsSelector', () => {
  it.each([
    ['counselingUpdateOrder', 'Labelled'],
    ['updateOrder', 'Labelled'],
    ['updateAllowance', 'Labelled'],
    ['counselingUpdateAllowance', 'Labelled'],
    ['updateMoveTaskOrder', 'Labelled'],
    ['updateMTOShipment', 'Labelled'],
    ['approveShipment', 'Approved shipment'],
    ['requestShipmentDiversion', 'Requested diversion'],
    ['updateMTOServiceItem', 'Service Items'],
    ['updateMTOServiceItemStatus', 'Service item status'],
    ['requestShipmentCancellation', 'Shipment cancelled'],
    ['createOrders', '-'],
    ['updateOrders', 'Labelled'],
    ['uploadAmendedOrders', '-'],
    ['submitMoveForApproval', '-'],
    ['submitAmendedOrders', 'Labelled'],
    ['createMTOShipment', '-'],
    ['updateMTOShipmentAddress', 'Labelled'],
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
    ['updateMoveTaskOrderStatus', 'Created Move Task Order (MTO)', [{ columnName: 'status', columnValue: 'APPROVED' }]],
  ])('for event name %s it renders %s', (eventName, text, changedValues) => {
    render(<MoveHistoryDetailsSelector eventName={eventName} changedValues={changedValues} />);

    expect(screen.getByText(text, { exact: false })).toBeInTheDocument();
  });
});
