import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './Address.css';

const formName = 'hhg_address';
const AddressWizardForm = reduxifyWizardForm(formName);
const hhgSchema = {
  properties: {
    has_secondary_pickup_address: {
      type: 'boolean',
      title: 'Do you have household goods at any other pickup location?',
      'x-nullable': true,
    },
    has_delivery_address: {
      type: 'boolean',
      title: 'Do you know your home address at your destination yet?',
      'x-nullable': true,
    },
    has_partial_sit_delivery_address: {
      type: 'boolean',
      title:
        'Do you want to deliver some of your household goods to an additional destination (such as a self-storage unit)?',
      'x-nullable': true,
    },
  },
};

export class HHGAddress extends Component {
  handleSubmit = () => {
    // Create new HHG pickup address
    // const newAddress = { residential_address: this.props.values };
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
    } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = get(currentServiceMember, 'residential_address');
    const hasSecondary = get(
      this.props,
      'formValues.has_secondary_pickup_address',
      false,
    );
    const hasDelivery = get(
      this.props,
      'formValues.has_delivery_address',
      false,
    );
    const hasSITDelivery = get(
      this.props,
      'formValues.has_partial_sit_delivery_address',
      false,
    );
    return (
      <AddressWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
      >
        <div className="usa-grid">
          <h3>Shipment 1 (HHG)</h3>
          <h3 className="instruction-heading">
            Now let's review your pickup and delivery locations
          </h3>

          <h4>Pickup Location</h4>

          <div className="usa-width-one-whole">
            <div className="address-segment">
              <SwaggerField
                fieldName="street_address_1"
                swagger={this.props.schema}
                required
              />
              <SwaggerField
                fieldName="street_address_2"
                swagger={this.props.schema}
              />
              <SwaggerField
                className="usa-width-one-fourth"
                fieldName="city"
                swagger={this.props.schema}
                required
              />
              <SwaggerField
                className="usa-width-one-sixth"
                fieldName="state"
                swagger={this.props.schema}
                required
              />
              <SwaggerField
                className="usa-width-one-fourth"
                fieldName="postal_code"
                swagger={this.props.schema}
                required
              />
            </div>
            <SwaggerField
              className="radio-title"
              fieldName="has_secondary_pickup_address"
              swagger={hhgSchema}
              component={YesNoBoolean}
            />
            {hasSecondary && (
              <Fragment>
                <div className="address-segment">
                  <SwaggerField
                    fieldName="street_address_1"
                    swagger={this.props.secondaryAddressSchema}
                    required={hasSecondary}
                  />
                  <SwaggerField
                    fieldName="street_address_2"
                    swagger={this.props.secondaryAddressSchema}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="city"
                    swagger={this.props.secondaryAddressSchema}
                    required={hasSecondary}
                  />
                  <SwaggerField
                    className="usa-width-one-sixth"
                    fieldName="state"
                    swagger={this.props.secondaryAddressSchema}
                    required={hasSecondary}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="postal_code"
                    swagger={this.props.secondaryAddressSchema}
                    required={hasSecondary}
                  />
                </div>
              </Fragment>
            )}
            <h4>Delivery location</h4>
            <SwaggerField
              className="radio-title"
              fieldName="has_delivery_address"
              swagger={hhgSchema}
              component={YesNoBoolean}
            />
            {hasDelivery && (
              <Fragment>
                <div className="address-segment">
                  <SwaggerField
                    fieldName="street_address_1"
                    swagger={this.props.deliveryAddressSchema}
                    required={hasDelivery}
                  />
                  <SwaggerField
                    fieldName="street_address_2"
                    swagger={this.props.deliveryAddressSchema}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="city"
                    swagger={this.props.deliveryAddressSchema}
                    required={hasDelivery}
                  />
                  <SwaggerField
                    className="usa-width-one-sixth"
                    fieldName="state"
                    swagger={this.props.deliveryAddressSchema}
                    required={hasDelivery}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="postal_code"
                    swagger={this.props.deliveryAddressSchema}
                    required={hasDelivery}
                  />
                </div>
              </Fragment>
            )}
            <SwaggerField
              className="radio-title"
              fieldName="has_partial_sit_delivery_address"
              swagger={hhgSchema}
              component={YesNoBoolean}
            />
            {hasSITDelivery && (
              <Fragment>
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={this.props.partialSITAddressSchema}
                  required={hasSITDelivery}
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={this.props.partialSITAddressSchema}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={this.props.partialSITAddressSchema}
                  required={hasSITDelivery}
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={this.props.partialSITAddressSchema}
                  required={hasSITDelivery}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={this.props.partialSITAddressSchema}
                  required={hasSITDelivery}
                />
              </Fragment>
            )}
          </div>
        </div>
      </AddressWizardForm>
    );
  }
}
HHGAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
function mapStateToProps(state) {
  return {
    schema: get(state, 'swagger.spec.definitions.Address', {}),
    secondaryAddressSchema: get(state, 'swagger.spec.definitions.Address', {}),
    deliveryAddressSchema: get(state, 'swagger.spec.definitions.Address', {}),
    partialSITAddressSchema: get(state, 'swagger.spec.definitions.Address', {}),
    formValues: getFormValues(formName)(state),
    ...state.serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(HHGAddress);
