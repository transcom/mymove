import { v4 } from 'uuid';

const createDocumentWithoutUploads = ({ serviceMemberId, uploads }) => {
  return {
    id: v4(),
    service_member_id: serviceMemberId || v4(),
    uploads: uploads || [],
  };
};

export default createDocumentWithoutUploads;
