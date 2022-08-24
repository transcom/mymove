import moment from 'moment';
import { v4 } from 'uuid';

import { ppmShipmentStatuses, shipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { createCompleteMovingExpense } from 'utils/test/factories/movingExpense';
import { createCompleteProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';

const mergePPMFieldOverrides = (defaultOverrides = {}, inputOverrides = {}) => {
  return {
    ...defaultOverrides,
    ...inputOverrides,
    ppmShipment: {
      ...defaultOverrides?.ppmShipment,
      ...inputOverrides?.ppmShipment,
    },
  };
};

const createBasePPMShipment = (fieldOverrides = {}) => {
  const mtoPPMShipmentId = v4();
  const mtoShipmentCreatedDate = new Date().toISOString();
  const ppmShipmentCreatedDate = moment(mtoShipmentCreatedDate).add(5, 'seconds').toISOString();

  return {
    id: mtoPPMShipmentId,
    shipmentType: SHIPMENT_OPTIONS.PPM,
    status: shipmentStatuses.SUBMITTED,
    moveTaskOrderId: v4(),
    createdAt: mtoShipmentCreatedDate,
    updatedAt: mtoShipmentCreatedDate,
    eTag: window.btoa(mtoShipmentCreatedDate),
    ...fieldOverrides,
    ppmShipment: {
      id: v4(),
      shipmentId: mtoPPMShipmentId,
      status: ppmShipmentStatuses.DRAFT,
      expectedDepartureDate: '2022-09-15',
      pickupPostalCode: '90210',
      destinationPostalCode: '30813',
      sitExpected: false,
      estimatedWeight: null,
      hasProGear: null,
      estimatedIncentive: null,
      hasRequestedAdvance: null,
      advanceAmountRequested: null,
      actualMoveDate: null,
      actualPickupPostalCode: null,
      actualDestinationPostalCode: null,
      hasReceivedAdvance: null,
      advanceAmountReceived: null,
      finalEstimatedIncentive: null,
      weightTickets: [],
      movingExpenses: [],
      proGearWeightTickets: [],
      createdAt: ppmShipmentCreatedDate,
      updatedAt: ppmShipmentCreatedDate,
      eTag: window.btoa(ppmShipmentCreatedDate),
      ...fieldOverrides?.ppmShipment,
    },
  };
};

const createFilledOutUnSubmittedPPMShipment = (fieldOverrides = {}) => {
  const fullFieldOverrides = mergePPMFieldOverrides(
    {
      ppmShipment: {
        sitExpected: false,
        estimatedWeight: 4000,
        hasProGear: false,
        estimatedIncentive: 10000000,
        hasRequestedAdvance: true,
        advanceAmountRequested: 5000000,
      },
    },
    fieldOverrides,
  );

  const shipment = createBasePPMShipment(fullFieldOverrides);

  if (shipment.ppmShipment.createdAt === shipment.ppmShipment.updatedAt) {
    const updatedAt = moment(shipment.ppmShipment.createdAt).add(1, 'hour').toISOString();

    shipment.ppmShipment.updatedAt = updatedAt;
    shipment.ppmShipment.eTag = window.btoa(updatedAt);
  }

  return shipment;
};

const createSubmittedPPMShipment = (fieldOverrides = {}) => {
  const fullFieldOverrides = mergePPMFieldOverrides(
    {
      ppmShipment: {
        status: ppmShipmentStatuses.SUBMITTED,
      },
    },
    fieldOverrides,
  );

  const shipment = createFilledOutUnSubmittedPPMShipment(fullFieldOverrides);

  if (shipment.ppmShipment.createdAt === shipment.ppmShipment.updatedAt) {
    const updatedAt = moment(shipment.ppmShipment.createdAt).add(1, 'hour').toISOString();

    shipment.ppmShipment.updatedAt = updatedAt;
    shipment.ppmShipment.eTag = window.btoa(updatedAt);
  }

  return shipment;
};

const createApprovedPPMShipment = (fieldOverrides = {}) => {
  const fullFieldOverrides = mergePPMFieldOverrides(
    {
      status: shipmentStatuses.APPROVED,
      ppmShipment: {
        status: ppmShipmentStatuses.WAITING_ON_CUSTOMER,
      },
    },
    fieldOverrides,
  );

  const shipment = createSubmittedPPMShipment(fullFieldOverrides);

  const approvedAt = fieldOverrides?.approvedAt
    ? fieldOverrides?.approvedAt
    : moment(shipment.ppmShipment.updatedAt).add(2, 'days').toISOString();

  return {
    ...shipment,
    ppmShipment: {
      ...shipment.ppmShipment,
      approvedAt,
      updatedAt: approvedAt,
      eTag: window.btoa(approvedAt),
    },
    updatedAt: approvedAt,
    eTag: window.btoa(approvedAt),
  };
};

const createPPMShipmentWithActualShipmentInfo = (fieldOverrides = {}) => {
  const shipment = createApprovedPPMShipment(fieldOverrides);

  const updatedAt = fieldOverrides?.updatedAt
    ? fieldOverrides?.updatedAt
    : moment(shipment.ppmShipment.updatedAt).add(15, 'days').toISOString();

  return {
    ...shipment,
    ppmShipment: {
      ...shipment.ppmShipment,
      actualMoveDate: shipment.ppmShipment.expectedDepartureDate,
      actualPickupPostalCode: shipment.ppmShipment.pickupPostalCode,
      actualDestinationPostalCode: shipment.ppmShipment.destinationPostalCode,
      hasReceivedAdvance: shipment.ppmShipment.hasRequestedAdvance,
      advanceAmountReceived: shipment.ppmShipment.advanceAmountRequested,
      updatedAt,
      eTag: window.btoa(updatedAt),
      ...fieldOverrides?.ppmShipment,
    },
  };
};

const createPPMShipmentWithDocuments = (fieldOverrides = {}) => {
  const shipment = createPPMShipmentWithActualShipmentInfo(fieldOverrides);

  const ppmShipmentId = shipment.ppmShipment.id;
  const serviceMemberId = v4();

  if (shipment.ppmShipment.weightTickets.length === 0) {
    shipment.ppmShipment.weightTickets.push(createCompleteWeightTicket({ serviceMemberId }, { ppmShipmentId }));
  }

  if (shipment.ppmShipment.movingExpenses.length === 0) {
    shipment.ppmShipment.movingExpenses.push(createCompleteMovingExpense({ serviceMemberId }, { ppmShipmentId }));
  }

  if (shipment.ppmShipment.proGearWeightTickets.length === 0) {
    shipment.ppmShipment.proGearWeightTickets.push(
      createCompleteProGearWeightTicket({ serviceMemberId }, { ppmShipmentId }),
    );
  }

  return shipment;
};

const createPPMShipmentWithFinalIncentive = (fieldOverrides = {}) => {
  const shipment = createPPMShipmentWithDocuments(fieldOverrides);

  const updatedAt = fieldOverrides?.updatedAt
    ? fieldOverrides?.updatedAt
    : moment(shipment.ppmShipment.updatedAt).add(3, 'hours').toISOString();

  return {
    ...shipment,
    ppmShipment: {
      ...shipment.ppmShipment,
      finalEstimatedIncentive: shipment.ppmShipment.estimatedIncentive,
      updatedAt,
      eTag: window.btoa(updatedAt),
      ...fieldOverrides?.ppmShipment,
    },
  };
};

export {
  createApprovedPPMShipment,
  createBasePPMShipment,
  createFilledOutUnSubmittedPPMShipment,
  createPPMShipmentWithActualShipmentInfo,
  createPPMShipmentWithDocuments,
  createPPMShipmentWithFinalIncentive,
  createSubmittedPPMShipment,
};
