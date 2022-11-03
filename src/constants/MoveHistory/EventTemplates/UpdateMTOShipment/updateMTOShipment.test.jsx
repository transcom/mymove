import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipment';

describe('when given an mto shipment update with mto shipment table history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'updateMTOShipment',
    tableName: 'mto_shipments',
    changedValues: {
      destination_address_type: 'HOME_OF_SELECTION',
      requested_delivery_date: '2020-04-14',
      requested_pickup_date: '2020-03-23',
      actual_pickup_date: '2020-04-15',
      approved_date: '2020-04-21',
      billable_weight_cap: '10000',
      billable_weight_justification: 'Heavy items',
      counselor_remarks: 'Remarks',
      customer_remarks: 'Words',
      diversion: 'false',
      first_available_delivery_date: '2020-04-12',
      id: 'b1cb1428-6d65-40fc-addd-b57f4ea120f1',
      prime_actual_weight: '14000',
      prime_estimated_weight: '12000',
      sac_type: 'HHG',
      scheduled_delivery_date: '2020-04-13',
      scheduled_pickup_date: '2020-03-22',
      service_order_number: '767567576',
      status: 'SUBMITTED',
      tac_type: 'NTS',
      uses_external_vendor: 'true',
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_id_abbr: 'b4b4b',
      },
    ],
  };
  it('correctly matches the Update mto shipment event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
  });
  describe('it correctly renders the details component for Create MTO shipments', () => {
    it.each([
      ['Requested delivery date', ': 14 Apr 2020'],
      ['Status', ': SUBMITTED'],
      ['Requested pickup date', ': 23 Mar 2020'],
      ['Destination type', ': Home of selection (HOS)'],
      ['Departure date', ': 15 Apr 2020'],
      ['Approved date', ': 21 Apr 2020'],
      ['Billable weight', ': 10,000 lbs'],
      ['Billable weight remarks', ': Heavy items'],
      ['Counselor remarks', ': Remarks'],
      ['Customer remarks', ': Words'],
      ['Diversion', ': false'],
      ['First available delivery date', ': 12 Apr 2020'],
      ['Shipment weight', ': 14,000 lbs'],
      ['Prime estimated weight', ': 12,000 lbs'],
      ['SAC type', ': HHG'],
      ['Scheduled pickup date', ': 22 Mar 2020'],
      ['Service order number', ': 767567576'],
      ['TAC type', ': NTS'],
      ['Uses external vendor', ': true'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
    it('displays the correct label for shipment', () => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText('PPM shipment #B4B4B')).toBeInTheDocument();
    });
  });
});
