import React from 'react';
import * as PropTypes from 'prop-types';
import { generatePath, useParams, useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './RequestedShipments.module.scss';

import { SERVICE_ITEM_CODES } from 'constants/serviceItems';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { tooRoutes } from 'constants/routes';
import { ADDRESS_UPDATE_STATUS, shipmentDestinationTypes } from 'constants/shipments';
import { shipmentTypeLabels } from 'content/shipments';
import shipmentCardsStyles from 'styles/shipmentCards.module.scss';
import { MTOServiceItemShape, OrdersInfoShape } from 'types/order';
import { ShipmentShape } from 'types/shipment';
import { formatDateFromIso } from 'utils/formatters';
import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';
import { SHIPMENT_OPTIONS_URL, FEATURE_FLAG_KEYS } from 'shared/constants';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

// nts defaults show preferred pickup date and pickup address, flagged items when collapsed
// ntsr defaults shows preferred delivery date, storage facility address, destination address, flagged items when collapsed
// Different things show when collapsed depending on if the shipment is an external vendor or not.
const showWhenCollapsedWithExternalVendor = {
  HHG_INTO_NTS_DOMESTIC: ['serviceOrderNumber', 'requestedDeliveryDate'],
  HHG_OUTOF_NTS_DOMESTIC: ['serviceOrderNumber', 'requestedPickupDate'],
};

const showWhenCollapsedWithGHCPrime = {
  HHG_INTO_NTS_DOMESTIC: ['tacType', 'requestedDeliveryDate'],
  HHG_OUTOF_NTS_DOMESTIC: ['ntsRecordedWeight', 'serviceOrderNumber', 'tacType', 'requestedPickupDate'],
};

const errorIfMissing = [
  {
    fieldName: 'destinationAddress',
    condition: (shipment) => shipment.deliveryAddressUpdate?.status === ADDRESS_UPDATE_STATUS.REQUESTED,
    optional: true,
  },
];

const ApprovedRequestedShipments = ({
  mtoShipments,
  closeoutOffice,
  ordersInfo,
  mtoServiceItems,
  displayDestinationType,
  isMoveLocked,
}) => {
  const ordersLOA = {
    tac: ordersInfo.tacMDC,
    sac: ordersInfo.sacSDN,
    ntsTac: ordersInfo.NTStac,
    ntsSac: ordersInfo.NTSsac,
  };

  const shipmentDisplayInfo = (shipment, dutyLocationPostal) => {
    const destType = displayDestinationType ? shipmentDestinationTypes[shipment.destinationType] : null;

    return {
      ...shipment,
      heading: shipmentTypeLabels[shipment.shipmentType],
      isDiversion: shipment.diversion,
      shipmentStatus: shipment.status,
      destinationAddress: shipment.destinationAddress || dutyLocationPostal,
      destinationType: destType,
      displayDestinationType,
      closeoutOffice,
    };
  };

  const { moveCode } = useParams();
  const navigate = useNavigate();
  const handleButtonDropdownChange = (e) => {
    const selectedOption = e.target.value;

    const addShipmentPath = `${generatePath(tooRoutes.SHIPMENT_ADD_PATH, {
      moveCode,
      shipmentType: selectedOption,
    })}`;

    navigate(addShipmentPath);
  };

  const dutyLocationPostal = { postalCode: ordersInfo.newDutyLocation?.address?.postalCode };

  const [enableBoat, setEnableBoat] = React.useState(false);
  const [enableMobileHome, setEnableMobileHome] = React.useState(false);
  React.useEffect(() => {
    const fetchData = async () => {
      setEnableBoat(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BOAT));
      setEnableMobileHome(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.MOBILE_HOME));
    };
    fetchData();
  }, []);

  const allowedShipmentOptions = () => {
    return (
      <>
        <option data-testid="hhgOption" value={SHIPMENT_OPTIONS_URL.HHG}>
          HHG
        </option>
        <option value={SHIPMENT_OPTIONS_URL.PPM}>PPM</option>
        <option value={SHIPMENT_OPTIONS_URL.NTS}>NTS</option>
        <option value={SHIPMENT_OPTIONS_URL.NTSrelease}>NTS-release</option>
        {enableBoat && <option value={SHIPMENT_OPTIONS_URL.BOAT}>Boat</option>}
        {enableMobileHome && <option value={SHIPMENT_OPTIONS_URL.MOBILE_HOME}>Mobile Home</option>}
      </>
    );
  };

  return (
    <div className={styles.RequestedShipments} data-testid="requested-shipments">
      <h2>Approved Shipments</h2>
      <div className={styles.dropdownButton}>
        {!isMoveLocked && (
          <Restricted to={permissionTypes.createTxoShipment}>
            <ButtonDropdown
              ariaLabel="Add a new shipment"
              data-testid="addShipmentButton"
              onChange={handleButtonDropdownChange}
            >
              <option value="" label="Add a new shipment">
                Add a new shipment
              </option>
              {allowedShipmentOptions()}
            </ButtonDropdown>
          </Restricted>
        )}
      </div>

      <div className={shipmentCardsStyles.shipmentCards}>
        {mtoShipments &&
          mtoShipments.map((shipment) => {
            const editUrl = `../${generatePath(tooRoutes.SHIPMENT_EDIT_PATH, {
              shipmentId: shipment.id,
            })}`;

            return (
              <ShipmentDisplay
                key={shipment.id}
                shipmentId={shipment.id}
                shipmentType={shipment.shipmentType}
                displayInfo={shipmentDisplayInfo(shipment, dutyLocationPostal)}
                ordersLOA={ordersLOA}
                showWhenCollapsed={
                  shipment.usesExternalVendor
                    ? showWhenCollapsedWithExternalVendor[shipment.shipmentType]
                    : showWhenCollapsedWithGHCPrime[shipment.shipmentType]
                }
                errorIfMissing={errorIfMissing}
                isSubmitted={false}
                editURL={editUrl}
                isMoveLocked={isMoveLocked}
              />
            );
          })}
      </div>

      <div className={styles.serviceItems}>
        <h3>Service Items</h3>

        <table className="table--stacked">
          <colgroup>
            <col style={{ width: '75%' }} />
            <col style={{ width: '25%' }} />
          </colgroup>
          <tbody>
            {mtoServiceItems &&
              mtoServiceItems
                .filter(
                  (serviceItem) =>
                    serviceItem.reServiceCode === SERVICE_ITEM_CODES.MS ||
                    serviceItem.reServiceCode === SERVICE_ITEM_CODES.CS,
                )
                .map((serviceItem) => (
                  <tr key={serviceItem.id}>
                    <td data-testid="basicServiceItemName">{serviceItem.reServiceName}</td>
                    <td data-testid="basicServiceItemDate">
                      {serviceItem.status === 'APPROVED' && (
                        <span>
                          <FontAwesomeIcon icon="check" className={styles.serviceItemApproval} />{' '}
                          {formatDateFromIso(serviceItem.approvedAt, 'DD MMM YYYY')}
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

ApprovedRequestedShipments.propTypes = {
  mtoShipments: PropTypes.arrayOf(ShipmentShape).isRequired,
  ordersInfo: OrdersInfoShape.isRequired,
  mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
  displayDestinationType: PropTypes.bool,
};

ApprovedRequestedShipments.defaultProps = {
  mtoServiceItems: [],
  displayDestinationType: false,
};

export default ApprovedRequestedShipments;
