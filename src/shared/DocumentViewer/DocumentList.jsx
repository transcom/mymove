import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Link } from 'react-router-dom-old';

import { renderStatusIcon, openLinkInNewWindow } from 'shared/utils';
import styles from 'shared/DocumentViewer/DocumentList.module.scss';
import { defaultRelativeWindowSize } from 'shared/constants';

const DocumentList = ({
  currentMoveDocumentId,
  moveDocuments,
  detailUrlPrefix,
  disableLinks,
  uploadDocumentUrl,
  moveId,
}) => (
  <div>
    {moveDocuments.map((doc) => {
      const chosenDocument = currentMoveDocumentId === doc.id ? styles['chosen-document'] : 'usa-link link-blue';
      const status = renderStatusIcon(doc.status);
      const detailUrl = `${detailUrlPrefix}/${doc.id}`;
      return (
        <div className={`panel-field ${chosenDocument}`} data-testid="doc-link" key={doc.id}>
          <span className="status">{status}</span>
          {!disableLinks &&
            (window.name === `docViewer-${moveId}` ? (
              // open in same window if already in document viewer
              <Link className={chosenDocument} to={detailUrl}>
                {doc.title}
              </Link>
            ) : (
              // open in new window if one is not already open
              <a
                href={detailUrl}
                target={`docViewer-${moveId}`}
                className={chosenDocument}
                onClick={openLinkInNewWindow.bind(
                  this,
                  detailUrl,
                  `docViewer-${moveId}`,
                  window,
                  defaultRelativeWindowSize,
                )}
              >
                {doc.title}
              </a>
            ))}
          {disableLinks && <span>{doc.title}</span>}
        </div>
      );
    })}
    {window.name === `docViewer-${moveId}` ? (
      <div className={`${styles['document-upload-link']} link-blue`}>
        <Link to={uploadDocumentUrl}>
          <FontAwesomeIcon className="icon link-blue" icon="plus-circle" />
          Upload new document
        </Link>
      </div>
    ) : (
      <div className={styles['document-upload-link']} data-testid="document-upload-link">
        <a
          href={uploadDocumentUrl}
          target={`docViewer-${moveId}`}
          onClick={openLinkInNewWindow.bind(
            this,
            uploadDocumentUrl,
            `docViewer-${moveId}`,
            window,
            defaultRelativeWindowSize,
          )}
          className="usa-link"
        >
          <FontAwesomeIcon className="icon link-blue" icon="plus-circle" />
          Upload new document
        </a>
      </div>
    )}
  </div>
);

DocumentList.propTypes = {
  currentMoveDocumentId: PropTypes.string,
  detailUrlPrefix: PropTypes.string.isRequired,
  disableLinks: PropTypes.bool,
  moveId: PropTypes.string,
  moveDocuments: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      status: PropTypes.string.isRequired,
      title: PropTypes.string.isRequired,
    }),
  ).isRequired,
  uploadDocumentUrl: PropTypes.string.isRequired,
};

export default DocumentList;
