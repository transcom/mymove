import React from 'react';
import PropTypes from 'prop-types';
import { renderStatusIcon } from 'shared/utils';
import { formatDate } from 'shared/formatters';
import { PanelSwaggerField } from 'shared/EditablePanel';

const DocumentDetailPanelView = ({ createdAt, notes, schema, status, title, type }) => (
  <div>
    <span className="panel-subhead">
      {renderStatusIcon(status)}
      {title}
    </span>
    <p className="uploaded-at">{`Uploaded ${formatDate(createdAt)}`}</p>
    <PanelSwaggerField title="Document Title" fieldName="title" required schema={schema} values={{ title }} />
    <PanelSwaggerField
      title="Document Type"
      fieldName="move_document_type"
      required
      schema={schema}
      values={{ move_document_type: type }}
    />
    <PanelSwaggerField title="Document Status" fieldName="status" required schema={schema} values={{ status }} />
    <PanelSwaggerField title="Notes" fieldName="notes" schema={schema} values={{ notes }} />
  </div>
);

DocumentDetailPanelView.propTypes = {
  createdAt: PropTypes.string.isRequired,
  notes: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  status: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
};

export default DocumentDetailPanelView;
