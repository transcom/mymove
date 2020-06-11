import React from 'react';
import propTypes from 'prop-types';

import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';
import { Button } from '@trussworks/react-uswds';
// import fakeDoc from 'shared/images/fake-doc.png';
import styles from './index.module.scss';

const testThumbnail = (
  <a className={styles.thumbnailItem}>
    <p>ASampleLongDocumentTitle.png</p>
    <div className={styles.thumbnailImage} />
  </a>
);

const DocViewerMenu = ({ menuClass }) => (
  <div className={`${styles.docViewerMenu} ${menuClass}`}>
    <div className={styles.menuHeader}>
      <h3>Documents</h3>
      <div className={styles.menuControls}>
        <Button unstyled className={styles.menuOpen}>
          <DocMenu />
        </Button>
        <Button unstyled className={styles.menuClose}>
          <XLightIcon />
        </Button>
      </div>
    </div>
    <div className={styles.thumbnailContainer}>
      {testThumbnail}
      {testThumbnail}
      {testThumbnail}
      {testThumbnail}
      {testThumbnail}
    </div>
  </div>
);

DocViewerMenu.propTypes = {
  menuClass: propTypes.string,
};

export default DocViewerMenu;
