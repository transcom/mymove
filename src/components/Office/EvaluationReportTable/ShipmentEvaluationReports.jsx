import React from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';

import EvaluationReportTable from './EvaluationReportTable';
import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';
import styles from './ShipmentEvaluationReports.module.scss';

import { CustomerShape } from 'types';

const ShipmentEvaluationReports = ({ shipments, reports, moveCode, customerInfo, grade }) => {
  const sortedShipments = shipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

  const shipmentRows = sortedShipments.map((shipment) => {
    return (
      <div key={shipment.id} className={styles.shipmentRow}>
        <EvaluationReportShipmentInfo shipment={shipment} />
        <EvaluationReportTable
          moveCode={moveCode}
          reports={reports.filter((r) => r.shipmentID === shipment.id)}
          customerInfo={customerInfo}
          grade={grade}
          shipmentId={shipment.id}
          emptyText="No QAE reports have been submitted for this shipment."
        />
      </div>
    );
  });

  return (
    <div className={styles.shipmentEvaluationReportsContainer}>
      <h2>Shipment QAE reports ({reports.length})</h2>
      <div className={styles.shipmentReportRows}>{shipmentRows}</div>
    </div>
  );
};

ShipmentEvaluationReports.propTypes = {
  reports: PropTypes.arrayOf(PropTypes.object),
  shipments: PropTypes.arrayOf(PropTypes.object),
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
};
ShipmentEvaluationReports.defaultProps = {
  reports: [],
  shipments: [],
};

export default ShipmentEvaluationReports;
