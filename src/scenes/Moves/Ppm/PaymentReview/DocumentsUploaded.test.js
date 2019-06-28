import React from 'react';
import { DocumentsUploaded } from './DocumentsUploaded';
import { mount } from 'enzyme';
import Alert from 'shared/Alert';

function generateWrapper(allDocuments, mockGetMoveDocumentsForMove) {
  return mount(<DocumentsUploaded allDocuments={allDocuments} getMoveDocumentsForMove={mockGetMoveDocumentsForMove} />);
}

describe('DocumentsUploaded Alert', () => {
  describe('No documents uploaded', () => {
    it('component does not render', () => {
      const mockGetMoveDocumentsForMove = jest.fn();
      const wrapper = generateWrapper([], mockGetMoveDocumentsForMove);

      expect(wrapper.find(Alert).length).toEqual(0);
    });
  });
  describe('One document uploaded', () => {
    it('component renders', () => {
      const mockGetMoveDocumentsForMove = jest.fn();
      const wrapper = generateWrapper(['document'], mockGetMoveDocumentsForMove);

      expect(mockGetMoveDocumentsForMove).toHaveBeenCalled();
      expect(wrapper.find(Alert).length).toEqual(1);
      expect(wrapper.find(Alert).text()).toContain('1 document added');
    });
  });
  describe('More than one document uploaded', () => {
    it('component renders and text uses "documents" instead of "document"', () => {
      const mockGetMoveDocumentsForMove = jest.fn();
      const wrapper = generateWrapper(['document 1', 'document 2'], mockGetMoveDocumentsForMove);

      expect(mockGetMoveDocumentsForMove).toHaveBeenCalled();
      expect(wrapper.find(Alert).length).toEqual(1);
      expect(wrapper.find(Alert).text()).toContain('2 documents added');
    });
  });
});
