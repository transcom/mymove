import React from 'react';
import PropTypes from 'prop-types';

import EvaluationReportTable from './EvaluationReportTable';
import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';

const ShipmentEvaluationReports = ({ shipments, reports }) => {
  const row = (shipment) => {
    return (
      <div key={shipment.id}>
        <EvaluationReportShipmentInfo shipment={shipment} />
        <EvaluationReportTable reports={reports.filter((r) => r.shipmentID === shipment.id)} />
      </div>
    );
  };

  const shipmentRows = shipments.map(row);

  return (
    <>
      <h2>Shipment QAE reports ({reports.length})</h2>
      {shipmentRows}
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
