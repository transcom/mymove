import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './Address.css';

export class ShipmentAddress extends Component {
  render() {
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
      <div className="usa-grid">
        <h3 className="instruction-heading">
          Now let's review your pickup and delivery locations
        </h3>

        <h4>Pickup Location</h4>

        <div className="usa-width-one-whole">
          <div className="address-segment">
            <SwaggerField
              fieldName="street_address_1"
              swagger={this.props.addressSchema}
              required
            />
            <SwaggerField
              fieldName="street_address_2"
              swagger={this.props.addressSchema}
            />
            <SwaggerField
              className="usa-width-one-fourth"
              fieldName="city"
              swagger={this.props.addressSchema}
              required
            />
            <SwaggerField
              className="usa-width-one-sixth"
              fieldName="state"
              swagger={this.props.addressSchema}
              required
            />
            <SwaggerField
              className="usa-width-one-fourth"
              fieldName="postal_code"
              swagger={this.props.addressSchema}
              required
            />
          </div>
          <SwaggerField
            className="radio-title"
            fieldName="has_secondary_pickup_address"
            swagger={this.props.schema}
            component={YesNoBoolean}
          />
          {hasSecondary && (
            <Fragment>
              <div className="address-segment">
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={this.props.addressSchema}
                  required={hasSecondary}
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={this.props.addressSchema}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={this.props.addressSchema}
                  required={hasSecondary}
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={this.props.addressSchema}
                  required={hasSecondary}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={this.props.addressSchema}
                  required={hasSecondary}
                />
              </div>
            </Fragment>
          )}
          <h4>Delivery location</h4>
          <SwaggerField
            className="radio-title"
            fieldName="has_delivery_address"
            swagger={this.props.schema}
            component={YesNoBoolean}
          />
          {hasDelivery && (
            <Fragment>
              <div className="address-segment">
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={this.props.addressSchema}
                  required={hasDelivery}
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={this.props.addressSchema}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={this.props.addressSchema}
                  required={hasDelivery}
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={this.props.addressSchema}
                  required={hasDelivery}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={this.props.addressSchema}
                  required={hasDelivery}
                />
              </div>
            </Fragment>
          )}
        </div>
      </div>
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
    addressSchema: get(state, 'swagger.spec.definitions.Address', {}),
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(ShipmentAddress);
