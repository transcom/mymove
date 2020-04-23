import React from 'react';
import ServiceItemTableHasImg from '../ServiceItemTableHasImg';
// import PropTypes from 'prop-types';
// import classNames from 'classnames/bind';
// import styles from './index.module.scss';
//
// const cx = classNames.bind(styles);

const RequestedServiceItemsTable = () => (
  <ServiceItemTableHasImg
    serviceItems={[
      {
        id: 'abc-123',
        dateRequested: '20 Nov 2020',
        serviceItem: 'Dom. Origin 1st Day SIT',
        code: 'DOMSIT',
        details: {
          text: "60612 and here's the reason",
          imgURL: null,
        },
      },
    ]}
  />
);

export default RequestedServiceItemsTable;
