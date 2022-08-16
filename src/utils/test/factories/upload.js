import { v4 } from 'uuid';

const createUpload = ({ fileName, createdAtDate }) => {
  const uploadId = v4();
  const uploadCreateDate = createdAtDate || new Date();

  return {
    id: uploadId,
    filename: fileName,
    status: 'PROCESSING',
    url: `/uploads/${uploadId}?contentType=application%2Fpdf`,
    content_type: 'application/pdf',
    bytes: 10596,
    created_at: uploadCreateDate.toISOString(),
    updated_at: uploadCreateDate.toISOString(),
  };
};

export default createUpload;
