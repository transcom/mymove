import React, { useState } from 'react';
import { PropTypes } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './ExpandableServiceItemRow.module.scss';

import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { allowedServiceItemCalculations, SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { PaymentServiceItemShape } from 'types';
import { MTOServiceItemShape } from 'types/order';
import { toDollarString, formatCents, formatDollarFromMillicents } from 'utils/formatters';
import ServiceItemCalculations from 'components/Office/ServiceItemCalculations/ServiceItemCalculations';

const ExpandableServiceItemRow = ({
  additionalServiceItemData,
  disableExpansion,
  index,
  paymentIsDeprecated,
  serviceItem,
  tppsDataExists,
}) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const canClickToExpandContent = (canShowExpandableContent, item) => {
    return canShowExpandableContent && (paymentIsDeprecated || item.status !== PAYMENT_SERVICE_ITEM_STATUS.REQUESTED);
  };
  const canShowExpandableContent =
    !disableExpansion &&
    (allowedServiceItemCalculations.includes(serviceItem.mtoServiceItemCode) || serviceItem.rejectionReason);

  const handleExpandClick = () => {
    setIsExpanded((prev) => !prev);
  };
  const expandableIconClasses = classnames({
    'chevron-down': isExpanded,
    'chevron-right': !isExpanded,
  });

  const tableRowClasses = classnames(styles.ExpandableServiceItemRow, styles.expandable, {
    [styles.expandedRow]: isExpanded,
    [styles.isExpandable]: canShowExpandableContent,
  });
  const tableDetailClasses = classnames(styles.ExpandableServiceItemRow, {
    [styles.expandedDetail]: isExpanded,
  });

  const colSpan =
    serviceItem.mtoServiceItemCode === SERVICE_ITEM_CODES.MS || serviceItem.mtoServiceItemCode === SERVICE_ITEM_CODES.CS
      ? 4
      : 2;

  return (
    <>
      <tr
        data-groupid={index}
        className={tableRowClasses}
        onClick={canClickToExpandContent(canShowExpandableContent, serviceItem) ? handleExpandClick : undefined}
        aria-expanded={isExpanded}
      >
        <td data-testid="serviceItemName">
          {canShowExpandableContent &&
            (paymentIsDeprecated || serviceItem.status !== PAYMENT_SERVICE_ITEM_STATUS.REQUESTED) && (
              <FontAwesomeIcon className={styles.icon} icon={expandableIconClasses} />
            )}
          {serviceItem.mtoServiceItemName}
          {additionalServiceItemData.standaloneCrate && ' - Standalone'}
        </td>
        <td data-testid="serviceItemAmount">{toDollarString(formatCents(serviceItem.priceCents))}</td>
        {tppsDataExists && (
          <td data-testid="serviceItemTPPSPaidAmount">
            {serviceItem.tppsInvoiceAmountPaidPerServiceItemMillicents > 0 &&
              toDollarString(formatDollarFromMillicents(serviceItem.tppsInvoiceAmountPaidPerServiceItemMillicents))}
          </td>
        )}
        <td data-testid="serviceItemStatus">
          {paymentIsDeprecated && (
            <div>
              <span data-testid="deprecated-marker">-</span>
            </div>
          )}
          {serviceItem.status === PAYMENT_SERVICE_ITEM_STATUS.REQUESTED && !paymentIsDeprecated && (
            <div className={styles.needsReview}>
              <FontAwesomeIcon icon="exclamation-circle" />
              <span>Needs review</span>
            </div>
          )}
          {serviceItem.status === PAYMENT_SERVICE_ITEM_STATUS.APPROVED && !paymentIsDeprecated && (
            <div className={styles.accepted}>
              <FontAwesomeIcon icon="check" />
              <span>Accepted</span>
            </div>
          )}
          {serviceItem.status === PAYMENT_SERVICE_ITEM_STATUS.DENIED && !paymentIsDeprecated && (
            <div className={styles.rejected}>
              <FontAwesomeIcon icon="times" />
              <span>Rejected</span>
            </div>
          )}
        </td>
      </tr>
      {isExpanded && (
        <tr data-testid="serviceItemCaclulations" data-groupdid={index} className={tableDetailClasses}>
          {Object.keys(additionalServiceItemData).length > 0 && (
            <td colSpan={1}>
              <ServiceItemCalculations
                itemCode={serviceItem.mtoServiceItemCode}
                totalAmountRequested={serviceItem.priceCents}
                serviceItemParams={serviceItem.paymentServiceItemParams}
                additionalServiceItemData={additionalServiceItemData}
                shipmentType={serviceItem.mtoShipmentType}
              />
            </td>
          )}
          {serviceItem.rejectionReason && (
            <td colSpan={colSpan} className={styles.rejectionReasonTd}>
              <div className={styles.rejectionReasonContainer}>
                <FontAwesomeIcon icon="times" />
                <h4 className={styles.title}>Rejection Reason</h4>
                <div className={styles.break} />
                <small className={styles.reasonText}>{serviceItem.rejectionReason}</small>
              </div>
            </td>
          )}
        </tr>
      )}
    </>
  );
};

ExpandableServiceItemRow.propTypes = {
  serviceItem: PaymentServiceItemShape.isRequired,
  index: PropTypes.number.isRequired,
  disableExpansion: PropTypes.bool,
  additionalServiceItemData: MTOServiceItemShape,
  paymentIsDeprecated: PropTypes.bool.isRequired,
};

ExpandableServiceItemRow.defaultProps = {
  disableExpansion: false,
  additionalServiceItemData: {},
};

export default ExpandableServiceItemRow;
