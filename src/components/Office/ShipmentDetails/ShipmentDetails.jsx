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
  handleSubmitSITExtension,
  handleEditFacilityInfo,
  handleEditServiceOrderNumber,
  handleEditAccountingCodes,
}) => {
  const { originDutyLocation, destinationDutyLocation, entitlement } = order;
  const ordersLOA = {
    tac: order.tac,
    sac: order.sac,
    ntsTac: order.ntsTac,
    ntsSac: order.ntsSac,
  };

  return (
    <div className={styles.ShipmentDetails}>
      <ShipmentDetailsMain
        className={styles.ShipmentDetailsMain}
        handleDivertShipment={handleDivertShipment}
        handleRequestReweighModal={handleRequestReweighModal}
        shipment={shipment}
        entitilement={entitlement}
        dutyStationAddresses={{
          originDutyLocationAddress: originDutyLocation?.address,
          destinationDutyLocationAddress: destinationDutyLocation?.address,
        }}
        handleReviewSITExtension={handleReviewSITExtension}
        handleSubmitSITExtension={handleSubmitSITExtension}
      />
      <ShipmentDetailsSidebar
        className={styles.ShipmentDetailsSidebar}
        shipment={shipment}
        ordersLOA={ordersLOA}
        handleEditFacilityInfo={handleEditFacilityInfo}
        handleEditServiceOrderNumber={handleEditServiceOrderNumber}
        handleEditAccountingCodes={handleEditAccountingCodes}
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
  handleSubmitSITExtension: PropTypes.func.isRequired,
  handleEditFacilityInfo: PropTypes.func.isRequired,
  handleEditServiceOrderNumber: PropTypes.func.isRequired,
  handleEditAccountingCodes: PropTypes.func.isRequired,
};

export default ShipmentDetails;
