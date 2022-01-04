import React from 'react';
import { Link } from 'react-router-dom';
import * as PropTypes from 'prop-types';

import SimpleSection from 'containers/SimpleSection/SimpleSection';
import { retrieveSAC, retrieveTAC, formatAgent, formatAddress, formatAccountingCode } from 'utils/shipmentDisplay';
import { ShipmentShape } from 'types/shipment';

const ShipmentDetailsSidebar = ({ className, shipment, ordersLOA }) => {
  const { agents, secondaryAddresses, serviceOrderNumber, storageFacility, sacType, tacType } = shipment;
  const tac = retrieveTAC(shipment.tacType, ordersLOA);
  const sac = retrieveSAC(shipment.sacType, ordersLOA);

  return (
    <div className={className}>
      {agents &&
        agents.map((agent) => (
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
            <>
              Facility info and address
              <Link to="" className="usa-link float-right">
                Edit
              </Link>
            </>
          }
          border
        >
          <div>{formatAddress(storageFacility)}</div>
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
          {tacType && tac && <div>{formatAccountingCode(tac, tacType)}</div>}
          {sacType && sac && <div>{formatAccountingCode(sac, sacType)}</div>}
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
  ordersLOA: PropTypes.shape({
    tac: PropTypes.string,
    sac: PropTypes.string,
    ntsTAC: PropTypes.string,
    ntsSAC: PropTypes.string,
  }),
};

ShipmentDetailsSidebar.defaultProps = {
  className: '',
  shipment: {},
  ordersLOA: {
    tac: '',
    sac: '',
    ntsTAC: '',
    ntsSAC: '',
  },
};

export default ShipmentDetailsSidebar;
