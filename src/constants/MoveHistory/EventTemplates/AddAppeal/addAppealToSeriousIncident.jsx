import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import { formatQAReportID } from 'utils/formatters';

export default {
  action: a.INSERT,
  eventName: o.addAppealToSeriousIncident,
  tableName: t.gsr_appeals,
  getEventNameDisplay: () => 'Appeal Decision on Serious Incident',
  getDetails: ({ changedValues, context }) => (
    <div data-testid="seriousIncidentAppealInfo">
      <b>Report ID</b>: {formatQAReportID(changedValues.evaluation_report_id)}
      <br />
      <b>Report Type</b>: {context[0].evaluation_report_type}
      <br />
      <b>Remarks</b>: {changedValues.remarks}
      <br />
      <b>Status</b>: {changedValues.appeal_status}
      <br />
    </div>
  ),
};
