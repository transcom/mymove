import React from 'react';
import { PropTypes } from 'prop-types';
import { useParams } from 'react-router-dom';

import { formatQAReportID } from 'utils/formatters';
import styles from 'pages/Office/TXOMoveInfo/TXOTab.module.scss';
import EVALUATION_REPORT_TYPE from 'constants/evaluationReports';

const QaeReportHeader = ({ report }) => {
  const { moveCode } = useParams();

  if (!report || !report.type) {
    return null;
  }

  const isShipment = report.type === EVALUATION_REPORT_TYPE.SHIPMENT;
  const reportId = formatQAReportID(report.id);
  const mtoRefId = report.moveReferenceID;
  return (
    <div className={styles.pageHeader}>
      <h1>{`${isShipment ? 'Shipment' : 'Counseling'} report`}</h1>
      <div className={styles.pageHeaderDetails}>
        <h6>{`REPORT ID ${reportId}`}</h6>
        <h6>{`MOVE CODE #${moveCode}`}</h6>
        <h6>{`MTO REFERENCE ID #${mtoRefId}`}</h6>
      </div>
    </div>
  );
};

QaeReportHeader.propTypes = {
  report: PropTypes.object.isRequired,
};

export default QaeReportHeader;
