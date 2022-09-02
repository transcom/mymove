import React from 'react';

import SelectedViolation from './SelectedViolation';

export default {
  title: 'Office Components/SelectedViolation',
  component: SelectedViolation,
};

const mockViolation = {
  category: 'Category 1',
  displayOrder: 1,
  id: '9cdc8dc3-6cf4-46fb-b272-1468ef40796f',
  paragraphNumber: '1.2.3',
  requirementStatement: 'Test requirement statement for violation 1',
  requirementSummary: 'Test requirement summary for violation 1',
  subCategory: 'SubCategory 1',
  title: 'Title for violation 1',
};

export const Default = () => <SelectedViolation violation={mockViolation} />;
