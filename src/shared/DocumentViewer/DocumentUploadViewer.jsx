import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { selectMoveDocument } from 'shared/Entities/modules/moveDocuments';
import DocumentContent from './DocumentContent';

export const DocumentUploadViewer = ({ moveDocument }) => {
  const uploadModels = get(moveDocument, 'document.uploads', []);
  return (
    <div className="document-contents">
      {uploadModels.map(({ url, filename, content_type }) => (
        <DocumentContent
          key={url}
          url={url}
          filename={filename}
          content_type={content_type}
        />
      ))}
    </div>
  );
};

DocumentUploadViewer.propTypes = {};

function mapStateToProps(state, props) {
  const moveDocumentId = props.match.params.moveDocumentId;
  return {
    moveDocument: selectMoveDocument(state, moveDocumentId),
  };
}
export default connect(mapStateToProps)(DocumentUploadViewer);
