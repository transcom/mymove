import React from 'react';

// import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
// import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';
// import { Button } from '@trussworks/react-uswds';
import fakeDoc from 'shared/images/fake-doc.png';
import styles from './content.module.scss';

const DocViewerContent = () => (
  <div className={`${styles.docViewerContent}`}>
    <div className={`${styles.titleBar}`}>I am title bar</div>
    <div className={`${styles.docArea}`}>
      <img src={`${fakeDoc}`} alt=" " />
    </div>
    <div className={`${styles.docControls}`}>I am document controls</div>
  </div>
);

export default DocViewerContent;
