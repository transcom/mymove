import React from 'react';
import { MemoryRouter, Route } from 'react-router';

import ShipmentEvaluationForm from './ShipmentEvaluationForm';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/ShipmentEvaluationForm',
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={['/moves/AMDORD/customer-support-remarks']}>
        <Route path="/moves/:moveCode/customer-support-remarks">
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

export const Default = () => <ShipmentEvaluationForm />;
