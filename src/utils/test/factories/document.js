import { v4 } from 'uuid';

import createUpload from 'utils/test/factories/upload';

const createDocumentWithoutUploads = ({ serviceMemberId, uploads } = {}) => {
  return {
    id: v4(),
    service_member_id: serviceMemberId || v4(),
    uploads: uploads || [],
  };
};

const createDocumentWithUploads = ({ serviceMemberId, uploadsArgs = [] } = {}) => {
  const uploads = [];

  uploadsArgs.forEach((uploadArgs) => {
    uploads.push(createUpload(uploadArgs));
  });

  if (uploads.length === 0) {
    uploads.push(createUpload({ fileName: 'testFile.pdf' }));
  }

  return createDocumentWithoutUploads({ serviceMemberId, uploads });
};

export { createDocumentWithoutUploads, createDocumentWithUploads };
