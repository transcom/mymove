import React from 'react';
import { render, screen } from '@testing-library/react';

import PlainTextDetails from './PlainTextDetails';

import { shipmentOptionToDisplay } from 'constants/historyLogUIDisplayName';

describe('PlainTextDetails', () => {
  it.each([
    ['approveShipment', 'Approved shipment', {}, {}],
    ['requestShipmentDiversion', 'Requested diversion', {}, {}],
    ['updateMTOServiceItemStatus', 'Service item status', {}, {}],
    ['setFinancialReviewFlag', 'Move flagged for financial review', {}, { financial_review_flag: 'true' }],
    ['setFinancialReviewFlag', 'Move unflagged for financial review', {}, { financial_review_flag: 'false' }],
    [
      'requestShipmentCancellation',
      `Requested cancellation for ${shipmentOptionToDisplay.HHG} shipment`,
      { shipment_type: 'HHG' },
      {},
    ],
    [
      'requestShipmentCancellation',
      `Requested cancellation for ${shipmentOptionToDisplay.HHG_OUTOF_NTS_DOMESTIC} shipment`,
      { shipment_type: 'HHG_OUTOF_NTS_DOMESTIC' },
      {},
    ],
    [
      'requestShipmentCancellation',
      `Requested cancellation for ${shipmentOptionToDisplay.HHG_INTO_NTS_DOMESTIC} shipment`,
      { shipment_type: 'HHG_INTO_NTS_DOMESTIC' },
      {},
    ],
    [
      'requestShipmentCancellation',
      `Requested cancellation for ${shipmentOptionToDisplay.PPM} shipment`,
      { shipment_type: 'PPM' },
      {},
    ],
    [
      'requestShipmentCancellation',
      `Requested cancellation for ${shipmentOptionToDisplay.HHG_SHORTHAUL_DOMESTIC} shipment`,
      { shipment_type: 'HHG_SHORTHAUL_DOMESTIC' },
      {},
    ],
    ['updateMoveTaskOrderStatus', 'Created Move Task Order (MTO)', {}, { status: 'APPROVED' }],
    ['updateMoveTaskOrderStatus', 'Rejected Move Task Order (MTO)', {}, { status: 'REJECTED' }],
  ])('for event name %s it renders %s', (eventName, text, oldValues, changedValues) => {
    render(<PlainTextDetails eventName={eventName} oldValues={oldValues} changedValues={changedValues} />);

    expect(screen.getByText(text, { exact: false })).toBeInTheDocument();
  });
});
