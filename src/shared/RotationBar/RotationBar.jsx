import React, { Component } from 'react';
import leftRotation from 'shared/images/left-rotation.png';
import rightRotation from 'shared/images/right-rotation.png';
import PropTypes from 'prop-types';
import styles from './RotationBar.module.scss';

export class RotationBar extends Component {
  render() {
    return (
      <div className={styles['rotation-bar']}>
        <button onClick={this.props.onLeftButtonClick}>
          <img src={leftRotation} alt="rotate-left" />
        </button>
        <button onClick={this.props.onRightButtonClick}>
          <img src={rightRotation} alt="rotate-right" />
        </button>
      </div>
    );
  }
}

RotationBar.propTypes = {
  onLeftButtonClick: PropTypes.func,
  onRightButtonClick: PropTypes.func,
};
