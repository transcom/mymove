import React from 'react';
import { Button, Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { useLocation, useNavigate } from 'react-router-dom';

import styles from './QAEReportTable.module.scss';

import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import { formatCustomerDate, formatEvaluationReportLocation, formatQAReportID } from 'utils/formatters';
import { EvaluationReportShape } from 'types';

const QAEReportTable = ({
  reports,
  emptyText,
  moveCode,
  setReportToDelete,
  setIsDeleteModalOpen,
  deleteReport,
  isDeleteModalOpen,
}) => {
  const location = useLocation();
  const navigate = useNavigate();

  // whether or not the delete report modal is displaying
  const toggleDeleteReportModal = (reportID) => {
    setReportToDelete(reports.find((report) => report.id === reportID));
    setIsDeleteModalOpen(!isDeleteModalOpen);
  };

  const handleViewReportClick = (reportID) => {
    navigate(`/moves/${moveCode}/evaluation-report/${reportID}`);
  };

  const row = (report) => {
    let violations;
    let seriousIncident;

    if (report.violationsObserved) {
      violations = 'Yes';
    } else if (report.violationsObserved === false) {
      violations = 'No';
    } else {
      violations = '';
    }

    if (report.seriousIncident) {
      seriousIncident = 'Yes';
    } else if (report.seriousIncident === false) {
      seriousIncident = 'No';
    } else {
      seriousIncident = '';
    }
    return (
      <tr key={report.id}>
        <td className={styles.reportIDColumn}>
          {formatQAReportID(report.id)} {report.submittedAt ? null : <Tag className={styles.draftTag}>DRAFT</Tag>}
        </td>
        <td className={styles.dateSubmittedColumn}>{report.submittedAt && formatCustomerDate(report.submittedAt)}</td>
        <td className={styles.locationColumn}>{formatEvaluationReportLocation(report.location)}</td>
        <td data-testid="violation-column" className={styles.violationsColumn}>
          {violations}
        </td>
        <td data-testid="incident-column" className={styles.seriousIncidentColumn}>
          {seriousIncident}
        </td>
        <td className={styles.viewReportColumn}>
          {report.submittedAt && (
            <Button
              type="button"
              id={report.id}
              className={classnames(styles.viewButton, 'text-blue usa-button--unstyled')}
              onClick={() => handleViewReportClick(report.id)}
              data-testid="viewReport"
            >
              View report
            </Button>
          )}
          {!report.submittedAt && (
            <a href={`${location.pathname}/${report.id}`} data-testid="editReport">
              Edit report
            </a>
          )}
        </td>
        {report.submittedAt && (
          <td className={styles.downloadColumn}>
            <a
              href={`/ghc/v1/evaluation-reports/${report.id}/download`}
              target="_blank"
              rel="noopener noreferrer"
              data-testid="downloadReport"
            >
              Download
            </a>
          </td>
        )}
        {!report.submittedAt && (
          <td className={styles.downloadColumn}>
            <Button
              className="usa-button--unstyled"
              onClick={() => toggleDeleteReportModal(report.id)}
              data-testid="deleteReport"
            >
              Delete
            </Button>
          </td>
        )}
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
    <div data-testid="evaluationReportTable">
      <ConnectedDeleteEvaluationReportConfirmationModal
        isOpen={isDeleteModalOpen}
        closeModal={toggleDeleteReportModal}
        submitModal={deleteReport}
        isDeleteFromTable
      />
      <table className={styles.qaeReportTable}>
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

QAEReportTable.propTypes = {
  reports: PropTypes.arrayOf(EvaluationReportShape),
  emptyText: PropTypes.string.isRequired,
  moveCode: PropTypes.string.isRequired,
  setIsDeleteModalOpen: PropTypes.func.isRequired,
  setReportToDelete: PropTypes.func.isRequired,
  deleteReport: PropTypes.func.isRequired,
  isDeleteModalOpen: PropTypes.bool.isRequired,
};

QAEReportTable.defaultProps = {
  reports: [],
};

export default QAEReportTable;
