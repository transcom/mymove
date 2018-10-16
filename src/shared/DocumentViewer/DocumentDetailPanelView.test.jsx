import React from 'react';
import { shallow } from 'enzyme';
import DocumentDetailPanelView from './DocumentDetailPanelView';

describe('DocumentDetailPanelView', () => {
  const renderDocumentPanel = ({
    createdAt = '',
    isExpenseDocument = false,
    notes = '',
    schema = {},
    status = '',
    title = '',
    type = '',
  }) =>
    shallow(
      <DocumentDetailPanelView
        createdAt={createdAt}
        notes={notes}
        schema={schema}
        status={status}
        title={title}
        type={type}
      />,
    );

  describe('panel-subhead', () => {
    it('includes the status of the move Document', () => {
      const status = 'SUBMITTED';
      const documentPanel = renderDocumentPanel({ status });
      // status is rendered as a FontAwesome icon
      expect(documentPanel.find('.panel-subhead').text()).toContain('FontAwesomeIcon');
    });

    it('includes the title of the document', () => {
      const title = 'My Title';
      const documentPanel = renderDocumentPanel({ title });
      expect(documentPanel.find('.panel-subhead').text()).toContain(title);
    });
  });

  it('includes the formatted uploaded-at time', () => {
    const createdAt = '2018-09-27 08:14:38.702434';
    const documentPanel = renderDocumentPanel({ createdAt });
    expect(documentPanel.find('.uploaded-at').text()).toEqual('Uploaded 27-Sep-18');
  });

  describe('swagger fields documents', () => {
    const schema = { properties: {} };

    it('includes the document title', () => {
      const title = 'My Title';
      const documentPanel = renderDocumentPanel({ schema, title });
      const swaggerField = documentPanel.find({ fieldName: 'title' });
      expect(swaggerField.prop('title')).toEqual('Document Title');
      expect(swaggerField.prop('required')).toEqual(true);
      expect(swaggerField.prop('values')).toEqual({ title });
      expect(swaggerField.prop('schema')).toEqual(schema);
    });

    it('includes the document type', () => {
      const type = 'GBL';
      const documentPanel = renderDocumentPanel({ schema, type });
      const swaggerField = documentPanel.find({
        fieldName: 'move_document_type',
      });
      expect(swaggerField.prop('title')).toEqual('Document Type');
      expect(swaggerField.prop('required')).toEqual(true);
      expect(swaggerField.prop('values')).toEqual({ move_document_type: type });
      expect(swaggerField.prop('schema')).toEqual(schema);
    });

    it('includes document status', () => {
      const status = 'SUBMITTED';
      const documentPanel = renderDocumentPanel({ schema, status });
      const swaggerField = documentPanel.find({ fieldName: 'status' });
      expect(swaggerField.prop('title')).toEqual('Document Status');
      expect(swaggerField.prop('required')).toEqual(true);
      expect(swaggerField.prop('values')).toEqual({ status });
      expect(swaggerField.prop('schema')).toEqual(schema);
    });

    it('includes document status', () => {
      const notes = 'More info about this doc';
      const documentPanel = renderDocumentPanel({ schema, notes });
      const swaggerField = documentPanel.find({ fieldName: 'notes' });
      expect(swaggerField.prop('title')).toEqual('Notes');
      expect(swaggerField.prop('values')).toEqual({ notes });
      expect(swaggerField.prop('schema')).toEqual(schema);
    });
  });
});
