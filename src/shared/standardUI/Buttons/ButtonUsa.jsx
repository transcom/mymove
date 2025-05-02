import { Button } from '@trussworks/react-uswds';
import React from 'react';
import styles from './buttons.module.scss';

const usaStyle = styles['usa-button'];

export const ButtonUsa = ({ className, ...props }) => {
  return React.cloneElement(<Button {...props} />, {
    class: [usaStyle, className || []].flat().join(' '),
    ...props,
  });
};
