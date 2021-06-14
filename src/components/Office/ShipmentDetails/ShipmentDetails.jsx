import React from 'react';
import { PropTypes } from 'prop-types';

import ShipmentDetailsMain from './ShipmentDetailsMain';
import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

import styles from 'components/Office/ShipmentDetails/ShipmentDetails.module.scss';
import { OrderShape } from 'types';
import { ShipmentShape } from 'types/shipment';

const ShipmentDetails = ({ shipment, order, handleDivertShipment }) => {
  const { originDutyStation, destinationDutyStation } = order;
  return (
    <div className={styles.ShipmentDetails}>
      <ShipmentDetailsMain
        className={styles.ShipmentDetailsMain}
        handleDivertShipment={handleDivertShipment}
        shipment={shipment}
        dutyStationAddresses={{
          originDutyStationAddress: originDutyStation?.address,
          destinationDutyStationAddress: destinationDutyStation?.address,
        }}
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
};

export default ShipmentDetails;
