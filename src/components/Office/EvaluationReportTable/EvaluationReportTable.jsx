import React from 'react';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import { formatCustomerDate, formatEvaluationReportLocation, formatQAReportID } from 'utils/formatters';
import { EvaluationReportShape } from 'types/evaluationReport';

const EvaluationReportTable = ({ reports }) => {
  const row = (report) => {
    return (
      <tr key={report.id}>
        <td>
          {formatQAReportID(report.id)} {report.submittedAt ? null : <Tag>DRAFT</Tag>}
        </td>
        <td>{formatCustomerDate(report.submittedAt)}</td>
        <td>{formatEvaluationReportLocation(report.location)}</td>
        <td>{report.violations ? 'Yes' : 'No'}</td>
        <td>No</td>
        <td>
          <a href={`/moves/${report.moveID}/evaluation-reports/${report.id}`}>View report</a>
        </td>
        <td>
          <a href={`/moves/${report.moveID}/evaluation-reports/${report.id}/download`}>Download</a>
        </td>
      </tr>
    );
  };
  let tableRows = (
    <tr>
      <td colSpan={5}>No QAE reports have been submitted for this shipment</td>
    </tr>
  );
  if (reports.length > 0) {
    tableRows = reports.map(row);
  }

  return (
    <table>
      <thead>
        <tr>
          <th>Report ID</th>
          <th>Date submitted</th>
          <th>Location</th>
          <th>Violations</th>
          <th>Serious Incident</th>
          <th>View Report</th>
          <th>Download</th>
        </tr>
      </thead>
      <tbody>{tableRows}</tbody>
    </table>
  );
};

EvaluationReportTable.propTypes = {
  reports: PropTypes.arrayOf(EvaluationReportShape),
};

EvaluationReportTable.defaultProps = {
  reports: [],
};

export default EvaluationReportTable;
