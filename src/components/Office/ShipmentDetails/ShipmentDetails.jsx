import React from 'react';
import { PropTypes } from 'prop-types';

import ShipmentDetailsMain from './ShipmentDetailsMain';
import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

import styles from 'components/Office/ShipmentDetails/ShipmentDetails.module.scss';
import { OrderShape } from 'types';
import { ShipmentShape } from 'types/shipment';

const ShipmentDetails = ({
  shipment,
  order,
  handleDivertShipment,
  handleRequestReweighModal,
  handleReviewSITExtension,
}) => {
  const { originDutyStation, destinationDutyStation, entitlement } = order;
  return (
    <div className={styles.ShipmentDetails}>
      <ShipmentDetailsMain
        className={styles.ShipmentDetailsMain}
        handleDivertShipment={handleDivertShipment}
        handleRequestReweighModal={handleRequestReweighModal}
        shipment={shipment}
        storageInTransit={entitlement.storageInTransit}
        dutyStationAddresses={{
          originDutyStationAddress: originDutyStation?.address,
          destinationDutyStationAddress: destinationDutyStation?.address,
        }}
        handleReviewSITExtension={handleReviewSITExtension}
      />
      <ShipmentDetailsSidebar
        className={styles.ShipmentDetailsSidebar}
        agents={shipment.agents}
        secondaryAddresses={shipment}
      />
    </div>
  );
};

ShipmentDetails.propTypes = {
  shipment: ShipmentShape.isRequired,
  order: OrderShape.isRequired,
  handleDivertShipment: PropTypes.func.isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
  handleReviewSITExtension: PropTypes.func.isRequired,
};

export default ShipmentDetails;
