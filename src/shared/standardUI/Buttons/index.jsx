import React from 'react';
import { ButtonUsa as Button } from './ButtonUsa';
import { Link } from 'react-router-dom';

const mainButtonClass = [];

export const Basic = ({ children, mainClassStyles: mainStyles = [mainButtonClass], className, ...props }) => {
  const mainClasses = Array.isArray(mainStyles) ? mainStyles : [mainStyles];
  const classNameValue = [mainClasses, className || []].flat().join(' ');
  return (
    <Button {...props} className={classNameValue}>
      {children}
    </Button>
  );
};

export const LinkButton = ({ children, href, to, state, ...props }) => {
  const elem = React.cloneElement(<Link {...{ to, state, ...props }} />, {
    ...props,
    style: { display: 'contents' },
    children: <button {...props}>{children}</button>,
  });

  return elem;
};
