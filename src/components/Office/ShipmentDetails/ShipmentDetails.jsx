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
  handleShowDiversionModal,
  handleRequestReweighModal,
  handleReviewSITExtension,
  handleSubmitSITExtension,
  handleEditFacilityInfo,
  handleEditServiceOrderNumber,
  handleEditAccountingCodes,
  handleUpdateSITServiceItemCustomerExpense,
  isMoveLocked,
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
        handleShowDiversionModal={handleShowDiversionModal}
        handleRequestReweighModal={handleRequestReweighModal}
        shipment={shipment}
        entitilement={entitlement}
        dutyLocationAddresses={{
          originDutyLocationAddress: originDutyLocation?.address,
          destinationDutyLocationAddress: destinationDutyLocation?.address,
        }}
        handleReviewSITExtension={handleReviewSITExtension}
        handleSubmitSITExtension={handleSubmitSITExtension}
        handleUpdateSITServiceItemCustomerExpense={handleUpdateSITServiceItemCustomerExpense}
        isMoveLocked={isMoveLocked}
      />
      <ShipmentDetailsSidebar
        className={styles.ShipmentDetailsSidebar}
        shipment={shipment}
        ordersLOA={ordersLOA}
        handleEditFacilityInfo={handleEditFacilityInfo}
        handleEditServiceOrderNumber={handleEditServiceOrderNumber}
        handleEditAccountingCodes={handleEditAccountingCodes}
        isMoveLocked={isMoveLocked}
      />
    </div>
  );
};

ShipmentDetails.propTypes = {
  shipment: ShipmentShape.isRequired,
  order: OrderShape.isRequired,
  handleShowDiversionModal: PropTypes.func.isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
  handleReviewSITExtension: PropTypes.func.isRequired,
  handleSubmitSITExtension: PropTypes.func.isRequired,
  handleEditFacilityInfo: PropTypes.func.isRequired,
  handleEditServiceOrderNumber: PropTypes.func.isRequired,
  handleEditAccountingCodes: PropTypes.func.isRequired,
  handleUpdateSITServiceItemCustomerExpense: PropTypes.func.isRequired,
};

export default ShipmentDetails;
