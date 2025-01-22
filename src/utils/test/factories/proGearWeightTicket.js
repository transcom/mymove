import moment from 'moment';
import { v4 } from 'uuid';

import createUpload from 'utils/test/factories/upload';
import { createDocumentWithoutUploads } from 'utils/test/factories/document';
import { PPM_DOCUMENT_STATUS } from 'shared/constants';

const createBaseProGearWeightTicket = ({ serviceMemberId, creationDate = new Date() } = {}, fieldOverrides = {}) => {
  const createdAt = creationDate.toISOString();

  const smId = serviceMemberId || v4();
  const document = createDocumentWithoutUploads({ serviceMemberId: smId });

  return {
    id: v4(),
    ppmShipmentId: v4(),
    belongsToSelf: null,
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

const createCompleteProGearWeightTicket = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    belongsToSelf: true,
    description: 'Work equipment',
    hasWeightTickets: true,
    weight: 1500,
    ...fieldOverrides,
  };

  const weightTicket = createBaseProGearWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);

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

const createCompleteProGearWeightTicketWithConstructedWeight = (
  { serviceMemberId, creationDate } = {},
  fieldOverrides = {},
) => {
  const fullFieldOverrides = {
    belongsToSelf: true,
    description: 'Work equipment',
    hasWeightTickets: false,
    weight: 1400,
    ...fieldOverrides,
  };

  const weightTicket = createBaseProGearWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);

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

const createRejectedProGearWeightTicket = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    belongsToSelf: true,
    description: 'Laptop',
    hasWeightTickets: true,
    weight: 150,
    ...fieldOverrides,
  };
  const weightTicket = createBaseProGearWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);
  weightTicket.status = PPM_DOCUMENT_STATUS.REJECTED;
  return weightTicket;
};

export {
  createBaseProGearWeightTicket,
  createCompleteProGearWeightTicket,
  createCompleteProGearWeightTicketWithConstructedWeight,
  createRejectedProGearWeightTicket,
};
