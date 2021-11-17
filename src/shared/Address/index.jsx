import React, { Component } from 'react';
import { FormSection } from 'redux-form';
import PropTypes from 'prop-types';
import style from './index.module.scss';

import { PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const AddressElementDisplay = ({ address, title }) => (
  <PanelField title={title}>
    {address.streetAddress1 && (
      <span>
        {address.streetAddress1}
        <br />
      </span>
    )}
    {address.streetAddress2 && (
      <span>
        {address.streetAddress2}
        <br />
      </span>
    )}
    {address.city}, {address.state} {address.postalCode}
  </PanelField>
);

AddressElementDisplay.defaultProps = {
  address: {},
};

AddressElementDisplay.propTypes = {
  address: PropTypes.shape({
    streetAddress1: PropTypes.string,
    streetAddress2: PropTypes.string,
    city: PropTypes.string.isRequired,
    state: PropTypes.string.isRequired,
    postalCode: PropTypes.string.isRequired,
  }).isRequired,
  title: PropTypes.string.isRequired,
};

export class AddressElementEdit extends Component {
  render() {
    return (
      <FormSection name={this.props.fieldName}>
        <div className="panel-subhead">{this.props.title}</div>
        <SwaggerField fieldName="streetAddress1" swagger={this.props.schema} required />
        <SwaggerField fieldName="streetAddress2" swagger={this.props.schema} />
        <SwaggerField fieldName="city" swagger={this.props.schema} required />
        <div className={style['state-zip-container']}>
          <div className="usa-width-one-half">
            <SwaggerField fieldName="state" swagger={this.props.schema} required />
          </div>
          <div className="usa-width-one-half">
            <SwaggerField
              fieldName="postalCode"
              swagger={this.props.schema}
              zipPattern={this.props.zipPattern}
              required
            />
          </div>
        </div>
      </FormSection>
    );
  }
}

AddressElementEdit.propTypes = {
  fieldName: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  zipPattern: PropTypes.string,
  title: PropTypes.string,
};
