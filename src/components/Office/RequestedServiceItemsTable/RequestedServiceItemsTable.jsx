import React from 'react';
import PropTypes from 'prop-types';

import { SERVICE_ITEM_STATUS } from '../../../shared/constants';

import styles from './RequestedServiceItemsTable.module.scss';

import ServiceItemTableHasImg from 'components/ServiceItemTableHasImg';

const RequestedServiceItemsTable = ({ serviceItems, handleUpdateMTOServiceItemStatus }) => {
  const chooseTitleText = (serviceItem) => {
    switch (serviceItem.status) {
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
        {chooseTitleText(serviceItems[0])} service items&nbsp;
        <span>
          ({serviceItems.length} {serviceItems.length === 1 ? 'item' : 'items'})
        </span>
      </h4>
      <ServiceItemTableHasImg
        serviceItems={serviceItems}
        handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
      />
    </div>
  );
};

RequestedServiceItemsTable.propTypes = {
  handleUpdateMTOServiceItemStatus: PropTypes.func.isRequired,
  serviceItems: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string,
      submittedAt: PropTypes.string,
      serviceItem: PropTypes.string,
      code: PropTypes.string,
      details: PropTypes.object,
    }),
  ).isRequired,
};

export default RequestedServiceItemsTable;
