import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import { formatQAReportID } from 'utils/formatters';

export default {
  action: a.INSERT,
  eventName: o.addAppealToViolation,
  tableName: t.gsr_appeals,
  getEventNameDisplay: ({ context }) => {
    return (
      <>
        <div>Appeal Decision on Violation</div>
        <div data-testid="violationAppealTitle">
          {context[0]?.violation_paragraph_number} {context[0]?.violation_title}
        </div>
      </>
    );
  },
  getDetails: ({ changedValues, context }) => {
    return (
      <div data-testid="violationAppealInfo">
        <b>Violation Summary</b>: {context[0]?.violation_summary}
        <br />
        <b>Report ID</b>: {formatQAReportID(changedValues.evaluation_report_id)}
        <br />
        <b>Report Type</b>: {context[0].evaluation_report_type === 'SHIPMENT' ? 'Shipment' : 'Counseling'}
        <br />
        <b>Remarks</b>: {changedValues.remarks}
        <br />
        <b>Status</b>: {changedValues.appeal_status === 'SUSTAINED' ? 'Sustained' : 'Rejected'}
        <br />
      </div>
    );
  },
};
