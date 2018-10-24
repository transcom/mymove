import React from 'react';
import { shallow } from 'enzyme';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { DeliveryDateFormView } from './ShipmentInfo';

const defaultProps = {
  schema: {},
  onCancel() {},
  handleSubmit() {},
  submitting: false,
  valid: false,
};

describe('DeliveryDateForm component test', () => {
  describe('renders with the correct fields', () => {
    const wrapper = shallow(<DeliveryDateFormView {...defaultProps} />);
    it('should have a header', () => {
      expect(wrapper.find('.infoPanel-wizard-header').text()).toEqual('Enter Delivery');
    });
    it('should have an actual delivery date swagger field', () => {
      expect(wrapper.find(SwaggerField).props()).toEqual({
        fieldName: 'actual_delivery_date',
        swagger: defaultProps.schema,
        required: true,
      });
    });
    it('should have upload origin documents help text', () => {
      expect(wrapper.find('.infoPanel-wizard-help').text()).toEqual(
        'After clicking "Done", please upload the destination docs. Use the "Upload new document" link in the Documents panel at right.',
      );
    });
    it('should have a cancel link', () => {
      expect(wrapper.find('.infoPanel-wizard-cancel').text()).toEqual('Cancel');
    });
    it('should have a done button', () => {
      expect(wrapper.find('.usa-button-primary').text()).toEqual('Done');
    });
  });
});
