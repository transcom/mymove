import React from 'react';
import PropTypes from 'prop-types';

import { MTOServiceItemDimensionShape, MTOServiceItemCustomerContactShape } from '../../../types/moveOrder';
import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import styles from './RequestedServiceItemsTable.module.scss';

import ServiceItemTableHasImg from 'components/ServiceItemTableHasImg/index';

const RequestedServiceItemsTable = ({
  serviceItems,
  handleUpdateMTOServiceItemStatus,
  handleShowRejectionDialog,
  statusForTableType,
}) => {
  const chooseTitleText = (status) => {
    switch (status) {
      case SERVICE_ITEM_STATUS.SUBMITTED:
        return 'Requested';
      case SERVICE_ITEM_STATUS.APPROVED:
        return 'Approved';
      case SERVICE_ITEM_STATUS.REJECTED:
        return 'Rejected';
      default:
        return 'Requested';
    }
  };

  return (
    <div className={styles.RequestedServiceItemsTable}>
      <h4>
        {chooseTitleText(statusForTableType)} service items&nbsp;
        <span>
          ({serviceItems.length} {serviceItems.length === 1 ? 'item' : 'items'})
        </span>
      </h4>
      <ServiceItemTableHasImg
        serviceItems={serviceItems}
        handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
        handleShowRejectionDialog={handleShowRejectionDialog}
        statusForTableType={statusForTableType}
      />
    </div>
  );
};

RequestedServiceItemsTable.propTypes = {
  handleUpdateMTOServiceItemStatus: PropTypes.func.isRequired,
  handleShowRejectionDialog: PropTypes.func.isRequired,
  statusForTableType: PropTypes.string.isRequired,
  serviceItems: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string,
      createdAt: PropTypes.string,
      approvedAt: PropTypes.string,
      rejectedAt: PropTypes.string,
      serviceItem: PropTypes.string,
      code: PropTypes.string,
      details: PropTypes.shape({
        reason: PropTypes.string,
        pickupPostalCode: PropTypes.string,
        imgURL: PropTypes.string,
        itemDimensions: MTOServiceItemDimensionShape,
        crateDimensions: MTOServiceItemDimensionShape,
        firstCustContact: MTOServiceItemCustomerContactShape,
        secondCustContact: MTOServiceItemCustomerContactShape,
      }),
    }),
  ).isRequired,
};

export default RequestedServiceItemsTable;
