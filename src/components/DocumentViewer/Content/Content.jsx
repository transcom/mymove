import React from 'react';
import PropTypes from 'prop-types';
import FileViewer from '@transcom/react-file-viewer';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './Content.module.scss';

const DocViewerContent = ({
  fileType,
  filePath,
  saveRotation,
  setRotationValue,
  rotationValue,
  disableSaveButton,
  onError,
}) => (
  <div data-testid="DocViewerContent" className={styles.DocViewerContent}>
    <FileViewer
      key={`fileViewer_${filePath}`}
      fileType={fileType}
      filePath={filePath}
      onError={onError}
      saveRotation={saveRotation}
      rotationValue={rotationValue}
      setRotationValue={setRotationValue}
      renderControls={({ handleZoomIn, handleZoomOut, handleRotateLeft, handleRotateRight, zoomPercentage }) => {
        return (
          <div className={styles.controls}>
            {zoomPercentage && (
              <Button
                data-testid="currentZoomPercentage"
                type="button"
                unstyled
                className={styles.zoomDisplayButton}
                aria-label={`Zoom level: ${zoomPercentage}%`}
              >
                Zoom: {zoomPercentage}%
              </Button>
            )}
            <Button data-testid="zoomOutButton" type="button" unstyled onClick={handleZoomOut}>
              <FontAwesomeIcon icon="search-minus" title="Zoom out" aria-label="Zoom out" />
              Zoom out
            </Button>
            <Button data-testid="zoomInButton" type="button" unstyled onClick={handleZoomIn}>
              <FontAwesomeIcon icon="search-plus" title="Zoom in" aria-label="Zoom in" />
              Zoom in
            </Button>
            {['jpg', 'jpeg', 'gif', 'png', 'pdf'].includes(fileType) && (
              <>
                <Button data-testid="rotateLeftButton" type="button" unstyled onClick={handleRotateLeft}>
                  <FontAwesomeIcon icon="rotate-left" title="Rotate left" aria-label="Rotate left" />
                  Rotate left
                </Button>
                <Button data-testid="rotateRightButton" type="button" unstyled onClick={handleRotateRight}>
                  <FontAwesomeIcon icon="rotate-right" title="Rotate right" aria-label="Rotate right" />
                  Rotate right
                </Button>
                <Button type="button" unstyled disabled={disableSaveButton} onClick={saveRotation}>
                  <svg
                    height="24"
                    viewBox="0 0 24 24"
                    style={{
                      textDecoration: 'none',
                      color: disableSaveButton ? 'transparent' : 'inherit',
                      visibility: disableSaveButton ? 'hidden' : 'visible',
                    }}
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M21 12.5L16.5 17L11 12.5L16.5 8L21 12.5Z" />
                  </svg>
                  <span style={{ textDecoration: 'none' }}>Save</span>
                </Button>
              </>
            )}
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
