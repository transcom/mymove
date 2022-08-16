import moment from 'moment';
import { v4 } from 'uuid';

import createUpload from 'utils/test/factories/upload';

const createBaseWeightTicket = ({ serviceMemberId, creationDate } = {}, fieldOverrides = {}) => {
  const weightTicketCreatedAtDate = creationDate || new Date();

  const smId = serviceMemberId || v4();
  const emptyDocumentId = v4();
  const fullDocumentId = v4();
  const trailerDocumentId = v4();

  return {
    id: v4(),
    ppmShipmentId: v4(),
    vehicleDescription: null,
    emptyWeight: null,
    missingEmptyWeightTicket: null,
    emptyDocumentId,
    emptyDocument: {
      id: emptyDocumentId,
      service_member_id: smId,
      uploads: [],
    },
    fullWeight: null,
    missingFullWeightTicket: null,
    fullDocumentId,
    fullDocument: {
      id: fullDocumentId,
      service_member_id: smId,
      uploads: [],
    },
    ownsTrailer: null,
    trailerMeetsCriteria: null,
    proofOfTrailerOwnershipDocumentId: trailerDocumentId,
    proofOfTrailerOwnershipDocument: {
      id: trailerDocumentId,
      service_member_id: smId,
      uploads: [],
    },
    createdAt: weightTicketCreatedAtDate.toISOString(),
    updatedAt: weightTicketCreatedAtDate.toISOString(),
    eTag: window.btoa(weightTicketCreatedAtDate.toISOString()),
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
