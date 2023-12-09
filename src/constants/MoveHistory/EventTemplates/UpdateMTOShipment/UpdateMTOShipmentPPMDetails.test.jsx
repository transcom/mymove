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
        shipment_id_abbr: '12992',
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
    expect(screen.getByText('PPM shipment #12992')).toBeInTheDocument();
  });
  describe('it correctly renders the shipment details', () => {
    it.each([
      ['Estimated weight', ': 2,233 lbs'],
      ['Expected departure date', ': 08 Dec 2023'],
      ['Sit expected', ': Yes'],
      ['Estimated storage start', ': 22 Mar 2020'],
      ['Estimated storage end', ': 13 Apr 2020'],
      ['Estimated storage weight', ': 6,877 lbs'],
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
});
