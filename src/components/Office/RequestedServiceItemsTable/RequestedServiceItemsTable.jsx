import React from 'react';
import PropTypes from 'prop-types';

import { MTOServiceItemDimensionShape, MTOServiceItemCustomerContactShape } from '../../../types/moveOrder';

import styles from './RequestedServiceItemsTable.module.scss';

import ServiceItemTableHasImg from 'components/ServiceItemTableHasImg';

const RequestedServiceItemsTable = ({ serviceItems, handleUpdateMTOServiceItemStatus }) => {
  return (
    <div className={styles.RequestedServiceItemsTable}>
      <h4>
        Requested service items&nbsp;
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
      details: PropTypes.shape({
        text: PropTypes.oneOfType([PropTypes.object, PropTypes.string]),
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
