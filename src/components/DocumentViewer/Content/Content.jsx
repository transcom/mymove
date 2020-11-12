import React from 'react';
import PropTypes from 'prop-types';
import FileViewer from '@trussworks/react-file-viewer';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearchPlus as ZoomIn, faSearchMinus as ZoomOut } from '@fortawesome/free-solid-svg-icons';

import styles from './Content.module.scss';
// TODO
/*
import { ReactComponent as RotateLeft } from 'shared/icon/rotate-counter-clockwise.svg';
import { ReactComponent as RotateRight } from 'shared/icon/rotate-clockwise.svg';
import { ReactComponent as ArrowLeft } from 'shared/icon/arrow-left.svg';
import { ReactComponent as ArrowRight } from 'shared/icon/arrow-right.svg';
*/

const DocViewerContent = ({ fileType, filePath }) => (
  <div data-testid="DocViewerContent" className={styles.DocViewerContent}>
    <FileViewer
      key={`fileViewer_${filePath}`}
      fileType={fileType}
      filePath={filePath}
      renderControls={({ handleZoomIn, handleZoomOut }) => {
        return (
          <div className={styles.controls}>
            <Button type="button" unstyled onClick={handleZoomOut}>
              <FontAwesomeIcon icon={ZoomOut} title="Zoom out" aria-label="Zoom out" />
              Zoom out
            </Button>
            <Button type="button" unstyled onClick={handleZoomIn}>
              <FontAwesomeIcon icon={ZoomIn} title="Zoom in" aria-label="Zoom in" />
              Zoom in
            </Button>
          </div>
        );
      }}
    />
  </div>
);

DocViewerContent.propTypes = {
  filePath: PropTypes.string.isRequired,
  fileType: PropTypes.string.isRequired,
};

export default DocViewerContent;

/**
 * TODO
 *
 * - add className prop to file viewer
 * - add rotate left/right:
 *  <Button unstyled>
        <RotateLeft />
        Rotate left
      </Button>
      <Button unstyled>
        <RotateRight />
        Rotate right
      </Button>
 * - implement pagination for multi-page PDFs & nav render prop:
 *  <div className={`${styles.docArrows}`}>
      <Button unstyled className={`${styles.arrowButton}`}>
        <ArrowLeft />
      </Button>
      <Button unstyled className={`${styles.arrowButton}`}>
        <ArrowRight />
      </Button>
    </div>
*/
