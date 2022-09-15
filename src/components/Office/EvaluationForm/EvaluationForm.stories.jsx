import React from 'react';
import { MemoryRouter, Route } from 'react-router';

import EvaluationForm from './EvaluationForm';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/EvaluationForm',
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={['/moves/AMDORD/evaluation-reports/6739d7fc-6067-4e84-996d-f4f70b8ec6fd']}>
        <Route path="/moves/:moveCode/evaluation-reports/:reportId">
          <MockProviders>
            <div style={{ padding: '40px', width: '850px' }}>
              <Story />
            </div>
          </MockProviders>
        </Route>
      </MemoryRouter>
    ),
  ],
};

export const Default = () => <EvaluationForm evaluationReport={{ type: 'SHIPMENT' }} />;
