import React from 'react';

import ShipmentEvaluationReports from './ShipmentEvaluationReports';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Office Components/ShipmentEvaluationReports',
  component: ShipmentEvaluationReports,
};

const shipments = [
  {
    id: '1',
    shipmentType: SHIPMENT_OPTIONS.HHG,
  },
  {
    id: '2',
    shipmentType: SHIPMENT_OPTIONS.HHG,
  },
];

const reports = [
  {
    id: '12354',
    shipmentID: '1',
  },
];

export const empty = () => (
  <div className="officeApp">
    <ShipmentEvaluationReports reports={[]} />
  </div>
);

export const single = () => (
  <div className="officeApp">
    <ShipmentEvaluationReports shipments={shipments} reports={reports} />
  </div>
);
