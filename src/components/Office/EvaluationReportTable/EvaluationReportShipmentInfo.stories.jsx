import React from 'react';

import { SHIPMENT_OPTIONS } from '../../../shared/constants';

import EvaluationReportShipmentInfo from './EvaluationReportShipmentInfo';

export default {
  title: 'Office Components/EvaluationReportShipmentInfo',
  component: EvaluationReportShipmentInfo,
};

const hhgShipment = { id: '111111', shipmentType: SHIPMENT_OPTIONS.HHG };
const ppmShipment = { id: '22222', shipmentType: SHIPMENT_OPTIONS.PPM };

const ntsShipment = { id: '3333333', shipmentType: SHIPMENT_OPTIONS.NTS };
const ntsrShipment = { id: '444444', shipmentType: SHIPMENT_OPTIONS.NTSR };

export const hhg = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={hhgShipment} />
  </div>
);

export const nts = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={ntsShipment} />
  </div>
);
export const ntsr = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={ntsrShipment} />
  </div>
);
export const ppm = () => (
  <div className="officeApp">
    <EvaluationReportShipmentInfo shipment={ppmShipment} />
  </div>
);
