import React from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';

import EvaluationReportTable from './EvaluationReportTable';
import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';
import styles from './ShipmentEvaluationReports.module.scss';

const ShipmentEvaluationReports = ({ shipments, reports }) => {
  const sortedShipments = shipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

  const shipmentRows = sortedShipments.map((shipment) => {
    return (
      <div key={shipment.id} className={styles.shipmentRow}>
        <EvaluationReportShipmentInfo shipment={shipment} />
        <EvaluationReportTable
          reports={reports.filter((r) => r.shipmentID === shipment.id)}
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
};
ShipmentEvaluationReports.defaultProps = {
  reports: [],
  shipments: [],
};

export default ShipmentEvaluationReports;
