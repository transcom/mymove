import React from 'react';
import { PropTypes } from 'prop-types';

import ExpandableServiceItemRow from '../ExpandableServiceItemRow/ExpandableServiceItemRow';
import ShipmentModificationTag from '../../ShipmentModificationTag/ShipmentModificationTag';

import styles from './PaymentRequestDetails.module.scss';

import { LOA_TYPE, PAYMENT_REQUEST_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { PaymentServiceItemShape } from 'types';
import { MTOServiceItemShape } from 'types/order';
import { formatDateFromIso } from 'utils/formatters';
import PAYMENT_REQUEST_STATUSES from 'constants/paymentRequestStatus';
import { shipmentModificationTypes } from 'constants/shipments';
import { AccountingCodesShape } from 'types/accountingCodes';
import { formatAccountingCode } from 'utils/shipmentDisplay';

const shipmentHeadingAndStyle = (mtoShipmentType) => {
  switch (mtoShipmentType) {
    case undefined:
    case null:
      return ['Basic service items', styles.basicServiceType];
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      return ['HHG', styles.hhgShipmentType];
    case SHIPMENT_OPTIONS.NTS:
      return ['Non-temp storage', styles.ntsrShipmentType];
    case SHIPMENT_OPTIONS.NTSR:
      return ['Non-temp storage release', styles.ntsrShipmentType];
    default:
      return [mtoShipmentType, styles.basicServiceType];
  }
};

const PaymentRequestAccountingCodes = ({
  tacs,
  sacs,
  tacType,
  sacType,
  shipmentType,
  mtoShipmentID,
  showEdit,
  onEditClick,
}) => {
  const handleEditClick = () => {
    onEditClick({
      mtoShipmentID,
      tacType,
      sacType,
      shipmentType,
    });
  };

  const tacValue = tacType && tacs[tacType] ? formatAccountingCode(tacs[tacType], tacType) : '—';
  const sacValue = sacType && sacs[sacType] ? formatAccountingCode(sacs[sacType], sacType) : '—';

  return (
    <span style={{ display: 'block' }}>
      <strong>TAC: </strong>
      <span data-testid="tac">{tacValue}</span>
      &nbsp;|&nbsp;
      <strong>SAC: </strong>
      <span data-testid="sac">{sacValue}</span>
      {showEdit && (
        <button type="button" className={styles.EditButton} onClick={handleEditClick}>
          Edit
        </button>
      )}
    </span>
  );
};

PaymentRequestAccountingCodes.propTypes = {
  tacs: AccountingCodesShape,
  sacs: AccountingCodesShape,
  tacType: PropTypes.string,
  sacType: PropTypes.string,
  shipmentType: PropTypes.string,
  mtoShipmentID: PropTypes.string,
  showEdit: PropTypes.bool,
  onEditClick: PropTypes.func,
};

PaymentRequestAccountingCodes.defaultProps = {
  tacs: {},
  sacs: {},
  tacType: null,
  sacType: null,
  shipmentType: null,
  mtoShipmentID: null,
  showEdit: false,
  onEditClick: () => {},
};

const PaymentRequestDetails = ({ serviceItems, shipment, paymentRequestStatus, tacs, sacs, onEditClick }) => {
  const mtoShipmentType = serviceItems?.[0]?.mtoShipmentType;
  const [headingType, shipmentStyle] = shipmentHeadingAndStyle(mtoShipmentType);
  const { modificationType, departureDate, address, mtoServiceItems } = shipment;

  const findAdditionalServiceItemData = (mtoServiceItemCode) => {
    return mtoServiceItems?.find((mtoServiceItem) => mtoServiceItem.reServiceCode === mtoServiceItemCode);
  };

  return (
    serviceItems.length > 0 && (
      <div className={styles.PaymentRequestDetails}>
        <div className="stackedtable-header">
          <div className={styles.shipmentType}>
            <div className={shipmentStyle} />
            <h3>
              {headingType} ({serviceItems.length} {serviceItems.length > 1 ? 'items' : 'item'})
              {modificationType && <ShipmentModificationTag shipmentModificationType={modificationType} />}
            </h3>
          </div>
          <div>
            <p>
              <small>
                {departureDate && (
                  <strong data-testid="departure-date">
                    Departed {formatDateFromIso(departureDate, 'DD MMM YYYY')}
                  </strong>
                )}{' '}
                {address && <span data-testid="pickup-to-destination">{address}</span>}
                {mtoShipmentType && (
                  <PaymentRequestAccountingCodes
                    tacs={tacs}
                    sacs={sacs}
                    tacType={shipment.tacType}
                    sacType={shipment.sacType}
                    shipmentType={mtoShipmentType}
                    mtoShipmentID={shipment.mtoShipmentID}
                    showEdit={headingType !== 'HHG'}
                    onEditClick={onEditClick}
                  />
                )}
              </small>
            </p>
          </div>
        </div>
        <table className="table--stacked">
          <colgroup>
            <col style={{ width: '50%' }} />
            <col style={{ width: '25%' }} />
            <col style={{ width: '25%' }} />
          </colgroup>
          <thead>
            <tr>
              <th>Service item</th>
              <th className="align-right">Amount</th>
              <th className="align-right">Status</th>
            </tr>
          </thead>
          <tbody>
            {serviceItems.map((item, index) => {
              return (
                <ExpandableServiceItemRow
                  serviceItem={item}
                  additionalServiceItemData={findAdditionalServiceItemData(item.mtoServiceItemCode)}
                  key={item.id}
                  index={index}
                  disableExpansion={paymentRequestStatus === PAYMENT_REQUEST_STATUSES.PENDING}
                  paymentIsDeprecated={paymentRequestStatus === PAYMENT_REQUEST_STATUS.DEPRECATED}
                />
              );
            })}
          </tbody>
        </table>
      </div>
    )
  );
};

PaymentRequestDetails.propTypes = {
  serviceItems: PropTypes.arrayOf(PaymentServiceItemShape).isRequired,
  shipment: PropTypes.shape({
    address: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
    modificationType: PropTypes.oneOfType([
      PropTypes.string,
      PropTypes.oneOf(Object.values(shipmentModificationTypes)),
    ]),
    departureDate: PropTypes.string,
    mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
    tacType: PropTypes.oneOf(Object.values(LOA_TYPE)),
    sacType: PropTypes.oneOf(Object.values(LOA_TYPE)),
    mtoShipmentID: PropTypes.string,
  }),
  paymentRequestStatus: PropTypes.oneOf(Object.values(PAYMENT_REQUEST_STATUSES)).isRequired,
  tacs: AccountingCodesShape,
  sacs: AccountingCodesShape,
  onEditClick: PropTypes.func,
};

PaymentRequestDetails.defaultProps = {
  shipment: {
    departureDate: '',
    address: '',
    modificationType: '',
    mtoServiceItems: [],
  },
  tacs: {},
  sacs: {},
  onEditClick: () => {},
};

export default PaymentRequestDetails;
