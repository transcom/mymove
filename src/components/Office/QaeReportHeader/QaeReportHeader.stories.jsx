import React from 'react';
import { MemoryRouter, Route } from 'react-router';

import QaeReportHeader from './QaeReportHeader';

import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';

export default {
  title: 'Office Components/QaeReportHeader',
  component: QaeReportHeader,
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={['/moves/AMDORD/evaluation-reports/ab30c135-1d6d-4a0d-a6d5-f408474f6ee2']}>
        <Route path="/moves/:moveCode/evaluation-reports/:reportId">
          <div style={{ padding: '40px', width: '900px', minWidth: '900px' }}>
            <Story />
          </div>
        </Route>
      </MemoryRouter>
    ),
  ],
};

export const ShipmentReport = () => (
  <QaeReportHeader
    report={{
      id: 'db30c135-1d6d-4a0d-a6d5-f408474f6ee2',
      type: EVALUATION_REPORT_TYPE.SHIPMENT,
      moveReferenceID: '6789-1234',
    }}
  />
);

export const CounselingReport = () => (
  <QaeReportHeader
    report={{
      id: 'ab30c135-1d6d-4a0d-a6d5-f408474f6ee2',
      type: EVALUATION_REPORT_TYPE.COUNSELING,
      moveReferenceID: '8789-1234',
    }}
  />
);
