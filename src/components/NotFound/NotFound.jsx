import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { useNavigate } from 'react-router-dom';

import styles from './NotFound.module.scss';

const NotFound = ({ handleOnClick }) => {
  const navigate = useNavigate();
  return (
    <div className={classnames('usa-grid', styles.notFound)}>
      <div className="grid-container">
        <div className={styles.preheader}>
          <b>Error - 404</b>
          <div className={styles.preheaderQuip}>
            <b>Let&apos;s move you in the right direction</b>
          </div>
        </div>
        <h1>
          <b>We can&apos;t find the page you&apos;re looking for</b>
        </h1>
        <div className={styles.description}>
          <p className={styles.explanation}>
            You are seeing this because the page you&apos;re looking for doesn&apos;t exist or has been removed.
          </p>
          <p className={styles.recommendation}>
            We suggest checking the spelling in the URL or return{' '}
            <Button unstyled className={styles.goBack} onClick={handleOnClick || (() => navigate(-1))}>
              back home.
            </Button>
          </p>
        </div>
      </div>
    </div>
  );
};

NotFound.propTypes = {
  handleOnClick: PropTypes.func,
};

NotFound.defaultProps = {
  handleOnClick: null,
};

export default NotFound;
