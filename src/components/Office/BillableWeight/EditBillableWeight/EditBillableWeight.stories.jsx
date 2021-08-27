import React from 'react';

import EditBillableWeight from './EditBillableWeight';

export default {
  title: 'Office Components/EditBillableWeight',
  component: EditBillableWeight,
};

export const Basic = () => (
  <div style={{ width: 336, margin: '0 auto' }}>
    <EditBillableWeight weightAllowance={8000} estimatedWeight={13750} />
  </div>
);
