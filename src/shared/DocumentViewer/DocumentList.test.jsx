import React from 'react';
import { shallow } from 'enzyme';
import DocumentList from './DocumentList';

describe('DocumentList tests', () => {
  it('has a link to upload a new document', () => {
    const newDocumentUrl = 'test-url-new';
    const defaultMoveDocument = {
      id: '',
      createdAt: '',
      notes: '',
      status: '',
      title: '',
      type: '',
    };
    let wrapper = shallow(
      <DocumentList
        currentMoveDocumentId=""
        detailUrlPrefix="/moves/1/documents"
        moveDocuments={[defaultMoveDocument]}
        uploadDocumentUrl={newDocumentUrl}
      />,
    );
    expect(
      wrapper
        .find('.document-upload-link')
        .find('Link')
        .prop('to'),
    ).toEqual(newDocumentUrl);
  });
});
