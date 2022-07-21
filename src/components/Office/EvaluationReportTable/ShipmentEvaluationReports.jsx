import React from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';

import EvaluationReportTable from './EvaluationReportTable';
import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';
import styles from './ShipmentEvaluationReports.module.scss';

const ShipmentEvaluationReports = ({ shipments, reports }) => {
  const sortedShipments = shipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));
  const shipmentNumbersByType = {};

  const shipmentRows = sortedShipments.map((shipment) => {
    const { shipmentType } = shipment;
    if (shipmentNumbersByType[shipmentType]) {
      shipmentNumbersByType[shipmentType] += 1;
    } else {
      shipmentNumbersByType[shipmentType] = 1;
    }
    const shipmentNumber = shipmentNumbersByType[shipmentType];
    return (
      <div key={shipment.id} className={styles.shipmentRow}>
        <EvaluationReportShipmentInfo shipment={shipment} shipmentNumber={shipmentNumber} />
        <EvaluationReportTable reports={reports.filter((r) => r.shipmentID === shipment.id)} />
      </div>
    );
  });

  return (
    <>
      <h2>Shipment QAE reports ({reports.length})</h2>
      <div className={styles.gridContainer}>{shipmentRows}</div>
    </>
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
