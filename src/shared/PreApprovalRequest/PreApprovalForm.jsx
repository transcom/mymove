import { get, includes, reject } from 'lodash';
import React, { Component, Fragment } from 'react';
import Select, { createFilter } from 'react-select';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm, Form, Field } from 'redux-form';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './PreApprovalRequest.css';

const getOptionValue = option => (option ? option.id : null);
const getOptionLabel = option => (option ? option.code + ' ' + option.item : '');
const stringify = option => option.label;
const filterOption = createFilter({ ignoreCase: true, stringify });
const sitCodes = ['17A', '17B', '17C', '17D', '17E', '17F', '17G', '185A', '185B', '210D', '210E'];

export class Tariff400ngItemSearch extends Component {
  constructor(props) {
    super(props);
    this.localOnChange = this.localOnChange.bind(this);
  }

  localOnChange(value) {
    if (value && value.id) {
      this.props.input.onChange(value);
      return value.id;
    } else {
      this.props.input.onChange(null);
      return null;
    }
  }
  render() {
    // Filtering out SIT-related codes until SIT support is fully implemented
    // Remove this when SIT is completely supported
    let filteredOptions = reject(this.props.tariff400ngItems, item => {
      return includes(sitCodes, item.code);
    });
    return (
      <Fragment>
        <label className="usa-input-label">{this.props.title}</label>
        <Select
          options={filteredOptions}
          getOptionLabel={getOptionLabel}
          getOptionValue={getOptionValue}
          value={this.props.input.value || null}
          onChange={this.localOnChange}
          placeholder="Select an item..."
          className={`tariff400-select ${this.props.input.name}`}
          classNamePrefix="tariff400"
          filterOption={filterOption}
          defaultValue={this.props.meta.initial || null}
        />
      </Fragment>
    );
  }
}

export class PreApprovalForm extends Component {
  render() {
    return (
      <Form onSubmit={this.props.handleSubmit(this.props.onSubmit)}>
        <div className="usa-grid-full">
          <div className="usa-width-one-third">
            <div className="tariff400-select">
              <Field
                name="tariff400ng_item"
                title="Code & Item"
                component={Tariff400ngItemSearch}
                tariff400ngItems={this.props.tariff400ngItems}
              />
            </div>
            {/* TODO andrea - set schema location enum array to tariff400ng_item selected location value */}
            <SwaggerField
              fieldName="location"
              className="rounded"
              swagger={this.props.ship_line_item_schema}
              required
            />
          </div>
          <div className="usa-width-one-third">
            <SwaggerField
              fieldName="quantity_1"
              className="half-width"
              swagger={this.props.ship_line_item_schema}
              required
            />
            <div className="bq-explanation">
              <p>
                Enter numbers only, no symbols or units. <em>Examples:</em>
              </p>
              <ul>
                <li>
                  Crating: enter "<strong>47.4</strong>" for crate size of 47.4 cu. ft.
                </li>
                <li>
                  {' '}
                  3rd-party service: enter "<strong>1299.99</strong>" for cost of $1,299.99.
                </li>
                <li>
                  Bulky item: enter "<strong>1</strong>" for a single item.
                </li>
              </ul>
            </div>
          </div>
          <div className="usa-width-one-third">
            <SwaggerField
              fieldName="notes"
              className="three-quarter-width"
              swagger={this.props.ship_line_item_schema}
            />
          </div>
        </div>
      </Form>
    );
  }
}

PreApprovalForm.propTypes = {
  tariff400ngItems: PropTypes.array,
  onSubmit: PropTypes.func.isRequired,
};

const validateItemSelect = validateAdditionalFields(['tariff400ng_item']);
export const formName = 'preapproval_request_form';

PreApprovalForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
  validate: validateItemSelect,
})(PreApprovalForm);

function mapStateToProps(state, props) {
  return {
    ship_line_item_schema: get(state, 'swaggerPublic.spec.definitions.ShipmentLineItem', {}),
  };
}

export default connect(mapStateToProps)(PreApprovalForm);
