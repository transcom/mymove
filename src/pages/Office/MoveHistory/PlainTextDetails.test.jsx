import React from 'react';
import { render, screen } from '@testing-library/react';

import PlainTextDetails from './PlainTextDetails';

import { shipmentOptionToDisplay } from 'constants/historyLogUIDisplayName';

describe('PlainTextDetails', () => {
  it.each([
    [{ eventName: 'approveShipment', oldValues: { shipmentType: 'HHG' } }, 'HHG shipment'],
    [
      { eventName: 'requestShipmentDiversion', oldValues: { shipment_type: 'HHG' } },
      `Requested diversion for ${shipmentOptionToDisplay.HHG} shipment`,
    ],
    [
      { eventName: 'setFinancialReviewFlag', changedValues: { financial_review_flag: 'true' } },
      'Move flagged for financial review',
    ],
    [
      { eventName: 'setFinancialReviewFlag', changedValues: { financial_review_flag: 'false' } },
      'Move unflagged for financial review',
      {},
    ],
    [
      {
        eventName: 'requestShipmentCancellation',
        oldValues: { shipment_type: 'HHG' },
      },
      `Requested cancellation for ${shipmentOptionToDisplay.HHG} shipment`,
    ],
    [
      { eventName: 'updateMoveTaskOrderStatus', changedValues: { status: 'APPROVED' } },
      'Created Move Task Order (MTO)',
    ],
    [
      { eventName: 'updateMoveTaskOrderStatus', changedValues: { status: 'REJECTED' } },
      'Rejected Move Task Order (MTO)',
    ],
    [
      {
        eventName: 'updateMTOServiceItemStatus',
        context: { name: 'Domestic origin price', shipment_type: 'HHG_INTO_NTS_DOMESTIC' },
      },
      'NTS shipment, Domestic origin price',
    ],
  ])('for history record %s it renders %s', (historyRecord, text) => {
    render(<PlainTextDetails historyRecord={historyRecord} />);

    expect(screen.getByText(text)).toBeInTheDocument();
  });
});
