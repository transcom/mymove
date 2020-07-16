import React from 'react';
import { array, text } from '@storybook/addon-knobs';

import HHGDetailsForm from './HHGDetailsForm';

export default {
  title: 'Customer Components | HHGDetailsForm',
};

const pageKey = 'pageKey';
const pageList = ['page1', 'anotherPage/:foo/:bar'];
export const Basic = () => <HHGDetailsForm pageList={array('pageList', pageList)} pageKey={text('pageKey', pageKey)} />;

// every named export is a test case
