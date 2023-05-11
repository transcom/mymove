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
    const wrapper = shallow(
      <DocumentList
        currentMoveDocumentId=""
        detailUrlPrefix="/moves/1/documents"
        moveDocuments={[defaultMoveDocument]}
        uploadDocumentUrl={newDocumentUrl}
        moveId="1"
      />,
    );
    expect(wrapper.find('.document-upload-link').find('a').prop('href')).toEqual(newDocumentUrl);
  });
});
