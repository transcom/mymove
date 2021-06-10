import React from 'react';

import ShipmentDetailsMain from './ShipmentDetailsMain';
import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

import styles from 'components/Office/ShipmentDetails/ShipmentDetails.module.scss';
import { MTOShipmentShape, OrderShape } from 'types';

const ShipmentDetails = ({ shipment, order }) => {
  return (
    <div className={styles.ShipmentDetails}>
      <ShipmentDetailsMain
        className={styles.ShipmentDetailsMain}
        shipment={{
          requestedPickupDate: shipment.requestedPickupDate,
          scheduledPickupDate: shipment.scheduledPickupDate,
          pickupAddress: shipment.pickupAddress,
          destinationAddress: shipment.destinationAddress,
          primeEstimatedWeight: shipment.primeEstimatedWeight,
          primeActualWeight: shipment.primeActualWeight,
        }}
        order={{
          originDutyStationAddress: order.originDutyStation?.address,
          destinationDutyStationAddress: order.destinationDutyStation?.address,
        }}
      />
      <ShipmentDetailsSidebar
        className={styles.ShipmentDetailsSidebar}
        agents={shipment.mtoAgents || shipment.agents}
        secondaryAddresses={{
          secondaryPickupAddress: shipment.secondaryPickupAddress,
          secondaryDeliveryAddress: shipment.secondaryDeliveryAddress,
        }}
      />
    </div>
  );
};

ShipmentDetails.propTypes = {
  shipment: MTOShipmentShape.isRequired,
  order: OrderShape.isRequired,
};

export default ShipmentDetails;
