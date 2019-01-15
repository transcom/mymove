import { get, includes, reject } from 'lodash';
import React, { Component, Fragment } from 'react';
import Select, { createFilter } from 'react-select';
import { connect } from 'react-redux';
import { withContext } from 'shared/AppContext';
import PropTypes from 'prop-types';
import { reduxForm, Form, Field } from 'redux-form';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { getDetailComponent } from './DetailsHelper';
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
          placeholder="Enter code or item"
          className={`tariff400-select ${this.props.input.name}`}
          classNamePrefix="tariff400"
          filterOption={filterOption}
          defaultValue={this.props.meta.initial || null}
        />
      </Fragment>
    );
  }
}

export class LocationSearch extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    if (this.props.options.length === 1) {
      this.props.input.onChange(this.props.options[0].value);
    }
  }

  componentDidUpdate() {
    if (this.props.options.length === 1) {
      this.props.input.onChange(this.props.options[0].value);
    }
  }

  render() {
    let content;
    if (this.props.options.length === 0) {
      content = null;
    } else if (this.props.options.length === 1) {
      content = this.props.options[0].label;
    } else {
      content = (
        <Fragment>
          <Select
            options={this.props.options}
            getOptionLabel={option => option.label}
            getOptionValue={option => option.value}
            value={this.props.input.value || null}
            onChange={this.props.input.onChange}
            placeholder="Select a location..."
            className={`location-select ${this.props.input.name} ${this.props.className}`}
            classNamePrefix="location"
            filterOption={filterOption}
            defaultValue={this.props.meta.initial || null}
            isSearchable={false}
          />
        </Fragment>
      );
    }
    return (
      <Fragment>
        <label className="usa-input-label">{this.props.title}</label>
        {content}
      </Fragment>
    );
  }
}
export class PreApprovalForm extends Component {
  constructor(props) {
    super(props);
    this.filterLocations = this.filterLocations.bind(this);
  }

  filterLocations() {
    const lineItemLocations = this.props.ship_line_item_schema.properties.location;
    if (!lineItemLocations.enum) {
      return [];
    }
    const lineItemLocationMap = lineItemLocations.enum.map(lineItemLocation => {
      return {
        value: lineItemLocation,
        label: lineItemLocations['x-display-value'][lineItemLocation],
      };
    });
    // Choose location options based on tariff400ng choice.
    if (this.props.tariff400ngItemLocation === 'EITHER') {
      return lineItemLocationMap.filter(lineItemLocationObject => {
        return lineItemLocationObject.value === 'ORIGIN' || lineItemLocationObject.value === 'DESTINATION';
      });
    } else {
      return lineItemLocationMap.filter(lineItemLocationObject => {
        return lineItemLocationObject.value === this.props.tariff400ngItemLocation;
      });
    }
  }

  render() {
    const DetailComponent = getDetailComponent(
      this.props.tariff400ng_item_code,
      this.props.context.flags.robustAccessorial,
    );

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
            <DetailComponent {...this.props} />
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
    tariff400ng_item_code: get(state, 'form.preapproval_request_form.values.tariff400ng_item.code'),
    ship_line_item_schema: get(state, 'swaggerPublic.spec.definitions.ShipmentLineItem', {}),
    tariff400ngItemLocation: get(state, 'form.preapproval_request_form.values.tariff400ng_item.location'),
  };
}

export default withContext(connect(mapStateToProps)(PreApprovalForm));
