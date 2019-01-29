import { get, includes, reject } from 'lodash';
import React, { Component, Fragment } from 'react';
import Select, { createFilter } from 'react-select';
import { connect } from 'react-redux';
import { withContext } from 'shared/AppContext';
import PropTypes from 'prop-types';
import { reduxForm, Form, Field, formValueSelector } from 'redux-form';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { getDetailComponent } from './DetailsHelper';
import { selectLocationFromTariff400ngItem } from 'shared/Entities/modules/shipmentLineItems';

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
  componentDidMount() {
    this.updateLocationValue();
  }

  componentDidUpdate() {
    this.updateLocationValue();
  }

  updateLocationValue() {
    if (
      this.props.filteredLocations &&
      this.props.filteredLocations.length === 1 &&
      this.props.filteredLocations[0] !== this.props.value
    ) {
      this.props.change('location', this.props.filteredLocations[0]);
    }
  }

  render() {
    return this.props.filteredLocations && this.props.filteredLocations.length === 1 ? (
      <Fragment>
        <label htmlFor="location" className="usa-input-label">
          Location
        </label>
        <div>
          <strong>
            {this.props.ship_line_item_schema.properties.location['x-display-value'][this.props.filteredLocations[0]]}
          </strong>
        </div>
      </Fragment>
    ) : (
      <SwaggerField
        fieldName="location"
        className="rounded"
        swagger={this.props.ship_line_item_schema}
        filteredEnumListOverride={this.props.filteredLocations}
        required
      />
    );
  }
}

export class PreApprovalForm extends Component {
  render() {
    const robustAccessorial = get(this.props, 'context.flags.robustAccessorial', false);
    const DetailComponent = getDetailComponent(this.props.tariff400ng_item_code, robustAccessorial);

    return (
      <Form className="pre-approval-form" onSubmit={this.props.handleSubmit(this.props.onSubmit)}>
        <div className="usa-grid-full">
          <div className="usa-width-one-third">
            <div className="tariff400-select usa-input">
              <Field
                name="tariff400ng_item"
                title="Code & Item"
                component={Tariff400ngItemSearch}
                tariff400ngItems={this.props.tariff400ngItems}
              />
            </div>
            {this.props.tariff400ngItem && (
              <div className="location-select">
                <LocationSearch
                  filteredLocations={this.props.filteredLocations}
                  ship_line_item_schema={this.props.ship_line_item_schema}
                  change={this.props.change}
                  value={this.props.selectedLocation}
                />
              </div>
            )}
          </div>
          {this.props.tariff400ngItem && (
            <Fragment>
              <div className="usa-width-one-third">
                <DetailComponent {...this.props} />
              </div>
              <div className="usa-width-one-third">
                <SwaggerField fieldName="notes" swagger={this.props.ship_line_item_schema} />
              </div>
            </Fragment>
          )}
        </div>
      </Form>
    );
  }
}

PreApprovalForm.propTypes = {
  tariff400ngItems: PropTypes.array,
  onSubmit: PropTypes.func.isRequired,
};

LocationSearch.propTypes = {
  filteredLocations: PropTypes.arrayOf(PropTypes.string),
  change: PropTypes.func,
  ship_line_item_schema: PropTypes.object,
};

const validateItemSelect = validateAdditionalFields(['tariff400ng_item']);
export const formName = 'preapproval_request_form';

PreApprovalForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
  validate: validateItemSelect,
})(PreApprovalForm);

const selector = formValueSelector(formName);
function mapStateToProps(state) {
  return {
    tariff400ng_item_code: get(state, 'form.preapproval_request_form.values.tariff400ng_item.code'),
    ship_line_item_schema: get(state, 'swaggerPublic.spec.definitions.ShipmentLineItem', {}),
    filteredLocations: selectLocationFromTariff400ngItem(state, selector(state, 'tariff400ng_item')),
    selectedLocation: selector(state, 'location'),
    tariff400ngItem: selector(state, 'tariff400ng_item'),
  };
}

export default withContext(connect(mapStateToProps)(PreApprovalForm));
