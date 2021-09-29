import React from 'react';
import PropTypes, { shape } from 'prop-types';
import classnames from 'classnames';
import moment from 'moment';
import { withRouter } from 'react-router-dom';

import styles from './PrimeSimulatorListMoveCard.module.scss';

import { formatDateFromIso } from 'shared/formatters';

const PrimeSimulatorListMoveCard = ({ listMove }) => {
  return (
    <div className={classnames(styles.PrimeSimulatorListMoveCard, 'container')}>
      <div className={styles.summary}>
        <div className={styles.header}>
          <h2>Move {listMove.id}</h2>
        </div>
        <div className={styles.footer}>
          <dl>
            <dt>Move Code</dt>
            <dd>{listMove.moveCode}</dd>
            <dt>Order ID</dt>
            <dd>{listMove.orderID}</dd>
            <dt>Reference ID</dt>
            <dd>{listMove.referenceId}</dd>
            <dt>PPM type</dt>
            <dd>{listMove.ppmType}</dd>
            <dt>PPM estimate weight</dt>
            <dd>{listMove.ppmEstimateWeight}</dd>
            <dt>Available To Prime At</dt>
            <dd>
              <span className={styles.dateAt}>
                {moment(listMove.availableToPrimeAt).fromNow()} on{' '}
                {formatDateFromIso(listMove.availableToPrimeAt, 'DD MMM YYYY')}
              </span>
            </dd>
          </dl>
        </div>
      </div>
    </div>
  );
};

PrimeSimulatorListMoveCard.propTypes = {
  listMove: shape({
    id: PropTypes.string,
    moveCode: PropTypes.string,
    createdAt: PropTypes.string,
    orderID: PropTypes.string,
    referenceId: PropTypes.string,
    availableToPrimeAt: PropTypes.string,
    updatedAt: PropTypes.string,
    ppmType: PropTypes.string,
    ppmEstimateWeight: PropTypes.number,
  }),
};

PrimeSimulatorListMoveCard.defaultProps = {
  listMove: {},
};

export default withRouter(PrimeSimulatorListMoveCard);
