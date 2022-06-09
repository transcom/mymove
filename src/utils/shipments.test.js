import { isPPMShipmentComplete } from './shipments';

import { SHIPMENT_OPTIONS } from 'shared/constants';

describe('shipments utils', () => {
  describe('isPPMShipmentComplete', () => {
    it('returns true when the hasRequestedAdvance field is set to true', () => {
      const completePPMShipment = {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: '2',
          expectedDepartureDate: '2022-04-01',
          pickupPostalCode: '90210',
          destinationPostalCode: '90211',
          sitExpected: false,
          estimatedWeight: 7999,
          hasProGear: false,
          estimatedIncentive: 1234500,
          hasRequestedAdvance: true,
          advanceAmountRequested: 487500,
        },
      };
      expect(isPPMShipmentComplete(completePPMShipment)).toBe(true);
    });

    it('returns true when the hasRequestedAdvance field is set to false', () => {
      const completePPMShipment = {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: '2',
          expectedDepartureDate: '2022-04-01',
          pickupPostalCode: '90210',
          destinationPostalCode: '90211',
          sitExpected: false,
          estimatedWeight: 7999,
          hasProGear: false,
          estimatedIncentive: 1234500,
          hasRequestedAdvance: false,
        },
      };
      expect(isPPMShipmentComplete(completePPMShipment)).toBe(true);
    });

    it('returns false when the hasRequestedAdvance field is undefined', () => {
      const incompletePPMShipment = {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: '2',
          expectedDepartureDate: '2022-04-01',
          pickupPostalCode: '90210',
          destinationPostalCode: '90211',
          sitExpected: false,
          estimatedWeight: 7999,
          hasProGear: false,
          estimatedIncentive: 1234500,
        },
      };
      expect(isPPMShipmentComplete(incompletePPMShipment)).toBe(false);
    });

    it('returns false when the hasRequestedAdvance field is null', () => {
      const incompletePPMShipment = {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: '2',
          expectedDepartureDate: '2022-04-01',
          pickupPostalCode: '90210',
          destinationPostalCode: '90211',
          sitExpected: false,
          estimatedWeight: 7999,
          hasProGear: false,
          estimatedIncentive: 1234500,
          hasRequestedAdvance: null,
        },
      };
      expect(isPPMShipmentComplete(incompletePPMShipment)).toBe(false);
    });
  });
});
