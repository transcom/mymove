import moment from 'moment';
import { v4 } from 'uuid';

import createUpload from 'utils/test/factories/upload';
import { createDocumentWithoutUploads } from 'utils/test/factories/document';

const createBaseWeightTicket = ({ serviceMemberId, creationDate = new Date() } = {}, fieldOverrides = {}) => {
  const weightTicketCreatedAtDate = creationDate.toISOString();

  const smId = serviceMemberId || v4();
  const emptyDocument = createDocumentWithoutUploads({ serviceMemberId: smId });
  const fullDocument = createDocumentWithoutUploads({ serviceMemberId: smId });
  const trailerDocument = createDocumentWithoutUploads({ serviceMemberId: smId });

  return {
    id: v4(),
    ppmShipmentId: v4(),
    vehicleDescription: null,
    emptyWeight: null,
    missingEmptyWeightTicket: null,
    emptyDocumentId: emptyDocument.id,
    emptyDocument,
    fullWeight: null,
    missingFullWeightTicket: null,
    fullDocumentId: fullDocument.id,
    fullDocument,
    ownsTrailer: null,
    trailerMeetsCriteria: null,
    proofOfTrailerOwnershipDocumentId: trailerDocument.id,
    proofOfTrailerOwnershipDocument: trailerDocument,
    createdAt: weightTicketCreatedAtDate,
    updatedAt: weightTicketCreatedAtDate,
    eTag: window.btoa(weightTicketCreatedAtDate),
    ...fieldOverrides,
  };
};

const createCompleteWeightTicket = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    vehicleDescription: '2022 Honda CR-V Hybrid',
    emptyWeight: 14500,
    fullWeight: 18500,
    ...fieldOverrides,
  };

  const weightTicket = createBaseWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);

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

const createCompleteWeightTicketWithTrailer = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const fullFieldOverrides = {
    ownsTrailer: true,
    trailerMeetsCriteria: true,
    ...fieldOverrides,
  };

  const weightTicket = createCompleteWeightTicket({ serviceMemberId, creationDate }, fullFieldOverrides);

  if (weightTicket.proofOfTrailerOwnershipDocument.uploads.length === 0) {
    weightTicket.proofOfTrailerOwnershipDocument.uploads.push(createUpload({ fileName: 'trailerDocument.pdf' }));
  }

  return weightTicket;
};

export { createBaseWeightTicket, createCompleteWeightTicket, createCompleteWeightTicketWithTrailer };
