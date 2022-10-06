import React from 'react';
import { MemoryRouter, Route } from 'react-router';

import EvaluationViolationsForm from './EvaluationViolationsForm';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/EvaluationViolationsForm',
  component: EvaluationViolationsForm,
  decorators: [
    (Story) => (
      <MemoryRouter
        initialEntries={['/moves/REWAYD/evaluation-reports/09132c3b-3ffe-41ec-9393-16e6f074adf7/violations']}
      >
        <Route path="/moves/:moveCode/evaluation-reports/:reportId/violations">
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

const mockEvaluationReport = {
  id: '9cdc8dc3-6cf4-46fb-b272-1468ef40796f',
  createdAt: '2022-09-07T15:17:37.484Z',
  eTag: 'MjAyMi0wOS0wN1QxODowNjozNy44NjQxNDJa',
  evaluationLengthMinutes: 240,
  inspectionDate: '2022-09-08',
  inspectionType: 'DATA_REVIEW',
  location: 'A12345',
  moveID: '551dd01f-90cf-44d6-addb-ff919433dd61',
  moveReferenceID: '4118-8295',
  officeUser: {
    email: 'qae_csr_role@office.mil',
    firstName: 'Leo',
    id: 'ef4f6d1f-4ac3-4159-a364-5403e7d958ff',
    lastName: 'Spaceman',
    phone: '415-555-1212',
  },
  remarks: 'test',
  shipmentID: '319e0751-1337-4ed9-b4b5-a15d4e6d272c',
  type: 'SHIPMENT',
  updatedAt: '2022-09-07T18:06:37.864Z',
  violationsObserved: true,
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

const mockCustomerInfo = {
  agency: 'ARMY',
  grade: 'E_4',
};

export const Default = () => (
  <EvaluationViolationsForm
    violations={[mockViolation]}
    evaluationReport={mockEvaluationReport}
    customerInfo={mockCustomerInfo}
    reportViolations={null}
  />
);
