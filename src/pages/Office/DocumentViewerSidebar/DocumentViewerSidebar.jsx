import React from 'react';
import { func, node, string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './DocumentViewerSidebar.module.scss';

export default function DocumentViewerSidebar({ children, description, title, subtitle, supertitle, onClose }) {
  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          {supertitle && <h2 className={styles.supertitle}>{supertitle}</h2>}
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
  supertitle: string,
  description: string,
  onClose: func.isRequired,
};

DocumentViewerSidebar.defaultProps = {
  subtitle: '',
  supertitle: '',
  description: '',
};

DocumentViewerSidebar.Content = function Content({ children }) {
  return <main>{children}</main>;
};

DocumentViewerSidebar.Content.propTypes = {
  children: node.isRequired,
};

DocumentViewerSidebar.Footer = function Footer({ children }) {
  return <footer className={styles.footer}>{children}</footer>;
};

DocumentViewerSidebar.Footer.propTypes = {
  children: node.isRequired,
};
