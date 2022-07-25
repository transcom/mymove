import React from 'react';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './EvaluationReportTable.module.scss';

import { formatCustomerDate, formatEvaluationReportLocation, formatQAReportID } from 'utils/formatters';
import { EvaluationReportShape } from 'types/evaluationReport';

const EvaluationReportTable = ({ reports, emptyText }) => {
  const row = (report) => {
    return (
      <tr key={report.id}>
        <td className={styles.reportIDColumn}>
          {formatQAReportID(report.id)} {report.submittedAt ? null : <Tag className={styles.draftTag}>DRAFT</Tag>}
        </td>
        <td className={styles.dateSubmittedColumn}>{report.submittedAt && formatCustomerDate(report.submittedAt)}</td>
        <td className={styles.locationColumn}>{formatEvaluationReportLocation(report.location)}</td>
        <td className={styles.violationsColumn}>{report.violationsObserved ? 'Yes' : 'No'}</td>
        <td className={styles.seriousIncidentColumn}>No</td>
        <td className={styles.viewReportColumn}>
          <a href={`/moves/${report.moveID}/evaluation-reports/${report.id}`}>View report</a>
        </td>
        <td className={styles.downloadColumn}>
          <a href={`/moves/${report.moveID}/evaluation-reports/${report.id}/download`}>Download</a>
        </td>
      </tr>
    );
  };
  let tableRows = (
    <tr className={styles.emptyTableRow}>
      <td className={styles.emptyTableRow} colSpan={7}>
        {emptyText}
      </td>
    </tr>
  );
  if (reports.length > 0) {
    tableRows = reports.map(row);
  }

  return (
    <table className={styles.evaluationReportTable}>
      <thead>
        <tr>
          <th className={styles.reportIDColumn}>Report ID</th>
          <th className={styles.dateSubmittedColumn}>Date submitted</th>
          <th className={styles.locationColumn}>Location</th>
          <th className={styles.violationsColumn}>Violations</th>
          <th className={styles.seriousIncidentColumn}>Serious Incident</th>
          <th className={styles.viewReportColumn} aria-label="View report" />
          <th className={styles.downloadColumn} aria-label="Download" />
        </tr>
      </thead>
      <tbody>{tableRows}</tbody>
    </table>
  );
};

EvaluationReportTable.propTypes = {
  reports: PropTypes.arrayOf(EvaluationReportShape),
  emptyText: PropTypes.string.isRequired,
};

EvaluationReportTable.defaultProps = {
  reports: [],
};

export default EvaluationReportTable;
