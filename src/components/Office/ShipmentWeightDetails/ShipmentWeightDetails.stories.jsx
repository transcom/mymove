import React from 'react';

import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';

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
};

const shipmentInfoNoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
};

export const WithNoDetails = () => <ShipmentWeightDetails shipmentInfo={shipmentInfoNoReweigh} />;

export const WithDetailsNoReweighRequested = () => (
  <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} shipmentInfo={shipmentInfoNoReweigh} />
);

export const WithDetailsReweighRequested = () => (
  <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} shipmentInfo={shipmentInfoReweighRequested} />
);
