import React from 'react';

import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

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
  <MockProviders permissions={[permissionTypes.createReweighRequest]}>
    <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} shipmentInfo={shipmentInfoNoReweigh} />
  </MockProviders>
);

export const WithDetailsReweighRequested = () => (
  <MockProviders permissions={[permissionTypes.createReweighRequest]}>
    <ShipmentWeightDetails estimatedWeight={1000} actualWeight={1000} shipmentInfo={shipmentInfoReweighRequested} />
  </MockProviders>
);

export const NTSRWithDetails = () => (
  <MockProviders permissions={[permissionTypes.createReweighRequest]}>
    <ShipmentWeightDetails
      estimatedWeight={null}
      actualWeight={1000}
      shipmentInfo={{ ...shipmentInfoNoReweigh, shipmentType: SHIPMENT_OPTIONS.NTSR }}
    />
  </MockProviders>
);
