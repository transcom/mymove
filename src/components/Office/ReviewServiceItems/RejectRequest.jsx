import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './RejectRequest.module.scss';

/**
 * This component represents a section shown in the ReviewDetailsCard at the end of navigation.
 * Only shows if all service items have been rejected.
 * */
const RejectRequest = ({ handleRejectBtn }) => {
  return (
    <div data-testid="RejectRequest" className={styles.RejectRequest}>
      <div>You&apos;re rejecting all service items. No payment will be authorized.</div>
      <Button type="button" onClick={handleRejectBtn}>
        Reject request
      </Button>
    </div>
  );
};

RejectRequest.propTypes = {
  handleRejectBtn: PropTypes.func,
};

RejectRequest.defaultProps = {
  handleRejectBtn: null,
};

export default RejectRequest;
