import { Button } from '@trussworks/react-uswds';
import React from 'react';
import classnames from 'classnames';

import styles from './ButtonUsa.module.scss';

export const usaButtonStyle = styles['usa-button'];
export const primaryButtonStyle = styles['primary'];
export const outlineButtonStyle = styles['outline'];
export const destructiveButtonStyle = styles['destructive'];
export const destructiveOutlineButtonStyle = styles['destructive-outline'];

export const ButtonUsa = ({ className, ...props }) => {
  const cn = classnames([usaButtonStyle, className || []].flat()); //.join(' ');

  return React.cloneElement(<Button {...props} />, {
    // className: cn,
    class: cn,
    ...props,
  });
};
