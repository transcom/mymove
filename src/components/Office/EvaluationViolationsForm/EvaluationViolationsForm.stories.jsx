import React from 'react';
import { MemoryRouter, Route } from 'react-router';

import EvaluationViolationsForm from './EvaluationViolationsForm';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/EvaluationViolationsForm',
  component: EvaluationViolationsForm,
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={['/moves/AMDORD/customer-support-remarks']}>
        <Route path="/moves/:moveCode/customer-support-remarks">
          <MockProviders>
            <div style={{ padding: '40px', width: '950px', minWidth: '950px' }}>
              <Story />
            </div>
          </MockProviders>
        </Route>
      </MemoryRouter>
    ),
  ],
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

export const Default = () => <EvaluationViolationsForm violations={[mockViolation]} />;
