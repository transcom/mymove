import React from 'react';
import PropTypes from 'prop-types';

import styles from './RotationBar.module.scss';

import leftRotation from 'shared/images/left-rotation.png';
import rightRotation from 'shared/images/right-rotation.png';

export const RotationBar = (props) => (
  <div className={styles['rotation-bar']}>
    <button className="usa-button" onClick={props.onLeftButtonClick}>
      <img src={leftRotation} alt="rotate-left" />
    </button>
    <button className="usa-button" onClick={props.onRightButtonClick}>
      <img src={rightRotation} alt="rotate-right" />
    </button>
  </div>
);

RotationBar.propTypes = {
  onLeftButtonClick: PropTypes.func.isRequired,
  onRightButtonClick: PropTypes.func.isRequired,
};
