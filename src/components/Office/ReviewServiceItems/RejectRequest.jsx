import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './RejectRequest.module.scss';

/**
 * This component represents a section shown in the ReviewDetailsCard at the end of navigation.
 * Only shows if all service items have been rejected.
 * */
const RejectRequest = ({ onClick }) => {
  return (
    <div data-testid="RejectRequest" className={styles.RejectRequest}>
      <p data-testid="content" className={styles.content}>
        You&apos;re rejecting all service items. No payment will be authorized.
      </p>
      <Button data-testid="rejectRequestBtn" type="button" onClick={onClick}>
        Reject request
      </Button>
    </div>
  );
};

RejectRequest.propTypes = {
  onClick: PropTypes.func,
};

RejectRequest.defaultProps = {
  onClick: null,
};

export default RejectRequest;
