import React, { useState } from 'react';
import { Button, Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { useLocation } from 'react-router';

import styles from './EvaluationReportTable.module.scss';

import ConnectedEvaluationReportConfirmationModal from 'components/ConfirmationModals/EvaluationReportConfirmationModal';
import { formatCustomerDate, formatEvaluationReportLocation, formatQAReportID } from 'utils/formatters';
import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';

const EvaluationReportTable = ({ reports, shipments, emptyText, moveCode, customerInfo, grade }) => {
  const location = useLocation();
  const [isViewReportModalVisible, setIsViewReportModalVisible] = useState(false);
  const [reportToView, setReportToView] = useState(undefined);

  const handleViewReportClick = (report) => {
    setReportToView(report);
    setIsViewReportModalVisible(true);
  };
  // Taken from https://mathiasbynens.github.io/rel-noopener/
  // tl;dr-- opening content in target _blank can leave parent window open to malicious code
  // below is a safer way to open content in a new tab
  function safeOpenInNewTab(url) {
    if (url) {
      const win = window.open();
      // win can be null if a pop-up blocker is used
      if (win) {
        win.opener = null;
        win.location = url;
      }
    }
  }

  const handleDownloadReportClick = (reportID) => {
    return safeOpenInNewTab(`/ghc/v1/evaluation-reports/${reportID}/download`);
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
        <td className={styles.downloadColumn}>
          <Button onClick={() => handleDownloadReportClick(report.id)}>Download</Button>
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
};

EvaluationReportTable.defaultProps = {
  reports: [],
  shipments: null,
};

export default EvaluationReportTable;
