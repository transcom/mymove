import React from 'react';
import { render, screen } from '@testing-library/react';

import MoveHistoryDetailsSelector from './MoveHistoryDetailsSelector';

describe('MoveHistoryDetailsSelector', () => {
  it.each([
    [{ eventName: 'counselingUpdateOrder' }, 'Labeled'],
    [{ eventName: 'updateOrder' }, 'Labeled'],
    [{ eventName: 'updateAllowance' }, 'Labeled'],
    [{ eventName: 'counselingUpdateAllowance' }, 'Labeled'],
    [{ eventName: 'updateMoveTaskOrder' }, 'Labeled'],
    [{ eventName: 'updateMTOShipment' }, 'Labeled'],
	[{ eventName: 'requestShipmentDiversion' }, 'Requested diversion'],
    [{ eventName: 'approveShipment' }, 'Approved shipment'],
    [{ eventName: 'updateMTOServiceItem' }, 'Service Items'],
    // [
    //   {
    //     eventName: 'updateMTOServiceItemStatus',
    //     context: { name: 'Domestic origin price', shipment_type: 'HHG_INTO_NTS_DOMESTIC' },
    //   },
    //   'Service item status',
    // ],
    [{ eventName: 'createOrders' }, '-'],
    [{ eventName: 'updateOrders' }, 'Labeled'],
    [{ eventName: 'uploadAmendedOrders' }, '-'],
    [{ eventName: 'submitMoveForApproval' }, '-'],
    [{ eventName: 'submitAmendedOrders' }, 'Labeled'],
    [{ eventName: 'createMTOShipment' }, '-'],
    [{ eventName: 'updateMTOShipmentAddress' }, 'Labeled'],
    [{ eventName: 'createMTOServiceItem' }, 'Service Items'],
    [{ eventName: 'default' }, '-'],
  ])('for history record %s it renders %s', (historyRecord, text) => {
    render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);

    expect(screen.getByText(text, { exact: false })).toBeInTheDocument();
  });

  it.each([
    [
      { eventName: 'setFinancialReviewFlag', changedValues: { financial_review_flag: 'true' } },
      'Move flagged for financial review',
    ],
    [
      { eventName: 'setFinancialReviewFlag', changedValues: { financial_review_flag: 'false' } },
      'Move unflagged for financial review',
    ],
    [
      { eventName: 'updateMoveTaskOrderStatus', changedValues: { status: 'APPROVED' } },
      'Created Move Task Order (MTO)',
    ],
    [
      { eventName: 'updateMoveTaskOrderStatus', changedValues: { status: 'Rejected' } },
      'Rejected Move Task Order (MTO)',
    ],
	[
		{ eventName: 'approveShipmentDiversion', oldValues: { shipment_type: 'HHG' }, changedValues: { status: 'APPROVED' } },
		'HHG shipment',
	],
    [
      { eventName: 'requestShipmentCancellation', oldValues: { shipment_type: 'HHG' } },
      'Requested cancellation for HHG shipment',
    ],
    [
      { eventName: 'requestShipmentDiversion', oldValues: { shipment_type: 'HHG' } },
      'Requested diversion for HHG shipment',
    ],
  ])('for historyRecord %s it renders %s', (historyRecord, text) => {
    render(<MoveHistoryDetailsSelector historyRecord={historyRecord} />);

    expect(screen.getByText(text)).toBeInTheDocument();
  });
});
