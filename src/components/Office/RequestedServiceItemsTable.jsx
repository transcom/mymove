import React from 'react';
import PropTypes from 'prop-types';

import ServiceItemTableHasImg from '../ServiceItemTableHasImg';

const RequestedServiceItemsTable = ({ serviceItems }) => {
  return (
    <>
      <h3>Requested service items ({serviceItems.length})</h3>
      <ServiceItemTableHasImg serviceItems={serviceItems} />
    </>
  );
};

RequestedServiceItemsTable.propTypes = {
  serviceItems: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string,
      dateRequested: PropTypes.string,
      serviceItem: PropTypes.string,
      code: PropTypes.string,
      details: PropTypes.object,
    }),
  ).isRequired,
};

export default RequestedServiceItemsTable;
