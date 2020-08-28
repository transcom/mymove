/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { ReviewSectionComponent as ReviewSection } from './ReviewSection';

const defaultProps = {
  fieldData: [{ label: 'Some heading' }, { value: 'Some value' }, { key: 'some key' }],
  title: 'Review section',
  editLink: 'linkToEditPath',
  useH4: true,
};

export default {
  title: 'Customer Components | ReviewSection',
};

export const Basic = () => <ReviewSection {...defaultProps} />;
