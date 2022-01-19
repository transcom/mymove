import React from 'react';

import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Office Components/ShipmentWeightDetails',
  decorators: [
    (storyFn) => (
      <div className="officeApp" id="containers" style={{ padding: '20px' }}>
        {storyFn()}
      </div>
    ),
  ],
};

const shipmentInfoReweighRequested = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  reweighID: 'reweighRequestID',
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

const shipmentInfoNoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

export const WithNoDetails = () => <ShipmentWeightDetails shipmentInfo={shipmentInfoNoReweigh} />;

export const WithDetailsNoReweighRequested = () => (
  <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} shipmentInfo={shipmentInfoNoReweigh} />
);

export const WithDetailsReweighRequested = () => (
  <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} shipmentInfo={shipmentInfoReweighRequested} />
);

export const NTSRWithDetails = () => (
  <ShipmentWeightDetails
    estimatedWeight={null}
    actualWeight={1000}
    shipmentInfo={{ ...shipmentInfoNoReweigh, shipmentType: SHIPMENT_OPTIONS.NTSR }}
  />
);
