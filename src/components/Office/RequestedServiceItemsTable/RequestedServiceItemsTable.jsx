import React from 'react';
import PropTypes from 'prop-types';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';
import { ServiceItemDetailsShape } from '../../../types/serviceItems';

import styles from './RequestedServiceItemsTable.module.scss';

import ServiceItemsTable from 'components/Office/ServiceItemsTable/ServiceItemsTable';

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

  const statusTitleText = chooseTitleText(statusForTableType);

  return (
    <div className={styles.RequestedServiceItemsTable} data-testid={`${statusTitleText}ServiceItemsTable`}>
      <h3>
        {statusTitleText} service items&nbsp;
        <span>
          ({serviceItems.length} {serviceItems.length === 1 ? 'item' : 'items'})
        </span>
      </h3>
      <ServiceItemsTable
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
  serviceItems: PropTypes.arrayOf(ServiceItemDetailsShape).isRequired,
};

export default RequestedServiceItemsTable;
