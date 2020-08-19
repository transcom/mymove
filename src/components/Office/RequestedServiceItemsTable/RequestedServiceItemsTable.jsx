import React from 'react';
import PropTypes from 'prop-types';

// eslint-disable-next-line no-unused-vars
import styles from './RequestedServiceItemsTable.module.scss';

import ServiceItemTableHasImg from 'components/ServiceItemTableHasImg';

const RequestedServiceItemsTable = ({ serviceItems }) => {
  return (
    <>
      <h4>
        Requested service items&nbsp;
        <span>
          ({serviceItems.length} {serviceItems.length === 1 ? 'item' : 'items'})
        </span>
      </h4>
      <ServiceItemTableHasImg serviceItems={serviceItems} />
    </>
  );
};

RequestedServiceItemsTable.propTypes = {
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
