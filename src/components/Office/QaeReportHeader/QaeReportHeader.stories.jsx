import React from 'react';

import QaeReportHeader from './QaeReportHeader';

import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';
import { MockRouting } from 'testUtils';
import { qaeCSRRoutes } from 'constants/routes';

export default {
  title: 'Office Components/QaeReportHeader',
  component: QaeReportHeader,
  decorators: [
    (Story) => (
      <MockRouting
        path={qaeCSRRoutes.BASE_EVALUATION_REPORT_PATH}
        params={{ moveCode: 'AMDORD', reportId: 'ab30c135-1d6d-4a0d-a6d5-f408474f6ee2' }}
      >
        <div style={{ padding: '40px', width: '900px', minWidth: '900px' }}>
          <Story />
        </div>
      </MockRouting>
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
