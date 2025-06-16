/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';

import styles from './LoadingButton.module.scss';

const LoadingButton = ({
  onClick,
  isLoading,
  labelText,
  loadingText,
  loadingIcon,
  iconSpin,
  buttonClassName,
  ...props
}) => {
  return (
    <Button
      className={classNames(styles.loadingButton, buttonClassName)}
      data-testid="loading-button"
      onClick={onClick}
      {...props}
    >
      {isLoading ? (
        <>
          {loadingText}
          <FontAwesomeIcon icon={loadingIcon} spin={iconSpin} role="presentation" id={styles.loadingButtonIcon} />
        </>
      ) : (
        labelText
      )}
    </Button>
  );
};

LoadingButton.defaultProps = {
  labelText: 'Save',
  loadingText: 'Saving',
  loadingIcon: 'spinner',
  iconSpin: true,
  buttonClassName: '',
};

LoadingButton.propTypes = {
  onClick: PropTypes.func.isRequired,
  isLoading: PropTypes.bool.isRequired,
  labelText: PropTypes.string,
  loadingText: PropTypes.string,
  loadingIcon: PropTypes.string,
  iconSpin: PropTypes.bool,
  buttonClassName: PropTypes.string,
};

export default LoadingButton;
