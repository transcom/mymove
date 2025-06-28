import moment from 'moment';
import { v4 } from 'uuid';

import createUpload from 'utils/test/factories/upload';
import { createDocumentWithoutUploads } from 'utils/test/factories/document';
import PPMDocumentsStatus from 'constants/ppms';

const createBaseGunSafeWeightTicket = ({ serviceMemberId, creationDate = new Date() } = {}, fieldOverrides = {}) => {
  const createdAt = creationDate.toISOString();

  const smId = serviceMemberId || v4();
  const document = createDocumentWithoutUploads({ serviceMemberId: smId });

  return {
    id: v4(),
    ppmShipmentId: v4(),
    description: null,
    hasWeightTickets: null,
    documentId: document.id,
    document,
    weight: null,
    status: null,
    reason: null,
    createdAt,
    updatedAt: createdAt,
    eTag: window.btoa(createdAt),
    ...fieldOverrides,
  };
};

const createCompleteGunSafeWeightTicket = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    description: 'Gun safe',
    hasWeightTickets: true,
    weight: 400,
    ...fieldOverrides,
  };

  const weightTicket = createBaseGunSafeWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);

  if (weightTicket.createdAt === weightTicket.updatedAt) {
    const updatedAt = moment(weightTicket.createdAt).add(1, 'hour').toISOString();

    weightTicket.updatedAt = updatedAt;
    weightTicket.eTag = window.btoa(updatedAt);
  }

  if (weightTicket.document.uploads.length === 0) {
    weightTicket.document.uploads.push(createUpload({ fileName: 'emptyDocument.pdf' }));
  }

  return weightTicket;
};

const createCompleteGunSafeWeightTicketWithConstructedWeight = (
  { serviceMemberId, creationDate } = {},
  fieldOverrides = {},
) => {
  const fullFieldOverrides = {
    description: 'Gun Safe',
    hasWeightTickets: false,
    weight: 1400,
    ...fieldOverrides,
  };

  const weightTicket = createBaseGunSafeWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);

  if (weightTicket.constructedWeightDocument.uploads.length === 0) {
    weightTicket.constructedWeightDocument.uploads.push(createUpload({ fileName: 'constructedWeight.pdf' }));
  }

  if (weightTicket.createdAt === weightTicket.updatedAt) {
    const updatedAt = moment(weightTicket.createdAt).add(1, 'hour').toISOString();

    weightTicket.updatedAt = updatedAt;
    weightTicket.eTag = window.btoa(updatedAt);
  }

  return weightTicket;
};

const createRejectedGunSafeWeightTicket = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    description: 'Gun Safe',
    hasWeightTickets: true,
    weight: 150,
    ...fieldOverrides,
  };
  const weightTicket = createBaseGunSafeWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);
  weightTicket.status = PPMDocumentsStatus.REJECTED;
  return weightTicket;
};

export {
  createBaseGunSafeWeightTicket,
  createCompleteGunSafeWeightTicket,
  createCompleteGunSafeWeightTicketWithConstructedWeight,
  createRejectedGunSafeWeightTicket,
};
