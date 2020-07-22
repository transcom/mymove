import React from 'react';

import HHGDetailsForm from './HHGDetailsForm';

export default {
  title: 'Customer Components | HHGDetailsForm',
};

const pageKey = 'pageKey';
const pageList = ['page1', 'anotherPage/:foo/:bar'];
export const Basic = () => <HHGDetailsForm pageList={pageList} pageKey={pageKey} />;
