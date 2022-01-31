import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from 'components/Office/ShipmentDetails/ShipmentDetailsSidebar.module.scss';
import SimpleSection from 'containers/SimpleSection/SimpleSection';
import ConnectedEditFacilityInfoModal from 'components/Office/EditFacilityInfoModal/EditFacilityInfoModal';
import ConnectedServiceOrderNumberModal from 'components/Office/ServiceOrderNumberModal/ServiceOrderNumberModal';
import { retrieveSAC, retrieveTAC, formatAgent, formatAddress, formatAccountingCode } from 'utils/shipmentDisplay';
import { ShipmentShape } from 'types/shipment';
import { OrdersLOAShape } from 'types/order';

const ShipmentDetailsSidebar = ({
  className,
  shipment,
  ordersLOA,
  handleEditFacilityInfo,
  handleEditServiceOrderNumber,
}) => {
  const { mtoAgents, secondaryAddresses, serviceOrderNumber, storageFacility, sacType, tacType } = shipment;
  const tac = retrieveTAC(shipment.tacType, ordersLOA);
  const sac = retrieveSAC(shipment.sacType, ordersLOA);

  const [isEditFacilityInfoModalVisible, setIsEditFacilityInfoModalVisible] = useState(false);
  const [isSonModalVisible, setIsSonModalVisible] = useState(false);

  const handleShowEditFacilityInfoModal = () => {
    setIsEditFacilityInfoModalVisible(true);
  };

  const handleShowSonModal = () => setIsSonModalVisible(true);
  const handleCloseSonModal = () => setIsSonModalVisible(false);
  const handleSubmitSonModal = (values) => {
    handleEditServiceOrderNumber(values, shipment);
    setIsSonModalVisible(false);
  };

  return (
    <div className={className}>
      <ConnectedEditFacilityInfoModal
        isOpen={isEditFacilityInfoModalVisible}
        onSubmit={(e) => {
          handleEditFacilityInfo(e, shipment);
          setIsEditFacilityInfoModalVisible(false);
        }}
        onClose={() => {
          setIsEditFacilityInfoModalVisible(false);
        }}
        storageFacility={shipment.storageFacility}
        serviceOrderNumber={shipment.serviceOrderNumber}
        shipmentType={shipment.shipmentType}
      />

      <ConnectedServiceOrderNumberModal
        isOpen={isSonModalVisible}
        onSubmit={handleSubmitSonModal}
        onClose={handleCloseSonModal}
        serviceOrderNumber={shipment.serviceOrderNumber}
      />

      {mtoAgents &&
        mtoAgents.map((agent) => (
          <SimpleSection
            key={`${agent.agentType}-${agent.email}`}
            header={agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}
            border
          >
            <div>{formatAgent(agent)}</div>
          </SimpleSection>
        ))}

      {storageFacility && storageFacility.facilityName && (
        <SimpleSection
          key="facility-info-and-address"
          header={
            <div className={styles.ShipmentDetailsSidebar}>
              Facility info and address
              <Button
                size="small"
                type="button"
                onClick={handleShowEditFacilityInfoModal}
                className="float-right usa-link modal-link"
                data-testid="edit-facility-info-modal-open"
                unstyled
              >
                Edit
              </Button>
            </div>
          }
          border
        >
          <div>{storageFacility.facilityName}</div>
          <div>{storageFacility.phone}</div>
          <div>{formatAddress(storageFacility.address)}</div>
          <div>Lot {storageFacility.lotNumber}</div>
        </SimpleSection>
      )}

      {serviceOrderNumber && (
        <SimpleSection
          key="service-order-number"
          header={
            <div className={styles.ShipmentDetailsSidebar}>
              Service order number
              <Button
                size="small"
                type="button"
                onClick={handleShowSonModal}
                className="float-right usa-link modal-link"
                data-testid="edit-facility-info-modal-open"
                unstyled
              >
                Edit
              </Button>
            </div>
          }
          border
        >
          <div>{serviceOrderNumber}</div>
        </SimpleSection>
      )}

      {(tacType || sacType) && (
        <SimpleSection
          key="accounting-codes"
          header={
            <>
              Accounting codes
              <Link to="" className="usa-link float-right">
                Edit
              </Link>
            </>
          }
          border
        >
          {tacType && tac && <div>TAC: {formatAccountingCode(tac, tacType)}</div>}
          {sacType && sac && <div>SAC: {formatAccountingCode(sac, sacType)}</div>}
        </SimpleSection>
      )}

      {(secondaryAddresses?.secondaryPickupAddress || secondaryAddresses?.secondaryDeliveryAddress) && (
        <SimpleSection header="Secondary addresses" border>
          {secondaryAddresses?.secondaryPickupAddress && (
            <SimpleSection header="Pickup">
              <div>{formatAddress(secondaryAddresses?.secondaryPickupAddress)}</div>
            </SimpleSection>
          )}

          {secondaryAddresses?.secondaryDeliveryAddress && (
            <SimpleSection header="Destination">
              <div>{formatAddress(secondaryAddresses?.secondaryDeliveryAddress)}</div>
            </SimpleSection>
          )}
        </SimpleSection>
      )}
    </div>
  );
};

ShipmentDetailsSidebar.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape,
  ordersLOA: OrdersLOAShape,
  handleEditFacilityInfo: PropTypes.func.isRequired,
  handleEditServiceOrderNumber: PropTypes.func,
};

ShipmentDetailsSidebar.defaultProps = {
  className: '',
  shipment: {},
  ordersLOA: {
    tac: '',
    sac: '',
    ntsTac: '',
    ntsSac: '',
  },
  handleEditServiceOrderNumber: () => {},
};

export default ShipmentDetailsSidebar;
