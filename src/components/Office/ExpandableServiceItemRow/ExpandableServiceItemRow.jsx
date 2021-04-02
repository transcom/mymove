import React, { useState } from 'react';
import { PropTypes } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './ExpandableServiceItemRow.module.scss';

import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { allowedServiceItemCalculations } from 'constants/serviceItems';
import { PaymentServiceItemShape } from 'types';
import { formatCents, toDollarString } from 'shared/formatters';
import ServiceItemCalculations from 'components/Office/ServiceItemCalculations/ServiceItemCalculations';

const ExpandableServiceItemRow = ({ serviceItem, index }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const canClickToExpandContent = (canShowExpandableContent, item) => {
    return canShowExpandableContent && item.status !== PAYMENT_SERVICE_ITEM_STATUS.REQUESTED;
  };
  const canShowExpandableContent = allowedServiceItemCalculations.includes(serviceItem.mtoServiceItemCode);

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

  return (
    <>
      <tr
        data-groupid={index}
        className={tableRowClasses}
        onClick={canClickToExpandContent(canShowExpandableContent, serviceItem) ? handleExpandClick : undefined}
        aria-expanded={isExpanded}
      >
        <td data-testid="serviceItemName">
          {canShowExpandableContent && serviceItem.status !== PAYMENT_SERVICE_ITEM_STATUS.REQUESTED && (
            <FontAwesomeIcon className={styles.icon} icon={expandableIconClasses} />
          )}
          {serviceItem.mtoServiceItemName}
        </td>
        <td data-testid="serviceItemAmount">{toDollarString(formatCents(serviceItem.priceCents))}</td>
        <td data-testid="serviceItemStatus">
          {serviceItem.status === PAYMENT_SERVICE_ITEM_STATUS.REQUESTED && (
            <div className={styles.needsReview}>
              <FontAwesomeIcon icon="exclamation-circle" />
              <span>Needs review</span>
            </div>
          )}
          {serviceItem.status === PAYMENT_SERVICE_ITEM_STATUS.APPROVED && (
            <div className={styles.accepted}>
              <FontAwesomeIcon icon="check" />
              <span>Accepted</span>
            </div>
          )}
          {serviceItem.status === PAYMENT_SERVICE_ITEM_STATUS.DENIED && (
            <div className={styles.rejected}>
              <FontAwesomeIcon icon="times" />
              <span>Rejected</span>
            </div>
          )}
        </td>
      </tr>
      {isExpanded && (
        <tr data-testid="serviceItemCaclulations" data-groupdid={index} className={tableDetailClasses}>
          <td colSpan={3}>
            <ServiceItemCalculations
              itemCode={serviceItem.mtoServiceItemCode}
              totalAmountRequested={serviceItem.priceCents}
              serviceItemParams={serviceItem.paymentServiceItemParams}
            />
          </td>
        </tr>
      )}
    </>
  );
};

ExpandableServiceItemRow.propTypes = {
  serviceItem: PaymentServiceItemShape.isRequired,
  index: PropTypes.number.isRequired,
};

export default ExpandableServiceItemRow;
