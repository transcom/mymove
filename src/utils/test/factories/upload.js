import { v4 } from 'uuid';

import { UPLOAD_SCAN_STATUS } from 'shared/constants';

const createUpload = ({ fileName, createdAtDate = new Date() } = {}, fieldOverrides = {}) => {
  const uploadId = v4();
  const uploadCreateDate = createdAtDate.toISOString();

  const contentType = fieldOverrides?.contentType ? fieldOverrides?.contentType : 'application/pdf';

  const url = `/uploads/${uploadId}?contentType=${encodeURIComponent(contentType)}`;

  return {
    id: uploadId,
    filename: fileName,
    status: UPLOAD_SCAN_STATUS.PROCESSING,
    contentType,
    url,
    bytes: 10596,
    createdAt: uploadCreateDate,
    updatedAt: uploadCreateDate,
    ...fieldOverrides,
  };
};

export default createUpload;
