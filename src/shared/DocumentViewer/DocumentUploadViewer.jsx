import React from 'react';
import PropTypes from 'prop-types';
import { get } from 'lodash';
import { connect } from 'react-redux';

import { selectMoveDocument } from 'shared/Entities/modules/moveDocuments';
import DocumentContent from './DocumentContent';
import { ExistingUploadsShape } from 'types/uploads';
import withRouter from 'utils/routing';

export const DocumentUploadViewer = ({ moveDocument }) => {
  const uploadModels = get(moveDocument, 'document.uploads', []);
  return uploadModels.map(({ url, filename, contentType, status }) => (
    <DocumentContent key={url} url={url} filename={filename} contentType={contentType} status={status} />
  ));
};

const { shape, string } = PropTypes;

DocumentUploadViewer.propTypes = {
  moveDocument: shape({
    document: shape({
      id: string.isRequired,
      service_member_id: string.isRequired,
      uploads: ExistingUploadsShape.isRequired,
    }),
    id: string.isRequired,
    move_document_type: string.isRequired,
    move_id: string.isRequired,
    notes: string,
    personally_procured_move_id: string,
    status: string.isRequired,
    title: string.isRequired,
  }).isRequired,
};

function mapStateToProps(state, { router: { params } }) {
  const moveDocumentId = params.moveDocumentId;
  return {
    moveDocument: selectMoveDocument(state, moveDocumentId),
  };
}
export default withRouter(connect(mapStateToProps)(DocumentUploadViewer));
