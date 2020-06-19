import React from 'react';
import { Button } from '@trussworks/react-uswds';

import styles from './menu.module.scss';

import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';

const testThumbnail = (
  <a className={styles.thumbnailItem} href="#">
    <p>ASamplePrettyLongDocumentTitle.png</p>
    <div className={styles.thumbnailImage} />
  </a>
);

const DocViewerMenu = () => (
  <div className={`${styles.docViewerMenu}`}>
    <div className={styles.fixedContainer}>
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
  </div>
);

export default DocViewerMenu;
