import React from 'react';
import { render, screen } from '@testing-library/react';

import PlainTextDetails from './PlainTextDetails';

describe('PlainTextDetails', () => {
  it.each([
    ['approveShipment', 'Approved shipment', {}, {}],
    ['requestShipmentDiversion', 'Requested diversion', {}, {}],
    ['updateMTOServiceItemStatus', 'Service item status', {}, {}],
    ['setFinancialReviewFlag', 'Move flagged for financial review', {}, { financial_review_flag: 'true' }],
    ['setFinancialReviewFlag', 'Move unflagged for financial review', {}, { financial_review_flag: 'false' }],
    ['requestShipmentCancellation', 'Requested cancellation for HHG Shipment', { shipment_type: 'HHG' }, {}],
    ['updateMoveTaskOrderStatus', 'Created Move Task Order (MTO)', {}, { status: 'APPROVED' }],
    ['updateMoveTaskOrderStatus', 'Rejected Move Task Order (MTO)', {}, { status: 'REJECTED' }],
  ])('for event name %s it renders %s', (eventName, text, oldValues, changedValues) => {
    render(<PlainTextDetails eventName={eventName} oldValues={oldValues} changedValues={changedValues} />);

    expect(screen.getByText(text, { exact: false })).toBeInTheDocument();
  });
});
