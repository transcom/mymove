import React, { Component, Fragment } from 'react';
import PropTypes from 'prop-types';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export class DimensionsField extends Component {
  render() {
    return (
      <Fragment>
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
                  fieldName={this.props.fieldName}
                  swagger={this.props.swagger}
                  dimensionComponent="length"
                  hideLabel={true}
                  className="dimensions-form-input"
                />
              </td>
              <td width="11%" className="multiplication-sign">
                x
              </td>
              <td width="26%" className="dimensions-form-input-cell">
                <SwaggerField
                  fieldName={this.props.fieldName}
                  swagger={this.props.swagger}
                  dimensionComponent="width"
                  hideLabel={true}
                  className="dimensions-form-input"
                />
              </td>
              <td width="11%" className="multiplication-sign">
                x
              </td>
              <td width="26%" className="dimensions-form-input-cell">
                <SwaggerField
                  fieldName={this.props.fieldName}
                  swagger={this.props.swagger}
                  dimensionComponent="height"
                  hideLabel={true}
                  className="dimensions-form-input"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </Fragment>
    );
  }
}

export const Code105Details = props => {
  return (
    <div>
      <DimensionsField fieldName="item_dimensions" swagger={props.swagger} labelText="Item Dimensions (inches)" />
      <DimensionsField fieldName="crate_dimensions" swagger={props.swagger} labelText="Crate Dimensions (inches)" />
      <div className="bq-explanation">
        <p>Crate can only exceed item size by:</p>
        <ul>
          <li>
            <em>Internal crate</em>: Up to 3" larger
          </li>
          <li>
            <em>External crate</em>: Up to 5" larger
          </li>
        </ul>
      </div>
    </div>
  );
};

DimensionsField.propTypes = {
  fieldName: PropTypes.string.isRequired,
  labelText: PropTypes.string.isRequired,
  swagger: PropTypes.object.isRequired,
};
