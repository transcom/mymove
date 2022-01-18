import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from 'components/Office/ShipmentDetails/ShipmentDetailsSidebar.module.scss';
import SimpleSection from 'containers/SimpleSection/SimpleSection';
import { EditFacilityInfoModal } from 'components/Office/EditFacilityInfoModal/EditFacilityInfoModal';
import { retrieveSAC, retrieveTAC, formatAgent, formatAddress, formatAccountingCode } from 'utils/shipmentDisplay';
import { ShipmentShape } from 'types/shipment';
import { OrdersLOAShape } from 'types/order';

const ShipmentDetailsSidebar = ({ className, shipment, ordersLOA, handleEditFacilityInfo }) => {
  const { mtoAgents, secondaryAddresses, serviceOrderNumber, storageFacility, sacType, tacType } = shipment;
  const tac = retrieveTAC(shipment.tacType, ordersLOA);
  const sac = retrieveSAC(shipment.sacType, ordersLOA);

  const [isEditFacilityInfoModalVisible, setIsEditFacilityInfoModalVisible] = useState(false);

  const handleShowEditFacilityInfoModal = () => {
    setIsEditFacilityInfoModalVisible(true);
  };

  return (
    <div className={className}>
      {isEditFacilityInfoModalVisible && (
        <EditFacilityInfoModal
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
      )}

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
            <>
              Service order number
              <Link to="" className="usa-link float-right">
                Edit
              </Link>
            </>
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
};

export default ShipmentDetailsSidebar;
