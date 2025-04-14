import { render, screen } from '@testing-library/react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/UpdateMTOShipmentPPMDetails';

describe('when given a PPM shipment update', () => {
  const historyRecord = {
    action: a.UPDATE,
    eventName: o.updateMTOShipment,
    tableName: t.ppm_shipments,
    changedValues: {
      destination_postal_code: '59801',
      estimated_weight: 2233,
      expected_departure_date: '2023-12-08',
      pickup_postal_code: '20906',
      secondary_destination_postal_code: '59802',
      secondary_pickup_postal_code: '20832',
      sit_estimated_departure_date: '2020-04-13',
      sit_estimated_entry_date: '2020-03-22',
      sit_estimated_weight: '6877',
      sit_expected: true,
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: '12992',
      },
    ],
  };

  const weightRecord = {
    action: a.UPDATE,
    eventName: o.updateMTOShipment,
    tableName: t.ppm_shipments,
    changedValues: {
      estimated_weight: '1234',
      has_pro_gear: true,
      pro_gear_weight: '321',
      spouse_pro_gear_weight: '34',
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: '12992',
      },
    ],
  };

  const advanceRecord = {
    action: a.UPDATE,
    eventName: o.updateMTOShipment,
    tableName: t.ppm_shipments,
    changedValues: {
      advance_amount_requested: 598600,
      advance_status: 'APPROVED',
      has_requested_advance: true,
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: '12992',
      },
    ],
  };

  const incentiveRecord = {
    action: a.UPDATE,
    eventName: o.updateMTOShipment,
    tableName: t.ppm_shipments,
    changedValues: {
      final_incentive: 46855239,
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: '12992',
      },
    ],
  };

  const w2Record = {
    action: a.UPDATE,
    eventName: o.updateMTOShipment,
    tableName: t.ppm_shipments,
    changedValues: {
      actual_destination_postal_code: '20889',
      actual_move_date: '2024-02-01',
      actual_pickup_postal_code: '59801',
      advance_amount_received: 2278600,
      has_received_advance: true,
      w2_address_id: '6c1df26c-76bc-42c4-9c7f-c626a78739d3',
    },
    context: [
      {
        shipment_id_abbr: '125d1',
        shipment_locator: 'RQ38D4-01',
        shipment_type: 'PPM',
        w2_address:
          '{"id":"6c1df26c-76bc-42c4-9c7f-c626a78739d3","street_address_1":"123 Dad St","street_address_2":null,"city":"Missoula","state":"MT","postal_code":"59801","created_at":"2024-02-15T08:39:41.692044","updated_at":"2024-02-16T03:36:51.388835","street_address_3":null,"country":null}',
      },
    ],
  };

  it('correctly matches the the event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
  });
  it('displays the correct label for shipment', () => {
    const result = getTemplate(historyRecord);
    render(result.getDetails(historyRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
  });
  describe('it correctly renders the shipment details', () => {
    it.each([
      ['Estimated weight', ': 2,233 lbs'],
      ['Expected departure date', ': 08 Dec 2023'],
      ['SIT expected', ': Yes'],
      ['SIT estimated storage start', ': 22 Mar 2020'],
      ['SIT estimated storage end', ': 13 Apr 2020'],
      ['SIT estimated storage weight', ': 6,877 lbs'],
      ['Pickup postal code', ': 20906'],
      ['Destination postal code', ': 59801'],
      ['Secondary pickup postal code', ': 20832'],
      ['Secondary destination postal code', ': 59802'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
  describe('it correctly renders the weight details', () => {
    it.each([
      ['Pro-gear', ': Yes'],
      ['Pro-gear weight', ': 321 lbs'],
      ['Spouse pro-gear weight', ': 34 lbs'],
      ['Estimated weight', ': 1,234 lbs'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(weightRecord);
      render(result.getDetails(weightRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
  describe('it correctly renders the advance details', () => {
    it.each([
      ['Advance amount requested', ': $5,986.00'],
      ['Advance requested', ': Yes'],
      ['Advance status', ': APPROVED'],
    ])('displays the correct details value for %s', async (label, value) => {
      const result = getTemplate(advanceRecord);
      render(result.getDetails(advanceRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
  it('it correctly renders the final incentive amount in dollars', () => {
    const result = getTemplate(incentiveRecord);
    render(result.getDetails(incentiveRecord));
    expect(screen.getByText('Final incentive')).toBeInTheDocument();
    expect(screen.getByText(': $468,552.39')).toBeInTheDocument();
  });
  it('it correctly renders details of a customer completing a PPM', () => {
    const result = getTemplate(w2Record);
    render(result.getDetails(w2Record));
    expect(result.getEventNameDisplay(w2Record)).toEqual('Customer Began PPM Document Process');
    expect(screen.getByText('W2 Address')).toBeInTheDocument();
    expect(screen.getByText(': 123 Dad St, Missoula, MT 59801')).toBeInTheDocument();
    expect(screen.getByText('Departure date')).toBeInTheDocument();
    expect(screen.getByText(': 01 Feb 2024')).toBeInTheDocument();
    expect(screen.getByText('Starting ZIP')).toBeInTheDocument();
    expect(screen.getByText(': 59801')).toBeInTheDocument();
    expect(screen.getByText('Ending ZIP')).toBeInTheDocument();
    expect(screen.getByText(': 20889')).toBeInTheDocument();
  });
});
