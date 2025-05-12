import { Button } from '@trussworks/react-uswds';
import React from 'react';
import styles from './buttons.module.scss';
import classnames from 'classnames';

const usaStyle = styles['usa-button'];

export const ButtonUsa = ({ className, ...props }) => {
  const cn = classnames([usaStyle, className || []].flat()); //.join(' ');
  return React.cloneElement(<Button {...props} />, {
    // className: cn,
    class: cn,
    ...props,
  });
};
