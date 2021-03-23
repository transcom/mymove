import React, { useState } from 'react';
import { PropTypes } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './PaymentRequestDetails.module.scss';

import { PAYMENT_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { allowedServiceItemCalculations } from 'constants/serviceItems';
import { formatCents, toDollarString } from 'shared/formatters';
import { PaymentServiceItemShape } from 'types';
import ServiceItemCalculations from 'components/Office/ServiceItemCalculations/ServiceItemCalculations';

const shipmentHeadingAndStyle = (mtoShipmentType) => {
  switch (mtoShipmentType) {
    case undefined:
    case null:
      return ['Basic service items', styles.basicServiceType];
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      return ['Household goods', styles.hhgShipmentType];
    case SHIPMENT_OPTIONS.NTS:
      return ['Non-temp storage', styles.ntsrShipmentType];
    case SHIPMENT_OPTIONS.NTSR:
      return ['Non-temp storage release', styles.ntsrShipmentType];
    default:
      return [mtoShipmentType, styles.basicServiceType];
  }
};

const PaymentRequestDetails = ({ serviceItems, shipmentAddress }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const mtoShipmentType = serviceItems?.[0]?.mtoShipmentType;
  const [headingType, shipmentStyle] = shipmentHeadingAndStyle(mtoShipmentType);
  const expandableIconClasses = classnames({
    'chevron-down': isExpanded,
    'chevron-right': !isExpanded,
  });

  const canClickToExpandContent = (canShowExpandableContent, item) => {
    if (canShowExpandableContent && item.status !== PAYMENT_SERVICE_ITEM_STATUS.REQUESTED) {
      return true;
    }
    return false;
  };

  return (
    serviceItems.length > 0 && (
      <div className={styles.PaymentRequestDetails}>
        <div className="stackedtable-header">
          <div className={styles.shipmentType}>
            <div className={shipmentStyle} />
            <h3>
              {headingType} ({serviceItems.length} {serviceItems.length > 1 ? 'items' : 'item'})
            </h3>
          </div>
          {shipmentAddress !== '' && <p data-testid="pickup-to-destination">{shipmentAddress}</p>}
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
              // TODO - temporary, will remove once all service item calculations are implemented
              const canShowExpandableContent = allowedServiceItemCalculations.includes(item.mtoServiceItemCode);

              const toggleExpandableRow = () => {
                setIsExpanded((prev) => !prev);
              };

              const tableRowClasses = classnames(styles.expandable, {
                [styles.expandedRow]: isExpanded,
              });
              const tableDetailClasses = classnames(styles.expandable, {
                [styles.expandedDetail]: isExpanded,
              });

              return (
                <React.Fragment key={item.id}>
                  <tr
                    data-groupid={index}
                    className={tableRowClasses}
                    onClick={canClickToExpandContent(canShowExpandableContent, item) ? toggleExpandableRow : undefined}
                  >
                    <td data-testid="serviceItemName">
                      {canShowExpandableContent && item.status !== PAYMENT_SERVICE_ITEM_STATUS.REQUESTED && (
                        <FontAwesomeIcon className={styles.icon} icon={expandableIconClasses} />
                      )}
                      {item.mtoServiceItemName}
                    </td>
                    <td data-testid="serviceItemAmount">{toDollarString(formatCents(item.priceCents))}</td>
                    <td data-testid="serviceItemStatus">
                      {item.status === PAYMENT_SERVICE_ITEM_STATUS.REQUESTED && (
                        <div className={styles.needsReview}>
                          <FontAwesomeIcon icon="exclamation-circle" />
                          <span>Needs review</span>
                        </div>
                      )}
                      {item.status === PAYMENT_SERVICE_ITEM_STATUS.APPROVED && (
                        <div className={styles.accepted}>
                          <FontAwesomeIcon icon="check" />
                          <span>Accepted</span>
                        </div>
                      )}
                      {item.status === PAYMENT_SERVICE_ITEM_STATUS.DENIED && (
                        <div className={styles.rejected}>
                          <FontAwesomeIcon icon="times" />
                          <span>Rejected</span>
                        </div>
                      )}
                    </td>
                  </tr>
                  {canShowExpandableContent && isExpanded && (
                    <tr data-testid="serviceItemCaclulations" data-groupdid={index} className={tableDetailClasses}>
                      <td colSpan={3}>
                        <ServiceItemCalculations
                          itemCode={item.mtoServiceItemCode}
                          totalAmountRequested={item.priceCents}
                          serviceItemParams={item.paymentServiceItemParams}
                        />
                      </td>
                    </tr>
                  )}
                </React.Fragment>
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
  shipmentAddress: PropTypes.string,
};

PaymentRequestDetails.defaultProps = {
  shipmentAddress: '',
};

export default PaymentRequestDetails;
