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
          text: {
            ZIP: '60612',
            Reason: "here's the reason",
          },
          imgURL: null,
        },
      },
      {
        id: 'abc-1234',
        dateRequested: '22 Nov 2020',
        serviceItem: 'Dom. Destination 1st Day SIT',
        code: 'DDFSIT',
        details: {
          text: {
            'First available delivery date': '22 Nov 2020',
            'First customer contact': '22 Nov 2020 12:00pm',
            'Second customer contact': '22 Nov 2020 12:00pm',
          },
          imgURL: null,
        },
      },
    ]}
  />
);

export default RequestedServiceItemsTable;
