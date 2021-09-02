import React from 'react';
import { func, node, string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './DocumentViewerSidebar.module.scss';

export default function DocumentViewerSidebar({ children, description, title, subtitle, onClose }) {
  return (
    <div className={styles.container}>
      <header>
        <div>
          <h1>{title}</h1>
          {subtitle && <h2>{subtitle}</h2>}
          {description && <h3>{description}</h3>}
        </div>
        <Button className={styles.closeButton} data-testid="closeSidebar" onClick={onClose} type="button" unstyled>
          <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
        </Button>
      </header>
      {children}
    </div>
  );
}

DocumentViewerSidebar.propTypes = {
  children: node.isRequired,
  title: string.isRequired,
  subtitle: string,
  description: string,
  onClose: func.isRequired,
};

DocumentViewerSidebar.defaultProps = {
  subtitle: '',
  description: '',
};

DocumentViewerSidebar.Content = function Content({ children }) {
  return <main>{children}</main>;
};

DocumentViewerSidebar.Content.propTypes = {
  children: node.isRequired,
};

DocumentViewerSidebar.Footer = function Footer({ children }) {
  return <footer>{children}</footer>;
};

DocumentViewerSidebar.Footer.propTypes = {
  children: node.isRequired,
};
