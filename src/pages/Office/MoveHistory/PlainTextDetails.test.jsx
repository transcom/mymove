import React from 'react';
import { render, screen } from '@testing-library/react';

import PlainTextDetails from './PlainTextDetails';

describe('PlainTextDetails', () => {
  it.each([
    ['approveShipment', 'HHG shipment', {}, { shipment_type: 'HHG' }],
    ['requestShipmentDiversion', 'Requested diversion', {}, {}],
    ['updateMTOServiceItemStatus', 'Service item status', {}, {}],
    ['setFinancialReviewFlag', 'Move flagged for financial review', { financial_review_flag: 'true' }, {}],
    ['setFinancialReviewFlag', 'Move unflagged for financial review', { financial_review_flag: 'false' }, {}],
    ['requestShipmentCancellation', 'Shipment cancelled', {}, {}],
    ['updateMoveTaskOrderStatus', 'Created Move Task Order (MTO)', { status: 'APPROVED' }, {}],
    ['updateMoveTaskOrderStatus', 'Rejected Move Task Order (MTO)', { status: 'REJECTED' }, {}],
  ])('for event name %s it renders %s', (eventName, text, changedValues, oldValues) => {
    render(<PlainTextDetails eventName={eventName} changedValues={changedValues} oldValues={oldValues} />);

    expect(screen.getByText(text, { exact: false })).toBeInTheDocument();
  });
});
