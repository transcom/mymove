import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { generatePath } from 'react-router';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import ConnectedAccountingCodesModal from '../AccountingCodesModal/AccountingCodesModal';

import styles from 'components/Office/ShipmentDetails/ShipmentDetailsSidebar.module.scss';
import SimpleSection from 'containers/SimpleSection/SimpleSection';
import ConnectedEditFacilityInfoModal from 'components/Office/EditFacilityInfoModal/EditFacilityInfoModal';
import { retrieveSAC, retrieveTAC, formatAgent, formatAddress, formatAccountingCode } from 'utils/shipmentDisplay';
import { ShipmentShape } from 'types/shipment';
import { OrdersLOAShape } from 'types/order';
import { tooRoutes } from 'constants/routes';

const ShipmentDetailsSidebar = ({
  history,
  className,
  shipment,
  ordersLOA,
  handleEditFacilityInfo,
  handleEditAccountingCodes,
}) => {
  const { mtoAgents, secondaryAddresses, serviceOrderNumber, storageFacility, sacType, tacType } = shipment;
  const tac = retrieveTAC(shipment.tacType, ordersLOA);
  const sac = retrieveSAC(shipment.sacType, ordersLOA);

  const moveCode = 'HGNTSR';
  const editOrdersPath = generatePath(tooRoutes.ORDERS_EDIT_PATH, { moveCode });

  const [isEditFacilityInfoModalVisible, setIsEditFacilityInfoModalVisible] = useState(false);
  const [isAccountingCodesModalVisible, setIsAccountingCodesModalVisible] = useState(false);

  const handleShowEditFacilityInfoModal = () => {
    setIsEditFacilityInfoModalVisible(true);
  };

  const handleShowAccountingCodesModal = () => {
    setIsAccountingCodesModalVisible(true);
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

      <ConnectedAccountingCodesModal
        isOpen={isAccountingCodesModalVisible}
        onSubmit={(accountingTypes) => {
          handleEditAccountingCodes(accountingTypes, shipment);
          setIsAccountingCodesModalVisible(false);
        }}
        onClose={() => {
          setIsAccountingCodesModalVisible(false);
        }}
        onEditCodesClick={() => history.push(editOrdersPath)}
        shipmentType={shipment.shipmentType}
        TACs={{
          HHG: ordersLOA.tac,
          NTS: ordersLOA.ntsTac,
        }}
        SACs={{
          HHG: ordersLOA.sac,
          NTS: ordersLOA.ntsSac,
        }}
        tacType={shipment.tacType}
        sacType={shipment.sacType}
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
              <Button
                size="small"
                type="button"
                onClick={handleShowAccountingCodesModal}
                className="float-right usa-link padding-right-0"
                data-testid="edit-accounting-code-modal-open"
                unstyled
              >
                Edit
              </Button>
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
  history: PropTypes.shape({
    push: PropTypes.func.isRequired,
  }),
  shipment: ShipmentShape,
  ordersLOA: OrdersLOAShape,
  handleEditFacilityInfo: PropTypes.func.isRequired,
  handleEditAccountingCodes: PropTypes.func.isRequired,
};

ShipmentDetailsSidebar.defaultProps = {
  className: '',
  history: {
    push: () => {},
  },
  shipment: {},
  ordersLOA: {
    tac: '',
    sac: '',
    ntsTac: '',
    ntsSac: '',
  },
};

export default ShipmentDetailsSidebar;
