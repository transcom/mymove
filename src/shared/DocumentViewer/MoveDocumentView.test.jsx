import React from 'react';
import { shallow } from 'enzyme';
import MoveDocumentView from './MoveDocumentView';

describe('MoveDocumentView', () => {
  const defaultMoveDocument = {
    createdAt: '',
    notes: '',
    status: '',
    title: '',
    type: '',
  };

  it('calls onDidMount when the component mounts', () => {
    const onDidMountSpy = jest.fn();
    shallow(
      <MoveDocumentView
        documentDetailUrlPrefix=""
        match={{ params: {} }}
        moveDocument={defaultMoveDocument}
        moveDocumentSchema={{}}
        moveDocumentId=""
        moveDocuments={[]}
        moveLocator=""
        newDocumentUrl=""
        onDidMount={onDidMountSpy}
        serviceMember={{ name: '', edipi: '' }}
        uploads={[]}
      />,
    );
    expect(onDidMountSpy).toBeCalled();
  });

  const renderMoveDocumentView = ({
    documentDetailUrlPrefix = '',
    moveDocument = defaultMoveDocument,
    moveDocumentSchema = {},
    moveDocumentId = '',
    moveDocuments = [],
    moveLocator = '',
    newDocumentUrl = '',
    serviceMember = { name: '', edipi: '' },
    uploads = [],
    match = { params: {} },
  }) =>
    shallow(
      <MoveDocumentView
        documentDetailUrlPrefix={documentDetailUrlPrefix}
        match={match}
        moveDocument={moveDocument}
        moveDocumentSchema={moveDocumentSchema}
        moveDocumentId={moveDocumentId}
        moveDocuments={moveDocuments}
        moveLocator={moveLocator}
        newDocumentUrl={newDocumentUrl}
        onDidMount={() => {}}
        serviceMember={serviceMember}
        uploads={uploads}
      />,
      { disableLifecycleMethods: true },
    );

  it('renders DocumentContent for each upload', () => {
    const uploads = [
      { url: 'http://test.pdf', filename: 'test.pdf', content_type: 'PDF' },
      { url: 'http://test2.pdf', filename: 'test2.pdf', content_type: 'PDF' },
    ];

    const documentView = renderMoveDocumentView({ uploads });
    expect(documentView.find('DocumentContent').length).toEqual(2);
  });

  it('renders the service member name', () => {
    const name = 'Doe, Jane';
    const serviceMember = { name, edipi: '' };
    const documentView = renderMoveDocumentView({ serviceMember });
    expect(
      documentView
        .find('.usa-width-one-third')
        .find('h3')
        .text(),
    ).toEqual(name);
  });

  it('renders the move locator', () => {
    const moveLocator = 'FBXY3M';
    const documentView = renderMoveDocumentView({ moveLocator });
    expect(documentView.find({ title: 'Move Locator' }).prop('children')).toEqual(moveLocator);
  });

  it('renders the serviceMember edipi', () => {
    const edipi = '999999999';
    const serviceMember = { edipi, name: '' };
    const documentView = renderMoveDocumentView({ serviceMember });
    expect(documentView.find({ title: 'DoD ID' }).prop('children')).toEqual(edipi);
  });

  describe('All Documents tab', () => {
    const moveDocument = { id: '', status: '', title: '' };
    const moveDocuments = [moveDocument, moveDocument, moveDocument];

    it('renders the All Documents tab header', () => {
      const documentView = renderMoveDocumentView({ moveDocuments });
      expect(
        documentView
          .find('Tab')
          .at(0)
          .dive()
          .text(),
      ).toEqual('All Documents (3)');
    });

    it('has a link to upload a new document', () => {
      const newDocumentUrl = 'test-url-new';
      const documentView = renderMoveDocumentView({ newDocumentUrl });
      const newDocumentLink = documentView
        .find('TabPanel')
        .at(0)
        .find('Link');
      expect(newDocumentLink.prop('to')).toEqual(newDocumentUrl);
      expect(newDocumentLink.prop('children')).toEqual('Upload new document');
    });

    it('renders a DocumentList with the appropriate props', () => {
      const documentDetailUrlPrefix = 'test-doc-prefix';
      const documentView = renderMoveDocumentView({
        documentDetailUrlPrefix,
        moveDocuments,
      });
      const documentList = documentView
        .find('TabPanel')
        .at(0)
        .find('DocumentList');
      expect(documentList.prop('detailUrlPrefix')).toEqual(documentDetailUrlPrefix);
      expect(documentList.prop('moveDocuments')).toEqual(moveDocuments);
    });
  });

  describe('Details tab', () => {
    it('renders the Details tab', () => {
      const documentView = renderMoveDocumentView({});
      expect(
        documentView
          .find('Tab')
          .at(1)
          .dive()
          .text(),
      ).toEqual('Details');
    });

    it('renders a DocumentDetailPanelView with the appropriate props', () => {
      const moveDocumentSchema = { properties: {} };
      const moveDocument = {
        ...defaultMoveDocument,
        title: 'My Document',
      };
      const documentView = renderMoveDocumentView({
        moveDocumentSchema,
        moveDocument,
      });
      const documentDetail = documentView
        .find('TabPanel')
        .at(1)
        .find('DocumentDetailPanelView');
      expect(documentDetail.prop('title')).toEqual('My Document');
      expect(documentDetail.prop('schema')).toEqual(moveDocumentSchema);
    });
  });
});
