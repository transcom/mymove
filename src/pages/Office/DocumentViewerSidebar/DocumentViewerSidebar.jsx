import React from 'react';
import { bool, func, object, node, string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './DocumentViewerSidebar.module.scss';

import ShipmentModificationTag from 'components/ShipmentModificationTag/ShipmentModificationTag';
import { shipmentModificationTypes } from 'constants/shipments';

export default function DocumentViewerSidebar({
  children,
  description,
  title,
  subtitle,
  supertitle,
  onClose,
  defaultH3,
  hyperlink,
  showDiversionModificationTag,
}) {
  return (
    <div
      className={classnames({
        [styles.container]: true,
        [styles.defaultH3]: defaultH3,
      })}
    >
      <header className={styles.header}>
        <div>
          {supertitle && <h2 className={styles.supertitle}>{supertitle}</h2>}
          <h1>
            {title}
            {showDiversionModificationTag && (
              <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.DIVERSION} />
            )}
          </h1>
          {subtitle && <h2>{subtitle}</h2>}
          {description && <h3>{description}</h3>}
          {hyperlink && <p data-testid="hyperlink">{hyperlink}</p>}
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
  defaultH3: bool,
  showDiversionModificationTag: bool,
};

DocumentViewerSidebar.defaultProps = {
  subtitle: '',
  supertitle: '',
  description: '',
  defaultH3: false,
  showDiversionModificationTag: false,
};

DocumentViewerSidebar.Content = function Content({ children, mainRef }) {
  return <main ref={mainRef}>{children}</main>;
};

DocumentViewerSidebar.Content.propTypes = {
  children: node.isRequired,
  mainRef: object,
};

DocumentViewerSidebar.Content.defaultProps = {
  mainRef: null,
};

DocumentViewerSidebar.Footer = function Footer({ children }) {
  return <footer className={styles.footer}>{children}</footer>;
};

DocumentViewerSidebar.Footer.propTypes = {
  children: node.isRequired,
};
