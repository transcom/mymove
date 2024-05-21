import React from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';

import EvaluationReportTable from './QAEReportTable';
import ShipmentQAEReportHeader from './ShipmentQAEReportHeader';
import styles from './ShipmentQAEReports.module.scss';

import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';

const ShipmentQAEReports = ({
  shipments,
  reports,
  moveCode,
  customerInfo,
  grade,
  setReportToDelete,
  setIsDeleteModalOpen,
  deleteReport,
  isDeleteModalOpen,
  destinationDutyLocationPostalCode,
  isMoveLocked,
}) => {
  const sortedShipments = shipments.sort((a, b) => moment(a.createdAt) - moment(b.createdAt));

  const shipmentRows = sortedShipments.map((shipment) => {
    return (
      <div key={shipment.id} className={styles.shipmentRow}>
        <ShipmentQAEReportHeader
          shipment={shipment}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
          isMoveLocked={isMoveLocked}
        />
        <EvaluationReportTable
          moveCode={moveCode}
          reports={reports.filter((r) => r.shipmentID === shipment.id)}
          customerInfo={customerInfo}
          grade={grade}
          shipments={[shipment]}
          emptyText="No QAE reports have been submitted for this shipment."
          setReportToDelete={setReportToDelete}
          setIsDeleteModalOpen={setIsDeleteModalOpen}
          isDeleteModalOpen={isDeleteModalOpen}
          deleteReport={deleteReport}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
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

ShipmentQAEReports.propTypes = {
  reports: PropTypes.arrayOf(EvaluationReportShape),
  shipments: PropTypes.arrayOf(ShipmentShape),
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  destinationDutyLocationPostalCode: PropTypes.string.isRequired,
};
ShipmentQAEReports.defaultProps = {
  reports: [],
  shipments: [],
};

export default ShipmentQAEReports;
