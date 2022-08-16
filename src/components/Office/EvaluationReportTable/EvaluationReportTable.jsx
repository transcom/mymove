import React, { useState } from 'react';
import { Button, Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { useLocation } from 'react-router';

import { EvaluationReportShape } from '../../../types/evaluationReport';

import styles from './EvaluationReportTable.module.scss';
import EvaluationReportContainer from './EvaluationReportContainer';

import { formatCustomerDate, formatEvaluationReportLocation, formatQAReportID } from 'utils/formatters';
import { CustomerShape } from 'types';

const EvaluationReportTable = ({ reports, emptyText, moveCode, customerInfo, grade, shipmentId }) => {
  const { pathname } = useLocation();
  const [isViewReportModalVisible, setIsViewReportModalVisible] = useState(false);
  const [reportToView, setReportToView] = useState(undefined);

  const handleViewReportClick = (report) => {
    setReportToView(report);
    setIsViewReportModalVisible(true);
  };

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
          {report.submittedAt && (
            <Button
              type="button"
              id={report.id}
              className={classnames(styles.viewButton, 'text-blue usa-button--unstyled')}
              onClick={() => handleViewReportClick(report)}
            >
              View report
            </Button>
          )}
          {!report.submittedAt && <a href={`${pathname}/${report.id}`}>Edit report</a>}
        </td>
        <td className={styles.downloadColumn}>
          <a href={`${pathname}/evaluation-reports/${report.id}/download`}>Download</a>
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
    <div>
      {isViewReportModalVisible && reportToView && (
        <EvaluationReportContainer
          reportType={reportToView.type}
          evaluationReportId={reportToView.id}
          moveCode={moveCode}
          customerInfo={customerInfo}
          grade={grade}
          shipmentId={shipmentId}
          setIsModalVisible={setIsViewReportModalVisible}
        />
      )}
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
    </div>
  );
};

EvaluationReportTable.propTypes = {
  reports: PropTypes.arrayOf(EvaluationReportShape),
  emptyText: PropTypes.string.isRequired,
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  shipmentId: PropTypes.string,
};

EvaluationReportTable.defaultProps = {
  reports: [],
  shipmentId: '',
};

export default EvaluationReportTable;
