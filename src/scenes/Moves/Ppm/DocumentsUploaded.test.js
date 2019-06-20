import React from 'react';
import { DocumentsUploaded } from './DocumentsUploaded';
import { mount } from 'enzyme';
import Alert from 'shared/Alert';

describe('Documents Uploaded Alert', () => {
  describe('No documents uploaded', () => {
    it('does not render DocumentUploaded', () => {
      const wrapper = mount(<DocumentsUploaded allDocuments={[]} getMoveDocumentsForMove={jest.fn()} />);
      expect(wrapper.find(Alert).length).toEqual(0);
    });
  });
  describe('One document uploaded', () => {
    it('renders DocumentUploaded', () => {
      const mockGetMoveDocumentsForMove = jest.fn();
      const wrapper = mount(
        <DocumentsUploaded allDocuments={['document']} getMoveDocumentsForMove={mockGetMoveDocumentsForMove} />,
      );
      expect(mockGetMoveDocumentsForMove).toHaveBeenCalled();
      expect(wrapper.find(Alert).length).toEqual(1);
      expect(wrapper.find(Alert).text()).toContain('1 document added');
    });
  });
  describe('More than one document uploaded', () => {
    it('renders DocumentUploaded with "documents" in message', () => {
      const mockGetMoveDocumentsForMove = jest.fn();
      const wrapper = mount(
        <DocumentsUploaded
          allDocuments={['document 1', 'document 2']}
          getMoveDocumentsForMove={mockGetMoveDocumentsForMove}
        />,
      );
      expect(mockGetMoveDocumentsForMove).toHaveBeenCalled();
      expect(wrapper.find(Alert).length).toEqual(1);
      expect(wrapper.find(Alert).text()).toContain('2 documents added');
    });
  });
});
