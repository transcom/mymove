import { FormSection } from 'redux-form';
import { get } from 'lodash';
import React, { Component } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import PropTypes from 'prop-types';

export class DimensionsField extends Component {
  render() {
    return (
      <FormSection name={this.props.fieldName}>
        <label htmlFor={this.props.fieldName} className="usa-input-label">
          <b>{this.props.labelText}</b>
        </label>
        <table className="dimensions-form">
          <tbody>
            <tr>
              <td width="26%" className="dimensions-form-header">
                Length
              </td>
              <td width="11%" />
              <td width="26%" className="dimensions-form-header">
                Width
              </td>
              <td width="11%" />
              <td width="26%" className="dimensions-form-header">
                Height
              </td>
            </tr>
            <tr>
              <td width="26%" className="dimensions-form-input-cell">
                <SwaggerField
                  fieldName="length"
                  swagger={get(this.props, 'swagger.properties.' + this.props.fieldName)}
                  hideLabel={true}
                  className="dimensions-form-input"
                  required={this.props.isRequired}
                />
              </td>
              <td width="11%" className="multiplication-sign" />
              <td width="26%" className="dimensions-form-input-cell">
                <SwaggerField
                  fieldName="width"
                  swagger={get(this.props, 'swagger.properties.' + this.props.fieldName)}
                  hideLabel={true}
                  className="dimensions-form-input"
                  required={this.props.isRequired}
                />
              </td>
              <td width="11%" className="multiplication-sign" />
              <td width="26%" className="dimensions-form-input-cell">
                <SwaggerField
                  fieldName="height"
                  swagger={get(this.props, 'swagger.properties.' + this.props.fieldName)}
                  hideLabel={true}
                  className="dimensions-form-input"
                  required={this.props.isRequired}
                />
              </td>
            </tr>
          </tbody>
        </table>
      </FormSection>
    );
  }
}

DimensionsField.propTypes = {
  fieldName: PropTypes.string.isRequired,
  labelText: PropTypes.string.isRequired,
  swagger: PropTypes.object.isRequired,
  isRequired: PropTypes.bool.isRequired,
};
