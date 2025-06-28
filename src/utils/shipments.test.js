import moment from 'moment';
import { v4 } from 'uuid';

import {
  isPPMAboutInfoComplete,
  isPPMShipmentComplete,
  isWeightTicketComplete,
  hasCompletedAllWeightTickets,
  isPPMOnly,
  isGunSafeComplete,
  blankContact,
} from './shipments';

import { ppmShipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { createDocumentWithoutUploads } from 'utils/test/factories/document';
import {
  createBaseWeightTicket,
  createCompleteWeightTicket,
  createCompleteWeightTicketWithTrailer,
} from 'utils/test/factories/weightTicket';

describe('shipments utils', () => {
  describe('isPPMShipmentComplete', () => {
    it('returns true when the hasRequestedAdvance field is set to true', () => {
      const completePPMShipment = {
        id: '1',
        shipmentType: SHIPMENT_OPTIONS.PPM,
        ppmShipment: {
          id: '2',
          expectedDepartureDate: '2022-04-01',
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

  describe('isPPMAboutInfoComplete', () => {
    const createdDate = new Date();
    const firstUpdatedDate = moment(createdDate).add(1, 'day');
    const approvedDate = moment(firstUpdatedDate).add(1, 'day');
    const secondUpdatedDate = moment(approvedDate).add(1, 'day');

    const approvedPPMShipment = {
      id: v4(),
      shipmentId: v4(),
      status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
      expectedDepartureDate: '2022-08-08',
      sitExpected: false,
      estimatedWeight: 4000,
      hasProGear: false,
      estimatedIncentive: 10000000,
      hasRequestedAdvance: true,
      advanceAmountRequested: 30000,
      actualMoveDate: null,
      hasReceivedAdvance: null,
      advanceAmountReceived: null,
      approvedAt: approvedDate.toISOString(),
      createdAt: createdDate.toISOString(),
      updatedAt: firstUpdatedDate.toISOString(),
      eTag: window.btoa(firstUpdatedDate.toISOString()),
    };

    const ppmShipmentWithAboutInfo = {
      ...approvedPPMShipment,
      actualMoveDate: approvedPPMShipment.expectedDepartureDate,
      hasReceivedAdvance: approvedPPMShipment.hasRequestedAdvance,
      advanceAmountReceived: approvedPPMShipment.advanceAmountRequested,
      updatedAt: secondUpdatedDate.toISOString(),
      eTag: window.btoa(secondUpdatedDate.toISOString()),
    };

    const ppmShipmentWithAboutInfoWithoutAdvance = {
      ...ppmShipmentWithAboutInfo,
      hasRequestedAdvance: false,
      advanceAmountRequested: null,
      hasReceivedAdvance: false,
      advanceAmountReceived: null,
    };

    it.each([
      [false, 'all about your ppm info is null', approvedPPMShipment],
      [false, 'actual move date is null', { ...ppmShipmentWithAboutInfo, actualMoveDate: null }],
      [false, 'has received advance is null', { ...ppmShipmentWithAboutInfo, hasReceivedAdvance: null }],
      [false, 'advance amount received is too low', { ...ppmShipmentWithAboutInfo, advanceAmountReceived: 0 }],
      [true, 'all about your ppm info is filled in (no advance)', ppmShipmentWithAboutInfoWithoutAdvance],
      [true, 'all about your ppm info is filled in (with advance)', ppmShipmentWithAboutInfo],
    ])('returns %s when %s', (expectedResult, scenarioDescription, ppmShipment) => {
      expect(isPPMAboutInfoComplete(ppmShipment)).toBe(expectedResult);
    });
  });

  describe('isWeightTicketComplete', () => {
    const serviceMemberId = v4();

    const completeWeightTicket = createCompleteWeightTicket({ serviceMemberId });
    const completeWeightTicketWithTrailer = createCompleteWeightTicketWithTrailer({ serviceMemberId });

    const emptyDocumentWithoutUploads = createDocumentWithoutUploads({ serviceMemberId });
    const fullDocumentWithoutUploads = createDocumentWithoutUploads({ serviceMemberId });
    const trailerDocumentWithoutUploads = createDocumentWithoutUploads({ serviceMemberId });

    it.each([
      [false, 'vehicle description is missing', { ...completeWeightTicket, vehicleDescription: null }],
      [false, 'empty weight is missing', { ...completeWeightTicket, emptyWeight: null }],
      [false, 'empty document has no uploads', { ...completeWeightTicket, emptyDocument: emptyDocumentWithoutUploads }],
      [false, 'full weight is missing', { ...completeWeightTicket, fullWeight: null }],
      [false, 'full document has no uploads', { ...completeWeightTicket, fullDocument: fullDocumentWithoutUploads }],
      [
        false,
        'owns trailer but missing trailer uploads',
        {
          ...completeWeightTicketWithTrailer,
          proofOfTrailerOwnershipDocument: trailerDocumentWithoutUploads,
        },
      ],
      [true, 'all required data is present (no trailer)', completeWeightTicket],
      [true, 'all required data is present (empty weight === 0)', { ...completeWeightTicket, emptyWeight: 0 }],
      [true, 'all required data is present (with trailer)', completeWeightTicketWithTrailer],
    ])('returns %s if %s', (expectedValue, scenarioDescription, weightTicket) => {
      expect(isWeightTicketComplete(weightTicket)).toEqual(expectedValue);
    });
  });

  describe('hasCompletedAllWeightTickets', () => {
    it('returns false when there are no weight tickets', () => {
      expect(hasCompletedAllWeightTickets()).toBe(false);
      expect(hasCompletedAllWeightTickets([])).toBe(false);
    });
    it('returns false when there is at least one incomplete weight ticket', () => {
      expect(hasCompletedAllWeightTickets([createBaseWeightTicket()])).toBe(false);
      expect(hasCompletedAllWeightTickets([createBaseWeightTicket(), createCompleteWeightTicket()])).toBe(false);
    });
    it('returns true when all weight tickets are complete', () => {
      expect(hasCompletedAllWeightTickets([createBaseWeightTicket()])).toBe(false);
      expect(hasCompletedAllWeightTickets([createCompleteWeightTicket(), createCompleteWeightTicket()])).toBe(true);
    });
  });

  describe('isPPMOnly', () => {
    const ppmShipment = {
      id: v4(),
      shipmentType: SHIPMENT_OPTIONS.PPM,
    };
    const secondPPMShipment = {
      id: v4(),
      shipmentType: SHIPMENT_OPTIONS.PPM,
    };
    const hhgShipment = {
      id: v4(),
      shipmentType: SHIPMENT_OPTIONS.HHG,
    };
    const ntsShipment = {
      id: v4(),
      shipmentType: SHIPMENT_OPTIONS.NTS,
    };
    it('returns false when there are no shipments', () => {
      expect(isPPMOnly()).toBe(false);
      expect(isPPMOnly([])).toBe(false);
    });
    it('returns false when not all shipments are PPM', () => {
      expect(isPPMOnly([ppmShipment, hhgShipment, ntsShipment])).toBe(false);
      expect(isPPMOnly([hhgShipment, ntsShipment])).toBe(false);
    });
    it('returns true when all shipments are PPM', () => {
      expect(isPPMOnly([ppmShipment])).toBe(true);
      expect(isPPMOnly([ppmShipment, secondPPMShipment])).toBe(true);
    });
  });

  describe('isGunSafeComplete', () => {
    const completeGunSafeTicket = {
      description: 'Test Description',
      weight: 500,
      document: {
        uploads: {
          filename: 'test.pdf',
        },
      },
    };

    const gunSafeNoDescription = {
      weight: 500,
      document: {
        uploads: {
          filename: 'test.pdf',
        },
      },
    };

    const gunSafeNoWeight = {
      description: 'test description',
      document: {
        uploads: {
          filename: 'test.pdf',
        },
      },
    };

    const gunSafeNoUpload = {
      description: 'test description',
      weight: 500,
      document: {},
    };

    it('returns true when a ticket has a weight, description, and document upload', () => {
      expect(isGunSafeComplete(completeGunSafeTicket));
    });

    it('returns false when a ticket has a description and document upload, but no weight', () => {
      expect(!isGunSafeComplete(gunSafeNoWeight));
    });

    it('returns false when a ticket has a weight and document upload, but no description', () => {
      expect(!isGunSafeComplete(gunSafeNoDescription));
    });

    it('returns false when a ticket has a description and weight, but no document uploads', () => {
      expect(!isGunSafeComplete(gunSafeNoUpload));
    });
  });

  describe('blankContact', () => {
    it('returns blank contact', () => {
      const contact = blankContact;
      expect(contact.contact.firstName).toBe('');
      expect(contact.contact.lastName).toBe('');
      expect(contact.contact.email).toBe('');
      expect(contact.contact.telephone).toBe('');
    });
  });
});
