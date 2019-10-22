import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const getUploadTagsLabel = 'UploadTags.getUploadTags';

export function getUploadTags(uploadID, label = getUploadTagsLabel) {
  const swaggerTag = 'uploads.getUploadTags';
  return swaggerRequest(getClient, swaggerTag, { uploadID }, { label });
}
