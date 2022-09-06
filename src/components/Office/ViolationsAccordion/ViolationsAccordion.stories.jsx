import React from 'react';

import ViolationsAccordion from './ViolationsAccordion';

export default {
  title: 'Office Components/ViolationsAccordion',
  component: ViolationsAccordion,
  decorators: [
    (Story) => (
      <div style={{ padding: '40px', width: '850px', minWidth: '850px' }}>
        <Story />
      </div>
    ),
  ],
};

const violations = [
  {
    category: 'Pre-Move Services',
    displayOrder: 1,
    id: 'c359ebc3-a506-4f41-8f91-409d59c97b22',
    paragraphNumber: '1.2.5.1',
    requirementStatement:
      'The contractor shall assign, during initial communication with each customer, a single POC responsible for coordination and communication throughout all phases of the move. The POC’s contact information shall be maintained throughout the entire shipment process and until all associated actions are final.',
    requirementSummary: 'Provide a single point of contact (POC)',
    subCategory: 'Customer Support',
    title: 'Point of Contact (POC)',
  },
];

export const Default = () => <ViolationsAccordion violations={violations} />;
