import React from 'react';
import propTypes from 'prop-types';
import { get } from 'lodash';

import classNames from 'classnames/bind';
import styles from 'components/CustomerRemarksAgentsDetails/customerRemarksAgentsDetails.module.scss';

const cx = classNames.bind(styles);

const CustomerRemarksAgentsDetails = ({ customerRemarks, releasingAgent, receivingAgent }) => (
  <div>
    <div className={cx('container--small')}>
      <div className={cx('container__heading')}>Customer remarks</div>
      <hr />
      {customerRemarks}
    </div>
    <div className={cx('container--small')}>
      <div className={cx('container__heading')}>Releasing agent</div>
      <hr />
      {(get(releasingAgent, 'firstName') || get(releasingAgent, 'lastName')) && (
        <>
          {`${get(releasingAgent, 'firstName')} ${get(releasingAgent, 'lastName')}`}
          <br />
        </>
      )}

      {get(releasingAgent, 'phone') && (
        <>
          {get(releasingAgent, 'phone')}
          <br />
        </>
      )}
      {get(releasingAgent, 'email')}
    </div>
    <div className={cx('container--small')}>
      <div className={cx('container__heading')}>Receiving agent</div>
      <hr />
      {(get(receivingAgent, 'firstName') || get(receivingAgent, 'lastName')) && (
        <>
          {`${get(receivingAgent, 'firstName')} ${get(receivingAgent, 'lastName')}`}
          <br />
        </>
      )}

      {get(receivingAgent, 'phone') && (
        <>
          {get(receivingAgent, 'phone')}
          <br />
        </>
      )}
      {get(receivingAgent, 'email')}
    </div>
  </div>
);

CustomerRemarksAgentsDetails.propTypes = {
  customerRemarks: propTypes.string,
  releasingAgent: propTypes.shape({
    firstName: propTypes.string,
    lastName: propTypes.string,
    phone: propTypes.string,
    email: propTypes.string,
  }),
  receivingAgent: propTypes.shape({
    firstName: propTypes.string,
    lastName: propTypes.string,
    phone: propTypes.string,
    email: propTypes.string,
  }),
};

export default CustomerRemarksAgentsDetails;
