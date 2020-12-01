import React from 'react';
import { DocumentsUploaded } from './DocumentsUploaded';
import { mount } from 'enzyme';

const initialProps = {
  moveId: 0,
  expenseDocs: [],
  weightTicketDocs: [],
  weightTicketSetDocs: [],
  getMoveDocumentsForMove: jest.fn(),
};
function generateWrapper(props) {
  return mount(<DocumentsUploaded {...initialProps} {...props} />);
}

describe('DocumentsUploaded Alert', () => {
  const documentsUploadedContainer = '[data-testid="documents-uploaded"]';
  describe('No documents uploaded', () => {
    it('component does not render', () => {
      const wrapper = generateWrapper();
      expect(wrapper.find(documentsUploadedContainer).length).toEqual(0);
    });
  });

  describe('One document uploaded', () => {
    it('component renders', () => {
      const mockGetMoveDocumentsForMove = jest.fn();
      const wrapper = generateWrapper({
        getMoveDocumentsForMove: mockGetMoveDocumentsForMove,
        expenseDocs: [{}],
      });
      expect(mockGetMoveDocumentsForMove).toHaveBeenCalled();
      expect(wrapper.find(documentsUploadedContainer).length).toEqual(1);
      expect(wrapper.find(documentsUploadedContainer).text()).toContain('1 document added');
    });
  });
  describe('More than one document uploaded', () => {
    it('component renders and text uses "documents" instead of "document"', () => {
      const mockGetMoveDocumentsForMove = jest.fn();
      const wrapper = generateWrapper({
        getMoveDocumentsForMove: mockGetMoveDocumentsForMove,
        expenseDocs: [{}],
        weightTicketSetDocs: [{}],
      });

      expect(mockGetMoveDocumentsForMove).toHaveBeenCalled();
      expect(wrapper.find(documentsUploadedContainer).length).toEqual(1);
      expect(wrapper.find(documentsUploadedContainer).text()).toContain('2 documents added');
    });
  });
});
