import React from 'react';

import { ReactComponent as ExternalLink } from 'shared/icon/external-link.svg';
import { ReactComponent as ZoomIn } from 'shared/icon/zoom-in.svg';
import { ReactComponent as ZoomOut } from 'shared/icon/zoom-out.svg';
import { ReactComponent as RotateLeft } from 'shared/icon/rotate-counter-clockwise.svg';
import { ReactComponent as RotateRight } from 'shared/icon/rotate-clockwise.svg';
import { ReactComponent as ArrowLeft } from 'shared/icon/arrow-left.svg';
import { ReactComponent as ArrowRight } from 'shared/icon/arrow-right.svg';
import { Button } from '@trussworks/react-uswds';
import fakeDoc from 'shared/images/fake-doc.png';
import styles from './content.module.scss';

const DocViewerContent = () => (
  <div className={`${styles.docViewerContent}`}>
    <div className={`${styles.titleBar}`}>
      <p>ThisIsAVeryLongDocumentTitle.png</p>
      <Button unstyled>
        Open in a new window
        <ExternalLink />
      </Button>
    </div>
    <div className={`${styles.docArrows}`}>
      <Button unstyled className={`${styles.arrowButton}`}>
        <ArrowLeft />
      </Button>
      <Button unstyled className={`${styles.arrowButton}`}>
        <ArrowRight />
      </Button>
    </div>
    <div className={`${styles.docArea}`}>
      <img src={`${fakeDoc}`} alt=" " />
    </div>
    <div className={`${styles.docControls}`}>
      <Button unstyled>
        <RotateLeft />
        Rotate left
      </Button>
      <Button unstyled>
        <RotateRight />
        Rotate right
      </Button>
      <Button unstyled>
        <ZoomOut />
        Zoom out
      </Button>
      <Button unstyled>
        <ZoomIn />
        Zoom in
      </Button>
    </div>
  </div>
);

export default DocViewerContent;
