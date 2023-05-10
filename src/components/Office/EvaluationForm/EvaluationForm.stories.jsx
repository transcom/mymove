import React from 'react';

import EvaluationForm from './EvaluationForm';

import { MockProviders } from 'testUtils';
import { qaeCSRRoutes } from 'constants/routes';

export default {
  title: 'Office Components/EvaluationForm',
  decorators: [
    (Story) => (
      <MockProviders
        path={qaeCSRRoutes.BASE_EVALUATION_REPORT_PATH}
        params={{ moveCode: 'AMDORD', reportId: '6739d7fc-6067-4e84-996d-f4f70b8ec6fd' }}
      >
        <div style={{ padding: '40px', width: '850px' }}>
          <Story />
        </div>
      </MockProviders>
    ),
  ],
};

export const Default = () => <EvaluationForm evaluationReport={{ type: 'SHIPMENT' }} />;
