import React from 'react';

import UploadsTable from './UploadsTable';

export const uploadsTable = () => (
  <UploadsTable
    uploads={[
      {
        bytes: 9043,
        contentType: 'image/png',
        createdAt: '2021-06-21T19:51:49.441Z',
        filename: 'orders1.png',
        id: '2e405ced-8b35-4298-a59a-426f5cef1d58',
        url: '/storage/user/72068b0c-ed23-4a77-beff-1d2bb3e713b0/uploads/2e405ced-8b35-4298-a59a-426f5cef1d58?contentType=image%2Fpng',
      },
      {
        bytes: 9043,
        contentType: 'application/pdf',
        createdAt: '2021-06-21T20:33:22.724Z',
        filename: 'orders2.pdf',
        id: '75bc82cd-c584-424d-83df-2045ff13611f',
        url: '/storage/user/72068b0c-ed23-4a77-beff-1d2bb3e713b0/uploads/75bc82cd-c584-424d-83df-2045ff13611f?contentType=application%2Fpdf',
      },
      {
        bytes: 9043,
        contentType: 'image/png',
        createdAt: '2021-06-21T20:40:52.246Z',
        filename: 'orders3.png',
        id: 'b469b0f9-9080-49cc-87cc-5a76c0e947f1',
        url: '/storage/user/72068b0c-ed23-4a77-beff-1d2bb3e713b0/uploads/b469b0f9-9080-49cc-87cc-5a76c0e947f1?contentType=image%2Fpng',
      },
    ]}
    onDelete={() => {}}
  />
);

export default { title: 'Customer Components/UploadsTable' };
