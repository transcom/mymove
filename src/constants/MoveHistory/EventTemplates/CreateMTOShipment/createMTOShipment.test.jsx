import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOShipment/createMTOShipment';

describe('when given a Create mto shipment history record', () => {
  const historyRecord = {
    action: 'INSERT',
    eventName: 'createMTOShipment',
    tableName: 'mto_shipments',
    changedValues: {
      customer_remarks: 'Redacted',
      requested_delivery_date: '2022-12-12',
      requested_pickup_date: '2022-12-10',
      status: 'SUBMITTED',
    },
    context: [
      {
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
      },
    ],
  };
  it('correctly matches the Create basic service item event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(result)).toEqual('Created shipment');
  });
  describe('it correctly renders the details component for Create MTO shipments', () => {
    it.each([
      ['Requested delivery date', ': 12 Dec 2022'],
      ['Status', ': SUBMITTED'],
      ['Customer remarks', ': Redacted'],
      ['Requested pickup date', ': 10 Dec 2022'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
    it('displays the correct label for shipment', () => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText('HHG shipment #A1B2C')).toBeInTheDocument();
    });
  });
});
