import moment from 'moment';
import { v4 } from 'uuid';

import createUpload from 'utils/test/factories/upload';
import { createDocumentWithoutUploads } from 'utils/test/factories/document';

const createBaseProGearWeightTicket = ({ serviceMemberId, creationDate = new Date() } = {}, fieldOverrides = {}) => {
  const createdAt = creationDate.toISOString();

  const smId = serviceMemberId || v4();
  const emptyDocument = createDocumentWithoutUploads({ serviceMemberId: smId });
  const fullDocument = createDocumentWithoutUploads({ serviceMemberId: smId });
  const constructedWeightDocument = createDocumentWithoutUploads({ serviceMemberId: smId });

  return {
    id: v4(),
    ppmShipmentId: v4(),
    belongsToSelf: null,
    description: null,
    hasWeightTickets: null,
    emptyWeight: null,
    emptyDocumentId: emptyDocument.id,
    emptyDocument,
    fullWeight: null,
    fullDocumentId: fullDocument.id,
    fullDocument,
    constructedWeight: null,
    constructedWeightDocumentId: constructedWeightDocument.id,
    constructedWeightDocument,
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
    emptyWeight: 14500,
    fullWeight: 16000,
    ...fieldOverrides,
  };

  const weightTicket = createBaseProGearWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);

  if (weightTicket.createdAt === weightTicket.updatedAt) {
    const updatedAt = moment(weightTicket.createdAt).add(1, 'hour').toISOString();

    weightTicket.updatedAt = updatedAt;
    weightTicket.eTag = window.btoa(updatedAt);
  }

  if (weightTicket.emptyDocument.uploads.length === 0) {
    weightTicket.emptyDocument.uploads.push(createUpload({ fileName: 'emptyDocument.pdf' }));
  }

  if (weightTicket.fullDocument.uploads.length === 0) {
    weightTicket.fullDocument.uploads.push(createUpload({ fileName: 'fullDocument.pdf' }));
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
    constructedWeight: 1400,
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

export {
  createBaseProGearWeightTicket,
  createCompleteProGearWeightTicket,
  createCompleteProGearWeightTicketWithConstructedWeight,
};
