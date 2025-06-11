import React from 'react';
import { shallow, mount } from 'enzyme';

import DocumentContent, { NonPDFImage, PDFImage } from './DocumentContent';

import Alert from 'shared/Alert';
import { UPLOAD_SCAN_STATUS } from 'shared/constants';

describe('DocumentContent', () => {
  describe('conditionally renders components based on content type', () => {
    it('renders a PDFImage when content type is pdf', () => {
      const wrapper = shallow(<DocumentContent contentType="application/pdf" url="www" filename="filename" />);
      expect(wrapper.find(PDFImage)).toHaveLength(1);
    });
    it('renders a NonPDFImage when content type is not pdf', () => {
      const wrapper = shallow(<DocumentContent contentType="image/jpeg" url="www" filename="filename" />);
      expect(wrapper.find(NonPDFImage)).toHaveLength(1);
    });
    it('renders an Alert when tags indicate document is infected', () => {
      const wrapper = shallow(
        <DocumentContent
          contentType="application/pdf"
          url="www"
          filename="filename"
          status={UPLOAD_SCAN_STATUS.LEGACY_INFECTED}
        />,
      );
      expect(wrapper.find(Alert)).toHaveLength(1);
    });
  });
});

describe('NonPDFImage', () => {
  describe('rotation', () => {
    it('renders with rotation of zero', () => {
      const wrapper = mount(<NonPDFImage src="url" />);
      expect(wrapper.state().rotation).toBe(0);
    });
    it('clicking rotate right rotates 90 degrees', () => {
      const wrapper = mount(<NonPDFImage src="url" />);
      const nonPdfImage = wrapper.find(NonPDFImage);
      nonPdfImage.instance().rotateRight();
      nonPdfImage.update();

      expect(wrapper.state().rotation).toBe(90);
    });

    it('clicking rotate left rotates 90 degrees to the left', () => {
      const wrapper = mount(<NonPDFImage src="url" />);
      const nonPdfImage = wrapper.find(NonPDFImage);
      nonPdfImage.instance().rotateLeft();
      nonPdfImage.update();

      expect(wrapper.state().rotation).toBe(-90);
    });

    it('clicking rotate left twice rotates 180 degrees to the left', () => {
      const wrapper = mount(<NonPDFImage src="url" />);
      const nonPdfImage = wrapper.find(NonPDFImage);
      nonPdfImage.instance().rotateLeft();
      nonPdfImage.instance().rotateLeft();
      nonPdfImage.update();

      expect(wrapper.state().rotation).toBe(-180);
    });
  });
});
