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
    ['requestShipmentDiversion', 'Requested diversion'],
    ['approveShipment', 'Approved shipment'],
    ['updateMTOServiceItem', 'Service Items'],
    ['updateMTOServiceItemStatus', 'Service item status'],
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

    ['setFinancialReviewFlag', 'Move flagged for financial review', {}, { financial_review_flag: 'true' }],
    ['setFinancialReviewFlag', 'Move unflagged for financial review', {}, { financial_review_flag: 'false' }],
    ['updateMoveTaskOrderStatus', 'Created Move Task Order (MTO)', {}, { status: 'APPROVED' }],
    ['updateMoveTaskOrderStatus', 'Rejected Move Task Order (MTO)', {}, { status: 'Rejected' }],
    ['approveShipment', 'HHG shipment', { status: 'APPROVED' }, { shipment_type: 'HHG' }],
    ['requestShipmentCancellation', 'Requested cancellation for HHG shipment', { shipment_type: 'HHG' }, {}],
    ['requestShipmentDiversion', 'Requested diversion for HHG shipment', { shipment_type: 'HHG' }, {}],
  ])('for event name %s it renders %s', (eventName, text, oldValues, changedValues) => {
    render(<MoveHistoryDetailsSelector eventName={eventName} oldValues={oldValues} changedValues={changedValues} />);

    expect(screen.getByText(text)).toBeInTheDocument();
  });
});
