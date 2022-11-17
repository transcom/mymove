import React from 'react';
import { func, object, node, string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ReviewWeightTicket.module.scss';

import ReviewItems from 'components/Customer/PPM/Closeout/ReviewItems/ReviewItems';
// import DocumentViewerSidebar from 'pages/Office/DocumentViewerSidebar/DocumentViewerSidebar';
import { formatAboutYourPPMItem } from 'utils/ppmCloseout';
// import officeShapes from 'types/officeShapes';

export default function ReviewWeightTicket({ mtoShipment, onClose }) {
  const aboutYourPPM = formatAboutYourPPMItem(mtoShipment?.ppmShipment);
  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1>PPM 1</h1>
          <ReviewItems contents={aboutYourPPM} />
        </div>
        <Button className={styles.closeButton} data-testid="closeSidebar" onClick={onClose} type="button" unstyled>
          <FontAwesomeIcon icon="times" title="Close sidebar" aria-label="Close sidebar" />
        </Button>
      </header>
      Hi
    </div>
  );
}

ReviewWeightTicket.propTypes = {
  mtoShipment: object.isRequired,
  onClose: func.isRequired,
};

// ReviewWeightTicket.defaultProps = {
//   description: '',
// };

// ReviewWeightTicket.Footer = function Footer({ children }) {
//   return <footer className={styles.footer}>{children}</footer>;
// };

// ReviewWeightTicket.Footer.propTypes = {
//   children: node.isRequired,
// };
