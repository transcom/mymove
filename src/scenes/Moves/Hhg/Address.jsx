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

const formName = 'shipment_address';
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
  },
};

export class ShipmentAddress extends Component {
  // handleSubmit = () => {
  //   return false;
  //   // Create new HHG primary (and secondary) pickup addresses, delivery address if it exists
  //   // const newAddresses = { ...this.props.values };
  //   // this.props.updateShipment(newAddress);
  // };

  render() {
    const handleSubmit = false;
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
                    swagger={this.props.schema}
                    required={hasSecondary}
                  />
                  <SwaggerField
                    fieldName="street_address_2"
                    swagger={this.props.schema}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="city"
                    swagger={this.props.schema}
                    required={hasSecondary}
                  />
                  <SwaggerField
                    className="usa-width-one-sixth"
                    fieldName="state"
                    swagger={this.props.schema}
                    required={hasSecondary}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="postal_code"
                    swagger={this.props.schema}
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
                    swagger={this.props.schema}
                    required={hasDelivery}
                  />
                  <SwaggerField
                    fieldName="street_address_2"
                    swagger={this.props.schema}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="city"
                    swagger={this.props.schema}
                    required={hasDelivery}
                  />
                  <SwaggerField
                    className="usa-width-one-sixth"
                    fieldName="state"
                    swagger={this.props.schema}
                    required={hasDelivery}
                  />
                  <SwaggerField
                    className="usa-width-one-fourth"
                    fieldName="postal_code"
                    swagger={this.props.schema}
                    required={hasDelivery}
                  />
                </div>
              </Fragment>
            )}
          </div>
        </div>
      </AddressWizardForm>
    );
  }
}
ShipmentAddress.propTypes = {
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
    formValues: getFormValues(formName)(state),
    ...state.serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(ShipmentAddress);
