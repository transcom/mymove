import { v4 } from 'uuid';

import { UPLOAD_SCAN_STATUS } from 'shared/constants';

const createUpload = ({ fileName, createdAtDate = new Date() } = {}) => {
  const uploadId = v4();
  const uploadCreateDate = createdAtDate.toISOString();

  return {
    id: uploadId,
    filename: fileName,
    status: UPLOAD_SCAN_STATUS.PROCESSING,
    url: `/uploads/${uploadId}?contentType=application%2Fpdf`,
    contentType: 'application/pdf',
    bytes: 10596,
    createdAt: uploadCreateDate,
    updatedAt: uploadCreateDate,
  };
};

export default createUpload;
