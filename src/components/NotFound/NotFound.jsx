import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './NotFound.module.scss';

const NotFound = ({ handleOnClick }) => {
  return (
    <div className={classnames('usa-grid', styles.notFound)}>
      <div className="grid-container">
        <b>Error - 404</b>
        <div>
          <b>Let&apos;s move you in the right direction</b>
        </div>
        <h1>
          <b>We can&apos;t find the page you&apos;re looking for</b>
        </h1>
        <div className={styles.explanation}>
          <p>You are seeing this because the page you&apos;re looking for doesn&apos;t exist or has been removed.</p>
          <p>
            We suggest checking the spelling in the URL or return{' '}
            <Button unstyled className={styles.goBack} onClick={handleOnClick}>
              back home.
            </Button>
          </p>
        </div>
      </div>
    </div>
  );
};

NotFound.propTypes = {
  handleOnClick: PropTypes.func.isRequired,
};

export default NotFound;
