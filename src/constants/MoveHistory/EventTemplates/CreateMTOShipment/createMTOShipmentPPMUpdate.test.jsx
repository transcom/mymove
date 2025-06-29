import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/CreateMTOShipment/createMTOShipmentPPMUpdate';

describe('When a PPM is created by the Prime and the move history is viewed', () => {
  const historyRecord = {
    action: 'UPDATE',
    eventName: 'createMTOShipment',
    tableName: 'ppm_shipments',
    changedValues: {
      destination_postal_code: '62269',
      estimated_weight: 6000,
      expected_departure_date: '2024-01-15',
      pickup_postal_code: '95630',
      has_pro_gear: true,
      pro_gear_weight: 2000,
      secondary_destination_postal_code: '95670',
      secondary_pickup_postal_code: '63108',
      spouse_pro_gear_weight: 500,
      status: 'SUBMITTED',
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'a1b2c',
      },
    ],
  };
  it('correctly matches the Create basic service item event', () => {
    const result = getTemplate(historyRecord);
    expect(result).toMatchObject(e);
    expect(result.getEventNameDisplay(result)).toEqual('Updated shipment');
  });
  describe('it correctly renders the details component for Create MTO shipments', () => {
    it('displays the correct detail values ', async () => {
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText('Status')).toBeInTheDocument();
      expect(screen.getByText(': SUBMITTED')).toBeInTheDocument();
      expect(screen.getByText('Pro-gear')).toBeInTheDocument();
      expect(screen.getByText(': Yes')).toBeInTheDocument();
      expect(screen.getByText('Pro-gear weight')).toBeInTheDocument();
      expect(screen.getByText(': 2,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Spouse pro-gear weight')).toBeInTheDocument();
      expect(screen.getByText(': 500 lbs')).toBeInTheDocument();
      expect(screen.getByText('Estimated weight')).toBeInTheDocument();
      expect(screen.getByText(': 6,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Expected departure date')).toBeInTheDocument();
      expect(screen.getByText(': 15 Jan 2024')).toBeInTheDocument();
      expect(screen.getByText('Pickup postal code')).toBeInTheDocument();
      expect(screen.getByText(': 95630')).toBeInTheDocument();
      expect(screen.getByText('Destination postal code')).toBeInTheDocument();
      expect(screen.getByText(': 62269')).toBeInTheDocument();
      expect(screen.getByText('Secondary pickup postal code')).toBeInTheDocument();
      expect(screen.getByText(': 63108')).toBeInTheDocument();
      expect(screen.getByText('Secondary destination postal code')).toBeInTheDocument();
      expect(screen.getByText(': 95670')).toBeInTheDocument();
    });
    it('displays the correct label for shipment', () => {
      const result = getTemplate(historyRecord);
      render(result.getDetails(historyRecord));
      expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
    });
  });
});
