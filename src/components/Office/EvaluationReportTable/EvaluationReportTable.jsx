import React, { useState } from 'react';
import { Button, Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { useLocation } from 'react-router';

import styles from './EvaluationReportTable.module.scss';

import ConnectedEvaluationReportConfirmationModal from 'components/ConfirmationModals/EvaluationReportConfirmationModal';
import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import { formatCustomerDate, formatEvaluationReportLocation, formatQAReportID } from 'utils/formatters';
import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';

const EvaluationReportTable = ({
  reports,
  shipments,
  emptyText,
  moveCode,
  customerInfo,
  grade,
  setReportToDelete,
  setIsDeleteModalOpen,
  deleteReport,
  isDeleteModalOpen,
}) => {
  const location = useLocation();
  const [isViewReportModalVisible, setIsViewReportModalVisible] = useState(false);
  const [reportToView, setReportToView] = useState(undefined);

  // whether or not the delete report modal is displaying
  const toggleDeleteReportModal = (reportID) => {
    setReportToDelete(reports.find((report) => report.id === reportID));
    setIsDeleteModalOpen(!isDeleteModalOpen);
  };

  const handleViewReportClick = (report) => {
    setReportToView(report);
    setIsViewReportModalVisible(true);
  };

  // this handles the close button at the bottom of the view report modal
  const toggleCloseModal = () => {
    setIsViewReportModalVisible(!isViewReportModalVisible);
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
          {!report.submittedAt && <a href={`${location.pathname}/${report.id}`}>Edit report</a>}
        </td>
        {report.submittedAt && (
          <td className={styles.downloadColumn}>
            <a href={`/ghc/v1/evaluation-reports/${report.id}/download`} target="_blank" rel="noopener noreferrer">
              Download
            </a>
          </td>
        )}
        {!report.submittedAt && (
          <td className={styles.downloadColumn}>
            <Button className="usa-button--unstyled" onClick={() => toggleDeleteReportModal(report.id)}>
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
    <div>
      <ConnectedDeleteEvaluationReportConfirmationModal
        isOpen={isDeleteModalOpen}
        closeModal={toggleDeleteReportModal}
        submitModal={deleteReport}
        isDeleteFromTable
      />
      {isViewReportModalVisible && reportToView && (
        <ConnectedEvaluationReportConfirmationModal
          isOpen={isViewReportModalVisible}
          evaluationReport={reportToView}
          moveCode={moveCode}
          customerInfo={customerInfo}
          grade={grade}
          mtoShipments={shipments}
          modalActions={
            <div className={styles.modalActions}>
              <Button
                type="button"
                onClick={toggleCloseModal}
                aria-label="Close"
                secondary
                className={styles.closeModalBtn}
              >
                Close
              </Button>
            </div>
          }
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
  shipments: PropTypes.arrayOf(ShipmentShape),
  setIsDeleteModalOpen: PropTypes.func.isRequired,
  setReportToDelete: PropTypes.func.isRequired,
  deleteReport: PropTypes.func.isRequired,
  isDeleteModalOpen: PropTypes.bool.isRequired,
};

EvaluationReportTable.defaultProps = {
  reports: [],
  shipments: null,
};

export default EvaluationReportTable;
