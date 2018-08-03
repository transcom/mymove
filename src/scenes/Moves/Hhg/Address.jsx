import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { FormSection } from 'redux-form';

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
    const addressSchema = get(
      this.props,
      'schema.properties.pickup_address',
      {},
    );

    return (
      <div className="form-section">
        <h3 className="instruction-heading">
          Now let's review your pickup and delivery locations
        </h3>

        <div className="usa-grid">
          <div className="usa-width-one-whole">
            <FormSection name="pickup_address">
              <div className="address-segment usa-grid">
                <SwaggerField
                  fieldName="street_address_1"
                  swagger={addressSchema}
                  required
                />
                <SwaggerField
                  fieldName="street_address_2"
                  swagger={addressSchema}
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="city"
                  swagger={addressSchema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-sixth"
                  fieldName="state"
                  swagger={addressSchema}
                  required
                />
                <SwaggerField
                  className="usa-width-one-fourth"
                  fieldName="postal_code"
                  swagger={addressSchema}
                  required
                />
              </div>
            </FormSection>
            <SwaggerField
              className="radio-title"
              fieldName="has_secondary_pickup_address"
              swagger={this.props.schema}
              component={YesNoBoolean}
            />
            {hasSecondary && (
              <Fragment>
                <FormSection name="secondary_pickup_address">
                  <div className="address-segment usa-grid">
                    <SwaggerField
                      fieldName="street_address_1"
                      swagger={addressSchema}
                      required={hasSecondary}
                    />
                    <SwaggerField
                      fieldName="street_address_2"
                      swagger={addressSchema}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="city"
                      swagger={addressSchema}
                      required={hasSecondary}
                    />
                    <SwaggerField
                      className="usa-width-one-sixth"
                      fieldName="state"
                      swagger={addressSchema}
                      required={hasSecondary}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="postal_code"
                      swagger={addressSchema}
                      required={hasSecondary}
                    />
                  </div>
                </FormSection>
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
                <FormSection name="delivery_address">
                  <div className="address-segment usa-grid">
                    <SwaggerField
                      fieldName="street_address_1"
                      swagger={addressSchema}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      fieldName="street_address_2"
                      swagger={addressSchema}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="city"
                      swagger={addressSchema}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      className="usa-width-one-sixth"
                      fieldName="state"
                      swagger={addressSchema}
                      required={hasDelivery}
                    />
                    <SwaggerField
                      className="usa-width-one-fourth"
                      fieldName="postal_code"
                      swagger={addressSchema}
                      required={hasDelivery}
                    />
                  </div>
                </FormSection>
              </Fragment>
            )}
          </div>
        </div>
      </div>
    );
  }
}
ShipmentAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  error: PropTypes.object,
  formValues: PropTypes.object,
};

export default ShipmentAddress;
